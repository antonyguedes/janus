// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.30;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

/// @title HTLCPayment — perna de pagamento da Abordagem A (HTLC)
/// @notice O comprador (Banco) trava Drex sob `hashlock = sha256(preimage)` até `timelock` (T2).
///         - `claim(swapId, preimage)` antes de T2 envia ao vendedor e **publica a preimagem** no evento.
///         - `refund(swapId)` a partir de T2 devolve ao comprador.
///         As duas funções são *permissionless*: qualquer um pode chamá-las, e os fundos sempre vão
///         para o destino fixado no lock.
/// @dev A fronteira de tempo é idêntica à do validador HTLC do Panurus: claim se `now < T`,
///      reversão se `now >= T` (ver ADR-0004). O hash usado é SHA-256 sobre exatamente 32 bytes,
///      compatível com `crypto/sha256` em Go (ver testdata/vectors/htlc.json).
contract HTLCPayment is ReentrancyGuard {
    using SafeERC20 for IERC20;

    enum State {
        None,
        Locked,
        Claimed,
        Refunded
    }

    struct Lock {
        address sender; // comprador (Banco)
        address recipient; // vendedor (Vendedora)
        uint256 amount;
        bytes32 hashlock;
        uint64 timelock; // T2 (unix s)
        State state;
    }

    IERC20 public immutable token;

    mapping(bytes32 swapId => Lock) private _locks;
    /// @notice Bloqueio de reuso de hashlock: uma preimagem revelada nunca pode destravar um segundo swap.
    mapping(bytes32 hashlock => bool) public hashlockUsed;
    mapping(bytes32 swapId => bytes32) public revealedPreimage;

    /// @notice Soma dos valores em locks ativos (usada nas invariantes).
    uint256 public totalLocked;

    event PaymentLocked(
        bytes32 indexed swapId,
        address indexed sender,
        address indexed recipient,
        uint256 amount,
        bytes32 hashlock,
        uint64 timelock
    );
    event PaymentClaimed(bytes32 indexed swapId, bytes32 preimage, address indexed recipient, uint256 amount);
    event PaymentRefunded(bytes32 indexed swapId, address indexed sender, uint256 amount);

    error SwapExists(bytes32 swapId);
    error HashlockReused(bytes32 hashlock);
    error InvalidParams();
    error TimelockInPast(uint64 timelock, uint256 nowTs);
    error NotLocked(bytes32 swapId, State state);
    error Expired(uint64 timelock, uint256 nowTs);
    error NotExpired(uint64 timelock, uint256 nowTs);
    error WrongPreimage(bytes32 swapId);
    error UnexpectedBalanceDelta(uint256 expected, uint256 received);

    constructor(IERC20 token_) {
        if (address(token_) == address(0)) revert InvalidParams();
        token = token_;
    }

    function lock(bytes32 swapId, address recipient, uint256 amount, bytes32 hashlock, uint64 timelock)
        external
        nonReentrant
    {
        if (swapId == bytes32(0) || hashlock == bytes32(0) || recipient == address(0) || amount == 0) {
            revert InvalidParams();
        }
        if (recipient == msg.sender) revert InvalidParams();
        if (_locks[swapId].state != State.None) revert SwapExists(swapId);
        if (hashlockUsed[hashlock]) revert HashlockReused(hashlock);
        if (timelock <= block.timestamp) revert TimelockInPast(timelock, block.timestamp);

        // efeitos antes da interação
        _locks[swapId] = Lock({
            sender: msg.sender,
            recipient: recipient,
            amount: amount,
            hashlock: hashlock,
            timelock: timelock,
            state: State.Locked
        });
        hashlockUsed[hashlock] = true;
        totalLocked += amount;

        uint256 before = token.balanceOf(address(this));
        token.safeTransferFrom(msg.sender, address(this), amount);
        uint256 received = token.balanceOf(address(this)) - before;
        if (received != amount) revert UnexpectedBalanceDelta(amount, received);

        emit PaymentLocked(swapId, msg.sender, recipient, amount, hashlock, timelock);
    }

    /// @notice Permissionless. Revela `preimage` e paga o vendedor. Só vale antes de T2.
    function claim(bytes32 swapId, bytes32 preimage) external nonReentrant {
        Lock storage l = _locks[swapId];
        if (l.state != State.Locked) revert NotLocked(swapId, l.state);
        if (block.timestamp >= l.timelock) revert Expired(l.timelock, block.timestamp);
        if (sha256(abi.encodePacked(preimage)) != l.hashlock) revert WrongPreimage(swapId);

        l.state = State.Claimed;
        revealedPreimage[swapId] = preimage;
        uint256 amount = l.amount;
        totalLocked -= amount;

        token.safeTransfer(l.recipient, amount);
        emit PaymentClaimed(swapId, preimage, l.recipient, amount);
    }

    /// @notice Permissionless. Devolve ao comprador a partir de T2.
    function refund(bytes32 swapId) external nonReentrant {
        Lock storage l = _locks[swapId];
        if (l.state != State.Locked) revert NotLocked(swapId, l.state);
        if (block.timestamp < l.timelock) revert NotExpired(l.timelock, block.timestamp);

        l.state = State.Refunded;
        uint256 amount = l.amount;
        totalLocked -= amount;

        token.safeTransfer(l.sender, amount);
        emit PaymentRefunded(swapId, l.sender, amount);
    }

    function getLock(bytes32 swapId) external view returns (Lock memory) {
        return _locks[swapId];
    }
}
