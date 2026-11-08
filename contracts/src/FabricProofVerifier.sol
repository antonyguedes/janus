// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.30;

import {AccessControl} from "@openzeppelin/contracts/access/AccessControl.sol";
import {P256} from "@openzeppelin/contracts/utils/cryptography/P256.sol";

/// @notice Atestação canônica de um lock HTLC do Panurus já comitado no Fabric (Abordagem B2, ADR-0003).
/// @dev Todos os campos são estáticos (32 bytes cada em `abi.encode`), o que torna o digest trivial
///      de reproduzir em Go (node/internal/attest).
struct FabricLockAttestation {
    bytes32 swapId;
    bytes32 hashlock; // H = sha256(S), igual ao HashInfo.Hash do script HTLC (encoding None)
    bytes32 assetIdHash; // sha256(ID do NFT do recibível)
    bytes32 recipientIdHash; // sha256(identidade FSC do recipient no script HTLC)
    uint64 deadline; // script.Deadline (T1), unix s
    bytes32 channelHash; // sha256(nome do canal)
    bytes32 namespaceHash; // sha256(namespace do Token Chaincode)
    bytes32 fabricTxId; // ID da transação de lock (hex de 32 bytes)
    uint64 blockNumber; // bloco do Fabric em que o lock foi comitado
}

/// @notice Assinatura P-256 de uma org atestadora.
struct OrgSignature {
    bytes32 mspIdHash; // keccak256(bytes(mspId))
    uint32 keyIndex;
    bytes32 r;
    bytes32 s; // low-s obrigatório (igual ao Fabric)
}

