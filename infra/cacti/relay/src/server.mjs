// Servidor Hyperledger Cacti (v3.0.1) do Janus (Abordagem B2). Hospeda, conforme a configuração:
//
// - conector Fabric, com a identidade da org atestadora (Registradora ou Banco). O atestador usa o
//   conector como fonte de prova do lock no ledger comitado (ADR-0005);
// - conector Besu, no relay da Vendedora. Submete o AttestedSettlement.settle(S, atestação) com a
//   conta de relay, cuja chave fica no keychain do servidor (ADR-0006).
//
// Uso: node src/server.mjs <config.json>
//   {
//     "apiHost": "127.0.0.1", "apiPort": 4000, "keychainId": "janus",
//     "fabric": {                                    // opcional
//       "connectionProfile": "<caminho do ccp.json>", "keychainRef": "registradora",
//       "identity": { "mspId": "RegistradoraMSP", "certificate": "<PEM>", "privateKey": "<PEM>" }
//     },
//     "besu": {                                      // opcional
//       "rpcApiHttpHost": "http://127.0.0.1:8545", "rpcApiWsHost": "ws://127.0.0.1:8546",
//       "keys": { "relay": "<chave privada hex>" }   // keychainEntryKey -> chave
//     }
//   }
//
// Compatibilidade: as chaves "connectionProfile", "keychainRef" e "identity" no topo do arquivo
// (formato do ADR-0005) continuam valendo como seção "fabric".
import { readFile } from "node:fs/promises";
import { createRequire } from "node:module";

const require = createRequire(import.meta.url);
const { ApiServer, ConfigService, AuthorizationProtocol } = require("@hyperledger-cacti/cactus-cmd-api-server");
const { PluginRegistry } = require("@hyperledger-cacti/cactus-core");
const { PluginKeychainMemory } = require("@hyperledger-cacti/cactus-plugin-keychain-memory");

const cfgPath = process.argv[2];
if (!cfgPath) {
  console.error("uso: node src/server.mjs <config.json>");
  process.exit(2);
}
const cfg = JSON.parse(await readFile(cfgPath, "utf8"));
const logLevel = cfg.logLevel ?? "WARN";
const fabricCfg = cfg.fabric ?? (cfg.connectionProfile
  ? { connectionProfile: cfg.connectionProfile, keychainRef: cfg.keychainRef, identity: cfg.identity }
  : undefined);
const besuCfg = cfg.besu;
if (!fabricCfg && !besuCfg) {
  console.error("configuração sem conector: informe 'fabric' e/ou 'besu'");
  process.exit(2);
}

// Keychain em memória: aceitável só em DEV (em produção, Vault/KMS; ver runbook).
const keychain = new PluginKeychainMemory({ instanceId: "janus-keychain", keychainId: cfg.keychainId, logLevel });
const pluginRegistry = new PluginRegistry({ plugins: [keychain] });
const labels = [];

if (fabricCfg) {
  const { PluginLedgerConnectorFabric, DefaultEventHandlerStrategy } = require(
    "@hyperledger-cacti/cactus-plugin-ledger-connector-fabric",
  );
  await keychain.set(
    fabricCfg.keychainRef,
    JSON.stringify({
      type: "X.509",
      mspId: fabricCfg.identity.mspId,
      credentials: { certificate: fabricCfg.identity.certificate, privateKey: fabricCfg.identity.privateKey },
    }),
  );
  const ccp = JSON.parse(await readFile(fabricCfg.connectionProfile, "utf8"));
  pluginRegistry.add(
    new PluginLedgerConnectorFabric({
      instanceId: `janus-fabric-${fabricCfg.keychainRef}`,
      pluginRegistry,
      connectionProfile: ccp,
      logLevel,
      discoveryOptions: { enabled: true, asLocalhost: true },
      eventHandlerOptions: { strategy: DefaultEventHandlerStrategy.NetworkScopeAnyfortx, commitTimeout: 300 },
    }),
  );
  labels.push(`fabric(${fabricCfg.identity.mspId})`);
}

if (besuCfg) {
  const { PluginLedgerConnectorBesu } = require("@hyperledger-cacti/cactus-plugin-ledger-connector-besu");
  for (const [ref, key] of Object.entries(besuCfg.keys ?? {})) {
    await keychain.set(ref, key.startsWith("0x") ? key : `0x${key}`);
  }
  pluginRegistry.add(
    new PluginLedgerConnectorBesu({
      instanceId: "janus-besu",
      rpcApiHttpHost: besuCfg.rpcApiHttpHost,
      rpcApiWsHost: besuCfg.rpcApiWsHost,
      pluginRegistry,
      logLevel,
    }),
  );
  labels.push(`besu(${besuCfg.rpcApiHttpHost})`);
}

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
console.log(`janus-cacti pronto: http://${opts.apiHost}:${opts.apiPort} [${labels.join(", ")}]`);
