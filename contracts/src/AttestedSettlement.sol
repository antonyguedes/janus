// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.30;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import {FabricProofVerifier, FabricLockAttestation, OrgSignature} from "./FabricProofVerifier.sol";

/// @title AttestedSettlement — perna de pagamento da Abordagem B2 ("HTLC atestado", ADR-0003)
/// @notice O Banco deposita Drex para um swap. O pagamento só é liberado ao vendedor com:
///           (i)  a preimagem `S`, com `sha256(S) == hashlock`; e
///           (ii) uma atestação do lock do NFT no Fabric que cumpra a política do FabricProofVerifier,
///                com os termos do lock (ativo, recipient, hashlock, deadline) batendo com o depósito.
///         Revelar `S` aqui permite ao Banco fazer o `htlc.Claim(S)` do Panurus no Fabric.
///         Depois de T2, o depósito é reversível (`refund`, permissionless).
contract AttestedSettlement is ReentrancyGuard {
    using SafeERC20 for IERC20;

    enum State {
        None,
        Deposited,
        Settled,
        Refunded
    }

    struct Deposit {
        address buyer;
        address seller;
        uint256 amount;
        bytes32 hashlock;
        bytes32 assetIdHash;
        bytes32 recipientIdHash;
        uint64 timelock; // T2
        uint64 minFabricDeadline; // menor T1 aceitável no lock do Fabric (T1 - tolerância)
        State state;
    }

    IERC20 public immutable token;
    FabricProofVerifier public immutable verifier;
    bytes32 public immutable channelHash;
    bytes32 public immutable namespaceHash;

    mapping(bytes32 swapId => Deposit) private _deposits;
    mapping(bytes32 hashlock => bool) public hashlockUsed;
    mapping(bytes32 swapId => bytes32) public revealedPreimage;
    uint256 public totalDeposited;

    event PaymentDeposited(
        bytes32 indexed swapId,
        address indexed buyer,
        address indexed seller,
        uint256 amount,
        bytes32 hashlock,
        uint64 timelock,
        uint64 minFabricDeadline
    );
    event PaymentSettled(
        bytes32 indexed swapId, bytes32 preimage, bytes32 fabricTxId, address indexed seller, uint256 amount
    );
    event PaymentRefunded(bytes32 indexed swapId, address indexed buyer, uint256 amount);

    error InvalidParams();
    error SwapExists(bytes32 swapId);
    error HashlockReused(bytes32 hashlock);
    error TimelockInPast(uint64 timelock, uint256 nowTs);
    error NotDeposited(bytes32 swapId, State state);
    error Expired(uint64 timelock, uint256 nowTs);
    error NotExpired(uint64 timelock, uint256 nowTs);
    error WrongPreimage(bytes32 swapId);
    error AttestationMismatch(string field);
    error UnexpectedBalanceDelta(uint256 expected, uint256 received);

    constructor(IERC20 token_, FabricProofVerifier verifier_, bytes32 channelHash_, bytes32 namespaceHash_) {
        if (address(token_) == address(0) || address(verifier_) == address(0)) revert InvalidParams();
        token = token_;
        verifier = verifier_;
        channelHash = channelHash_;
        namespaceHash = namespaceHash_;
    }

    function deposit(
        bytes32 swapId,
        address seller,
        uint256 amount,
        bytes32 hashlock,
        bytes32 assetIdHash,
        bytes32 recipientIdHash,
        uint64 timelock,
        uint64 minFabricDeadline
    ) external nonReentrant {
        if (
            swapId == bytes32(0) || hashlock == bytes32(0) || seller == address(0) || seller == msg.sender
                || amount == 0 || assetIdHash == bytes32(0) || recipientIdHash == bytes32(0)
        ) revert InvalidParams();
        if (minFabricDeadline <= timelock) revert InvalidParams(); // T1 > T2 obrigatório
        if (_deposits[swapId].state != State.None) revert SwapExists(swapId);
        if (hashlockUsed[hashlock]) revert HashlockReused(hashlock);
        if (timelock <= block.timestamp) revert TimelockInPast(timelock, block.timestamp);

        _deposits[swapId] = Deposit({
            buyer: msg.sender,
            seller: seller,
            amount: amount,
            hashlock: hashlock,
            assetIdHash: assetIdHash,
            recipientIdHash: recipientIdHash,
            timelock: timelock,
            minFabricDeadline: minFabricDeadline,
            state: State.Deposited
        });
        hashlockUsed[hashlock] = true;
        totalDeposited += amount;

        uint256 before = token.balanceOf(address(this));
        token.safeTransferFrom(msg.sender, address(this), amount);
        uint256 received = token.balanceOf(address(this)) - before;
        if (received != amount) revert UnexpectedBalanceDelta(amount, received);

        emit PaymentDeposited(swapId, msg.sender, seller, amount, hashlock, timelock, minFabricDeadline);
    }

    /// @notice Permissionless (normalmente chamada pela Vendedora ou pelo relay). Antes de T2.
    function settle(
        bytes32 swapId,
        bytes32 preimage,
        FabricLockAttestation calldata att,
        OrgSignature[] calldata sigs
    ) external nonReentrant {
        Deposit storage d = _deposits[swapId];
        if (d.state != State.Deposited) revert NotDeposited(swapId, d.state);
        if (block.timestamp >= d.timelock) revert Expired(d.timelock, block.timestamp);
        if (sha256(abi.encodePacked(preimage)) != d.hashlock) revert WrongPreimage(swapId);

        if (att.swapId != swapId) revert AttestationMismatch("swapId");
        if (att.hashlock != d.hashlock) revert AttestationMismatch("hashlock");
        if (att.assetIdHash != d.assetIdHash) revert AttestationMismatch("assetId");
        if (att.recipientIdHash != d.recipientIdHash) revert AttestationMismatch("recipient");
        if (att.deadline < d.minFabricDeadline) revert AttestationMismatch("deadline");
        if (att.channelHash != channelHash) revert AttestationMismatch("channel");
        if (att.namespaceHash != namespaceHash) revert AttestationMismatch("namespace");
        verifier.requireValid(att, sigs);

        d.state = State.Settled;
        revealedPreimage[swapId] = preimage;
        uint256 amount = d.amount;
        totalDeposited -= amount;

        token.safeTransfer(d.seller, amount);
        emit PaymentSettled(swapId, preimage, att.fabricTxId, d.seller, amount);
    }

    function refund(bytes32 swapId) external nonReentrant {
        Deposit storage d = _deposits[swapId];
        if (d.state != State.Deposited) revert NotDeposited(swapId, d.state);
        if (block.timestamp < d.timelock) revert NotExpired(d.timelock, block.timestamp);

        d.state = State.Refunded;
        uint256 amount = d.amount;
        totalDeposited -= amount;

        token.safeTransfer(d.buyer, amount);
        emit PaymentRefunded(swapId, d.buyer, amount);
    }

    function getDeposit(bytes32 swapId) external view returns (Deposit memory) {
        return _deposits[swapId];
    }
}
