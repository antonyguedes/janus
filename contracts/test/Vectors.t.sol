// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.30;

import {Base} from "./Base.t.sol";
import {FabricProofVerifier, FabricLockAttestation, OrgSignature} from "../src/FabricProofVerifier.sol";

/// @notice Vetores compartilhados com Go (gerados por `go run ./cmd/vectors`, ver testdata/vectors).
///         Se este teste e `go test ./internal/...` passarem com o mesmo JSON, o hashing e o digest de
///         atestação são interoperáveis.
contract VectorsTest is Base {
    function test_htlcVectors() public view {
        string memory json = vm.readFile(string.concat(vm.projectRoot(), "/../testdata/vectors/htlc.json"));
        uint256 n = _count(json, ".vectors");
        assertGt(n, 0);
        for (uint256 i = 0; i < n; ++i) {
            string memory k = string.concat(".vectors[", vm.toString(i), "]");
            bytes32 pre = vm.parseJsonBytes32(json, string.concat(k, ".preimage"));
            bytes32 img = vm.parseJsonBytes32(json, string.concat(k, ".hashlock"));
            assertEq(sha256(abi.encodePacked(pre)), img, "vetor htlc divergente");
        }
    }

    /// @dev O JSONPath do Foundry não aceita curingas: conta os elementos por índice.
    function _count(string memory json, string memory arr) internal view returns (uint256 n) {
        while (vm.keyExistsJson(json, string.concat(arr, "[", vm.toString(n), "]"))) ++n;
    }

    function _loadAttestation(string memory json) internal pure returns (FabricLockAttestation memory a) {
        a.swapId = vm.parseJsonBytes32(json, ".swapId");
        a.hashlock = vm.parseJsonBytes32(json, ".hashlock");
        a.assetIdHash = vm.parseJsonBytes32(json, ".assetIdHash");
        a.recipientIdHash = vm.parseJsonBytes32(json, ".recipientIdHash");
        a.deadline = uint64(vm.parseJsonUint(json, ".deadline"));
        a.channelHash = vm.parseJsonBytes32(json, ".channelHash");
        a.namespaceHash = vm.parseJsonBytes32(json, ".namespaceHash");
        a.fabricTxId = vm.parseJsonBytes32(json, ".fabricTxId");
        a.blockNumber = uint64(vm.parseJsonUint(json, ".blockNumber"));
    }

    function test_attestationVector() public {
        string memory json =
            vm.readFile(string.concat(vm.projectRoot(), "/../testdata/vectors/attestation.json"));
        vm.chainId(vm.parseJsonUint(json, ".chainId"));
        address at = vm.parseJsonAddress(json, ".verifier");
        deployCodeTo("FabricProofVerifier.sol:FabricProofVerifier", abi.encode(admin), at);
        FabricProofVerifier v = FabricProofVerifier(at);

        assertEq(v.DOMAIN_TAG(), vm.parseJsonBytes32(json, ".domainTag"), "domain tag");
        assertEq(sha256(bytes(vm.parseJsonString(json, ".assetId"))), vm.parseJsonBytes32(json, ".assetIdHash"));
        assertEq(sha256(bytes(vm.parseJsonString(json, ".channel"))), vm.parseJsonBytes32(json, ".channelHash"));
        assertEq(
            sha256(abi.encodePacked(vm.parseJsonBytes32(json, ".preimage"))), vm.parseJsonBytes32(json, ".hashlock")
        );

        FabricLockAttestation memory a = _loadAttestation(json);
        assertEq(v.digest(a), vm.parseJsonBytes32(json, ".digest"), "digest Go != Solidity");

        uint256 n = _count(json, ".orgs");
        assertEq(n, 2);
        OrgSignature[] memory sigs = new OrgSignature[](n);
        vm.startPrank(admin);
        for (uint256 i = 0; i < n; ++i) {
            string memory o = string.concat(".orgs[", vm.toString(i), "]");
            bytes32 h = v.registerOrg(vm.parseJsonString(json, string.concat(o, ".mspId")));
            assertEq(h, vm.parseJsonBytes32(json, string.concat(o, ".mspIdHash")), "keccak(mspId) Go != Solidity");
            uint32 k = v.addAttestorKey(
                h,
                vm.parseJsonBytes32(json, string.concat(o, ".qx")),
                vm.parseJsonBytes32(json, string.concat(o, ".qy")),
                bytes32(0)
            );
            sigs[i] = OrgSignature({
                mspIdHash: h,
                keyIndex: k,
                r: vm.parseJsonBytes32(json, string.concat(o, ".r")),
                s: vm.parseJsonBytes32(json, string.concat(o, ".s"))
            });
        }
        v.setThreshold(uint32(n));
        vm.stopPrank();

        assertTrue(v.isValid(a, sigs), "assinaturas Go nao verificam no Solidity");

        a.deadline += 1; // qualquer alteração invalida
        assertFalse(v.isValid(a, sigs));
    }
}
