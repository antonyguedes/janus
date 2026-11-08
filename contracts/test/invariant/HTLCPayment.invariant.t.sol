// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.30;

import {Test} from "forge-std/Test.sol";
import {DrexToken} from "../../src/DrexToken.sol";
import {HTLCPayment} from "../../src/HTLCPayment.sol";

/// @dev Handler: sequências aleatórias de lock/claim/refund/warp por vários atores.
contract HTLCHandler is Test {
    DrexToken public drex;
    HTLCPayment public htlc;
    address[] public actors;

    bytes32[] public swapIds;
    mapping(bytes32 => bytes32) public preimageOf;
    mapping(bytes32 => bool) public everClaimed;
    mapping(bytes32 => bool) public everRefunded;
    uint256 public nonce;

    constructor(DrexToken drex_, HTLCPayment htlc_, address[] memory actors_) {
        drex = drex_;
        htlc = htlc_;
        actors = actors_;
    }

    function _actor(uint256 seed) internal view returns (address) {
        return actors[seed % actors.length];
    }

    function lock(uint256 fromSeed, uint256 toSeed, uint256 amount, uint32 ttl) external {
        address from = _actor(fromSeed);
        address to = _actor(toSeed);
        if (from == to) return;
        amount = bound(amount, 1, drex.balanceOf(from));
        if (amount == 0 || drex.balanceOf(from) == 0) return;
        bytes32 s = keccak256(abi.encode("s", nonce));
        bytes32 id = keccak256(abi.encode("id", nonce));
        ++nonce;
        uint64 t = uint64(block.timestamp + bound(ttl, 1, 6 hours));
        vm.startPrank(from);
        drex.approve(address(htlc), amount);
        htlc.lock(id, to, amount, sha256(abi.encodePacked(s)), t);
        vm.stopPrank();
        swapIds.push(id);
        preimageOf[id] = s;
    }

    function claim(uint256 idx, bool wrong) external {
        if (swapIds.length == 0) return;
        bytes32 id = swapIds[idx % swapIds.length];
        bytes32 s = wrong ? keccak256("errado") : preimageOf[id];
        try htlc.claim(id, s) {
            everClaimed[id] = true;
        } catch {}
    }

    function refund(uint256 idx) external {
        if (swapIds.length == 0) return;
        bytes32 id = swapIds[idx % swapIds.length];
        try htlc.refund(id) {
            everRefunded[id] = true;
        } catch {}
    }

    function warp(uint32 dt) external {
        vm.warp(block.timestamp + bound(dt, 1, 3 hours));
    }

    function swapCount() external view returns (uint256) {
        return swapIds.length;
    }
}

contract HTLCPaymentInvariantTest is Test {
    DrexToken internal drex;
    HTLCPayment internal htlc;
    HTLCHandler internal handler;
    address[] internal actors;
    uint256 internal constant SUPPLY_EACH = 1_000_000_00;

    function setUp() public {
        address admin = makeAddr("bcb");
        drex = new DrexToken(admin);
        htlc = new HTLCPayment(drex);
        for (uint256 i = 0; i < 4; ++i) actors.push(makeAddr(string.concat("ator", vm.toString(i))));
        vm.startPrank(admin);
        drex.setAllowed(address(htlc), true);
        for (uint256 i = 0; i < actors.length; ++i) {
            drex.setAllowed(actors[i], true);
            drex.mint(actors[i], SUPPLY_EACH);
        }
        vm.stopPrank();
        handler = new HTLCHandler(drex, htlc, actors);
        targetContract(address(handler));
    }

    /// Saldo do contrato == soma dos locks ativos.
    function invariant_balanceEqualsActiveLocks() public view {
        uint256 sum;
        uint256 n = handler.swapCount();
        for (uint256 i = 0; i < n; ++i) {
            HTLCPayment.Lock memory l = htlc.getLock(handler.swapIds(i));
            if (l.state == HTLCPayment.State.Locked) sum += l.amount;
        }
        assertEq(drex.balanceOf(address(htlc)), sum);
        assertEq(htlc.totalLocked(), sum);
    }

    /// Nenhum swap fica claimed e refunded ao mesmo tempo.
    function invariant_claimXorRefund() public view {
        uint256 n = handler.swapCount();
        for (uint256 i = 0; i < n; ++i) {
            bytes32 id = handler.swapIds(i);
            assertFalse(handler.everClaimed(id) && handler.everRefunded(id));
        }
    }

    /// A oferta total se conserva (o HTLC não cria nem destrói Drex).
    function invariant_supplyConserved() public view {
        uint256 total = drex.balanceOf(address(htlc));
        for (uint256 i = 0; i < actors.length; ++i) total += drex.balanceOf(actors[i]);
        assertEq(total, SUPPLY_EACH * actors.length);
        assertEq(drex.totalSupply(), SUPPLY_EACH * actors.length);
    }

    /// Todo swap reclamado publicou exatamente a preimagem cujo sha256 é o hashlock.
    function invariant_revealedPreimageMatches() public view {
        uint256 n = handler.swapCount();
        for (uint256 i = 0; i < n; ++i) {
            bytes32 id = handler.swapIds(i);
            HTLCPayment.Lock memory l = htlc.getLock(id);
            if (l.state == HTLCPayment.State.Claimed) {
                assertEq(sha256(abi.encodePacked(htlc.revealedPreimage(id))), l.hashlock);
            }
        }
    }
}