/// @title FabricProofVerifier — verifica atestações de lock no Fabric por política k-de-n orgs
/// @notice Cada org cadastrada tem uma ou mais chaves P-256 de atestador, emitidas pela CA da org
///         (X.509). A cadeia de certificados é validada no cadastro, feito por governança fora da cadeia.
///         O contrato guarda a chave pública e o hash do certificado, para auditoria.
///         A verificação usa o precompile P256VERIFY (EIP-7951 / RIP-7212, endereço 0x100) quando ele
///         existe, com fallback em Solidity puro (OpenZeppelin P256).
contract FabricProofVerifier is AccessControl {
    bytes32 public constant GOVERNOR_ROLE = keccak256("GOVERNOR_ROLE");

    /// @notice Separador de domínio: sha256("janus.fabric-lock-attestation.v1")
    bytes32 public constant DOMAIN_TAG = sha256("janus.fabric-lock-attestation.v1");

    struct Org {
        bool registered;
        bool active;
        uint32 keyCount;
        string mspId;
    }

    struct AttestorKey {
        bytes32 qx;
        bytes32 qy;
        bytes32 certHash; // sha256 do certificado X.509 (DER)
        bool active;
    }

    mapping(bytes32 mspIdHash => Org) private _orgs;
    mapping(bytes32 mspIdHash => mapping(uint32 keyIndex => AttestorKey)) private _keys;
    bytes32[] private _orgList;

    uint32 public activeOrgCount;
    uint32 public threshold;

    event OrgRegistered(bytes32 indexed mspIdHash, string mspId);
    event OrgStatusChanged(bytes32 indexed mspIdHash, bool active);
    event AttestorKeyAdded(bytes32 indexed mspIdHash, uint32 indexed keyIndex, bytes32 qx, bytes32 qy, bytes32 certHash);
    event AttestorKeyRevoked(bytes32 indexed mspIdHash, uint32 indexed keyIndex);
    event ThresholdChanged(uint32 threshold);

    error OrgAlreadyRegistered(bytes32 mspIdHash);
    error UnknownOrg(bytes32 mspIdHash);
    error UnknownKey(bytes32 mspIdHash, uint32 keyIndex);
    error InvalidPublicKey();
    error InvalidThreshold(uint32 threshold, uint32 activeOrgs);
    error SignaturesNotSorted();
    error InvalidSignature(bytes32 mspIdHash, uint32 keyIndex);
    error InactiveSigner(bytes32 mspIdHash, uint32 keyIndex);
    error ThresholdNotMet(uint32 valid, uint32 threshold);

    constructor(address admin) {
        _grantRole(DEFAULT_ADMIN_ROLE, admin);
        _grantRole(GOVERNOR_ROLE, admin);
    }

    // ------------------------------------------------------------------ governança

    function registerOrg(string calldata mspId) external onlyRole(GOVERNOR_ROLE) returns (bytes32 mspIdHash) {
        mspIdHash = keccak256(bytes(mspId));
        if (_orgs[mspIdHash].registered) revert OrgAlreadyRegistered(mspIdHash);
        _orgs[mspIdHash] = Org({registered: true, active: true, keyCount: 0, mspId: mspId});
        _orgList.push(mspIdHash);
        activeOrgCount += 1;
        emit OrgRegistered(mspIdHash, mspId);
    }

    function setOrgActive(bytes32 mspIdHash, bool active) external onlyRole(GOVERNOR_ROLE) {
        Org storage o = _orgs[mspIdHash];
        if (!o.registered) revert UnknownOrg(mspIdHash);
        if (o.active == active) return;
        o.active = active;
        activeOrgCount = active ? activeOrgCount + 1 : activeOrgCount - 1;
        if (threshold > activeOrgCount) revert InvalidThreshold(threshold, activeOrgCount);
        emit OrgStatusChanged(mspIdHash, active);
    }

    function addAttestorKey(bytes32 mspIdHash, bytes32 qx, bytes32 qy, bytes32 certHash)
        external
        onlyRole(GOVERNOR_ROLE)
        returns (uint32 keyIndex)
    {
        Org storage o = _orgs[mspIdHash];
        if (!o.registered) revert UnknownOrg(mspIdHash);
        if (!P256.isValidPublicKey(qx, qy)) revert InvalidPublicKey();
        keyIndex = o.keyCount++;
        _keys[mspIdHash][keyIndex] = AttestorKey({qx: qx, qy: qy, certHash: certHash, active: true});
        emit AttestorKeyAdded(mspIdHash, keyIndex, qx, qy, certHash);
    }

    function revokeAttestorKey(bytes32 mspIdHash, uint32 keyIndex) external onlyRole(GOVERNOR_ROLE) {
        AttestorKey storage k = _keys[mspIdHash][keyIndex];
        if (k.qx == bytes32(0)) revert UnknownKey(mspIdHash, keyIndex);
        k.active = false;
        emit AttestorKeyRevoked(mspIdHash, keyIndex);
    }

    function setThreshold(uint32 newThreshold) external onlyRole(GOVERNOR_ROLE) {
        if (newThreshold == 0 || newThreshold > activeOrgCount) revert InvalidThreshold(newThreshold, activeOrgCount);
        threshold = newThreshold;
        emit ThresholdChanged(newThreshold);
    }

    // ------------------------------------------------------------------ verificação

    /// @notice Digest SHA-256 assinado pelos atestadores. Vinculado a esta cadeia e a este contrato,
    ///         o que impede replay entre deployments.
    function digest(FabricLockAttestation calldata a) public view returns (bytes32) {
        return sha256(
            abi.encode(
                DOMAIN_TAG,
                block.chainid,
                address(this),
                a.swapId,
                a.hashlock,
                a.assetIdHash,
                a.recipientIdHash,
                a.deadline,
                a.channelHash,
                a.namespaceHash,
                a.fabricTxId,
                a.blockNumber
            )
        );
    }

    /// @notice Reverte se a atestação não cumprir a política.
    /// @dev As assinaturas precisam vir em ordem estritamente crescente de `mspIdHash`. Isso garante
    ///      no máximo uma assinatura por org sem precisar de storage temporário.
    function requireValid(FabricLockAttestation calldata a, OrgSignature[] calldata sigs) public view {
        uint32 t = threshold;
        if (t == 0) revert InvalidThreshold(0, activeOrgCount);
        bytes32 h = digest(a);
        bytes32 prev;
        uint32 valid;
        for (uint256 i = 0; i < sigs.length; ++i) {
            OrgSignature calldata sig = sigs[i];
            if (i > 0 && sig.mspIdHash <= prev) revert SignaturesNotSorted();
            prev = sig.mspIdHash;

            Org storage o = _orgs[sig.mspIdHash];
            AttestorKey storage k = _keys[sig.mspIdHash][sig.keyIndex];
            if (!o.active || !k.active) revert InactiveSigner(sig.mspIdHash, sig.keyIndex);
            if (!P256.verify(h, sig.r, sig.s, k.qx, k.qy)) revert InvalidSignature(sig.mspIdHash, sig.keyIndex);
            unchecked {
                ++valid;
            }
        }
        if (valid < t) revert ThresholdNotMet(valid, t);
    }

    function isValid(FabricLockAttestation calldata a, OrgSignature[] calldata sigs) external view returns (bool) {
        try this.requireValid(a, sigs) {
            return true;
        } catch {
            return false;
        }
    }

    // ------------------------------------------------------------------ leitura

    function getOrg(bytes32 mspIdHash) external view returns (Org memory) {
        return _orgs[mspIdHash];
    }

    function getAttestorKey(bytes32 mspIdHash, uint32 keyIndex) external view returns (AttestorKey memory) {
        return _keys[mspIdHash][keyIndex];
    }

    function orgs() external view returns (bytes32[] memory) {
        return _orgList;
    }
}
