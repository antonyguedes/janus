// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.30;

import {Script, console2} from "forge-std/Script.sol";
import {DrexToken} from "../src/DrexToken.sol";
import {HTLCPayment} from "../src/HTLCPayment.sol";
import {FabricProofVerifier} from "../src/FabricProofVerifier.sol";
import {AttestedSettlement} from "../src/AttestedSettlement.sol";

/// @notice Deploy local no Besu QBFT (ou anvil):
///   forge script script/Deploy.s.sol --rpc-url http://localhost:8545 --broadcast --legacy --with-gas-price 0
/// Variáveis (padrões = contas de DEV do anvil, também na allowlist do Besu):
///   DEPLOYER_PK, SELLER_ADDR, BUYER_ADDR, ATTESTOR_ADDR, RELAY_ADDR, BUYER_MINT (centavos),
///   FABRIC_CHANNEL (recebiveis), FABRIC_NAMESPACE (tcc)
/// Escreve ../deployments/<chainId>.local.json.
contract Deploy is Script {
    function run() external {
        uint256 pk = vm.envOr(
            "DEPLOYER_PK", uint256(0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80)
        );
        address seller = vm.envOr("SELLER_ADDR", address(0x70997970C51812dc3A010C7d01b50e0d17dc79C8));
        address buyer = vm.envOr("BUYER_ADDR", address(0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC));
        address attestor = vm.envOr("ATTESTOR_ADDR", address(0x90F79bf6EB2c4f870365E785982E1f101E93b906));
        address relay = vm.envOr("RELAY_ADDR", address(0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65));
        uint256 mint = vm.envOr("BUYER_MINT", uint256(100_000_000_00)); // R$ 100 mi
        string memory channel = vm.envOr("FABRIC_CHANNEL", string("recebiveis"));
        string memory ns = vm.envOr("FABRIC_NAMESPACE", string("tcc"));
        address admin = vm.addr(pk);

        vm.startBroadcast(pk);
        DrexToken drex = new DrexToken(admin);
        HTLCPayment htlc = new HTLCPayment(drex);
        FabricProofVerifier verifier = new FabricProofVerifier(admin);
        AttestedSettlement settlement =
            new AttestedSettlement(drex, verifier, sha256(bytes(channel)), sha256(bytes(ns)));

        drex.setAllowed(seller, true);
        drex.setAllowed(buyer, true);
        drex.setAllowed(attestor, true);
        drex.setAllowed(relay, true);
        drex.setAllowed(address(htlc), true);
        drex.setAllowed(address(settlement), true);
        drex.mint(buyer, mint);
        vm.stopBroadcast();

        console2.log("DrexToken          ", address(drex));
        console2.log("HTLCPayment        ", address(htlc));
        console2.log("FabricProofVerifier", address(verifier));
        console2.log("AttestedSettlement ", address(settlement));

        string memory o = "deploy";
        vm.serializeUint(o, "chainId", block.chainid);
        vm.serializeAddress(o, "DrexToken", address(drex));
        vm.serializeAddress(o, "HTLCPayment", address(htlc));
        vm.serializeAddress(o, "FabricProofVerifier", address(verifier));
        vm.serializeString(o, "fabricChannel", channel);
        vm.serializeString(o, "fabricNamespace", ns);
        string memory json = vm.serializeAddress(o, "AttestedSettlement", address(settlement));
        vm.writeJson(json, string.concat(vm.projectRoot(), "/../deployments/", vm.toString(block.chainid), ".local.json"));
    }
}
