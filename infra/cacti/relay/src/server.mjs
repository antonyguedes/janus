// Servidor Hyperledger Cacti (v3.0.1) de uma org atestadora do Janus (Abordagem B2).
//
// Hospeda o conector Fabric do Cacti com a identidade da própria org. O atestador do nó Janus
// usa este servidor como fonte de prova do lock no ledger comitado: o recibo da transação (bloco e
// endossantes) e o TRANSACTIONS_FILTER do bloco (validade). Ver ADR-0005.
//
// Uso: node src/server.mjs <config.json>
//   {
//     "apiHost": "127.0.0.1", "apiPort": 4000,
//     "connectionProfile": "<caminho do ccp.json>",
//     "keychainId": "janus", "keychainRef": "registradora",
//     "identity": { "mspId": "RegistradoraMSP", "certificate": "<PEM>", "privateKey": "<PEM>" }
//   }
import { readFile } from "node:fs/promises";
import { createRequire } from "node:module";

const require = createRequire(import.meta.url);
const { ApiServer, ConfigService, AuthorizationProtocol } = require("@hyperledger-cacti/cactus-cmd-api-server");
const { PluginRegistry } = require("@hyperledger-cacti/cactus-core");
const { PluginKeychainMemory } = require("@hyperledger-cacti/cactus-plugin-keychain-memory");
const {
  PluginLedgerConnectorFabric,
  DefaultEventHandlerStrategy,
} = require("@hyperledger-cacti/cactus-plugin-ledger-connector-fabric");

const cfgPath = process.argv[2];
if (!cfgPath) {
  console.error("uso: node src/server.mjs <config.json>");
  process.exit(2);
}
const cfg = JSON.parse(await readFile(cfgPath, "utf8"));
const ccp = JSON.parse(await readFile(cfg.connectionProfile, "utf8"));
const logLevel = cfg.logLevel ?? "WARN";

// Keychain em memória: aceitável só em DEV (em produção, Vault/KMS; ver runbook).
const keychain = new PluginKeychainMemory({ instanceId: "janus-keychain", keychainId: cfg.keychainId, logLevel });
await keychain.set(
  cfg.keychainRef,
  JSON.stringify({
    type: "X.509",
    mspId: cfg.identity.mspId,
    credentials: { certificate: cfg.identity.certificate, privateKey: cfg.identity.privateKey },
  }),
);

const pluginRegistry = new PluginRegistry({ plugins: [keychain] });
const fabric = new PluginLedgerConnectorFabric({
  instanceId: `janus-fabric-${cfg.keychainRef}`,
  pluginRegistry,
  connectionProfile: ccp,
  logLevel,
  discoveryOptions: { enabled: true, asLocalhost: true },
  eventHandlerOptions: { strategy: DefaultEventHandlerStrategy.NetworkScopeAnyfortx, commitTimeout: 300 },
});
pluginRegistry.add(fabric);

const configService = new ConfigService();
const opts = await configService.newExampleConfig();
Object.assign(opts, {
  configFile: "",
  logLevel,
  authorizationProtocol: AuthorizationProtocol.NONE,
  apiHost: cfg.apiHost ?? "127.0.0.1",
  apiPort: cfg.apiPort,
  apiCorsDomainCsv: "*",
  apiTlsEnabled: false,
  apiMtlsEnabled: false,
  cockpitEnabled: false,
  grpcPort: 0,
  crpcPort: 0,
  plugins: [],
});
const config = await configService.newExampleConfigConvict(opts);
const apiServer = new ApiServer({ config: config.getProperties(), pluginRegistry });

const shutdown = async () => {
  try {
    await apiServer.shutdown();
  } finally {
    process.exit(0);
  }
};
process.on("SIGINT", shutdown);
process.on("SIGTERM", shutdown);

await apiServer.start();
console.log(`janus-cacti pronto: http://${opts.apiHost}:${opts.apiPort} (${cfg.identity.mspId})`);
