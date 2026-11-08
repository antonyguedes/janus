// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.30;

import {Base} from "./Base.t.sol";
import {HTLCPayment} from "../src/HTLCPayment.sol";

contract HTLCPaymentTest is Base {
    HTLCPayment internal htlc;

    bytes32 internal constant SWAP = keccak256("swap-1");
    bytes32 internal constant S = keccak256("segredo da vendedora");
    bytes32 internal H;
    uint256 internal constant AMOUNT = 1_250_000_00; // R$ 1,25 mi
    uint64 internal T2;

    event PaymentClaimed(bytes32 indexed swapId, bytes32 preimage, address indexed recipient, uint256 amount);

    function setUp() public {
        _deployDrex();
        htlc = new HTLCPayment(drex);
        _allow(address(htlc));
        H = sha256(abi.encodePacked(S));
        T2 = uint64(block.timestamp + 2 hours);
        vm.prank(buyer);
        drex.approve(address(htlc), type(uint256).max);
    }

    function _lock() internal {
        vm.prank(buyer);
        htlc.lock(SWAP, seller, AMOUNT, H, T2);
    }

    // ------------------------------------------------------------------ caminho feliz

    function test_lockAndClaim() public {
        _lock();
        assertEq(drex.balanceOf(address(htlc)), AMOUNT);
        assertEq(htlc.totalLocked(), AMOUNT);

        vm.expectEmit(true, true, false, true);
        emit PaymentClaimed(SWAP, S, seller, AMOUNT);
        vm.prank(stranger); // permissionless
        htlc.claim(SWAP, S);

        assertEq(drex.balanceOf(seller), AMOUNT);
        assertEq(htlc.revealedPreimage(SWAP), S);
        assertEq(uint8(htlc.getLock(SWAP).state), uint8(HTLCPayment.State.Claimed));
        assertEq(htlc.totalLocked(), 0);
    }

    function test_refundAfterT2() public {
        _lock();
        vm.warp(T2);
        vm.prank(stranger); // permissionless
        htlc.refund(SWAP);
        assertEq(drex.balanceOf(buyer), 10_000_000_00);
        assertEq(uint8(htlc.getLock(SWAP).state), uint8(HTLCPayment.State.Refunded));
    }

    // ------------------------------------------------------------------ fronteiras de tempo

    function test_claimAtLastSecond() public {
        _lock();
        vm.warp(T2 - 1);
        htlc.claim(SWAP, S);
        assertEq(drex.balanceOf(seller), AMOUNT);
    }

    function test_revert_claimAtT2() public {
        _lock();
        vm.warp(T2);
        vm.expectRevert(abi.encodeWithSelector(HTLCPayment.Expired.selector, T2, T2));
        htlc.claim(SWAP, S);
    }

    function test_revert_refundBeforeT2() public {
        _lock();
        vm.warp(T2 - 1);
        vm.expectRevert(abi.encodeWithSelector(HTLCPayment.NotExpired.selector, T2, T2 - 1));
        htlc.refund(SWAP);
    }

    /// @dev Fuzz: em qualquer instante, exatamente uma das duas saídas é possível.
    function testFuzz_timeBoundaryIsExclusive(uint64 dt) public {
        _lock();
        uint256 t = block.timestamp + bound(dt, 0, 4 hours);
        vm.warp(t);
        uint256 snap = vm.snapshotState();
        bool claimOk;
        bool refundOk;
        try htlc.claim(SWAP, S) {
            claimOk = true;
        } catch {}
        vm.revertToState(snap);
        try htlc.refund(SWAP) {
            refundOk = true;
        } catch {}
        assertTrue(claimOk != refundOk, "claim e refund nao podem ser ambos possiveis/impossiveis");
        assertEq(claimOk, t < T2);
    }

    // ------------------------------------------------------------------ preimagem / hashlock

    function test_revert_wrongPreimage() public {
        _lock();
        vm.expectRevert(abi.encodeWithSelector(HTLCPayment.WrongPreimage.selector, SWAP));
        htlc.claim(SWAP, keccak256("errado"));
    }

    function testFuzz_revert_wrongPreimage(bytes32 guess) public {
        vm.assume(guess != S);
        _lock();
        vm.expectRevert(abi.encodeWithSelector(HTLCPayment.WrongPreimage.selector, SWAP));
        htlc.claim(SWAP, guess);
    }

    function test_revert_hashlockReuse() public {
        _lock();
        vm.prank(buyer);
        vm.expectRevert(abi.encodeWithSelector(HTLCPayment.HashlockReused.selector, H));
        htlc.lock(keccak256("swap-2"), seller, 1, H, T2);
    }

    function test_revert_hashlockReuseAfterClaim() public {
        _lock();
        htlc.claim(SWAP, S);
        vm.prank(buyer);
        vm.expectRevert(abi.encodeWithSelector(HTLCPayment.HashlockReused.selector, H));
        htlc.lock(keccak256("swap-2"), seller, 1, H, T2);
    }

    // ------------------------------------------------------------------ estados

    function test_revert_doubleClaim() public {
        _lock();
        htlc.claim(SWAP, S);
        vm.expectRevert(abi.encodeWithSelector(HTLCPayment.NotLocked.selector, SWAP, HTLCPayment.State.Claimed));
        htlc.claim(SWAP, S);
    }

    function test_revert_refundAfterClaim() public {
        _lock();
        htlc.claim(SWAP, S);
        vm.warp(T2);
        vm.expectRevert(abi.encodeWithSelector(HTLCPayment.NotLocked.selector, SWAP, HTLCPayment.State.Claimed));
        htlc.refund(SWAP);
    }

    function test_revert_claimAfterRefund() public {
        _lock();
        vm.warp(T2);
        htlc.refund(SWAP);
        vm.expectRevert(abi.encodeWithSelector(HTLCPayment.NotLocked.selector, SWAP, HTLCPayment.State.Refunded));
        htlc.claim(SWAP, S);
    }

    function test_revert_swapIdReuse() public {
        _lock();
        vm.prank(buyer);
        vm.expectRevert(abi.encodeWithSelector(HTLCPayment.SwapExists.selector, SWAP));
        htlc.lock(SWAP, seller, 1, keccak256("h2"), T2);
    }

    function test_revert_unknownSwap() public {
        vm.expectRevert(abi.encodeWithSelector(HTLCPayment.NotLocked.selector, SWAP, HTLCPayment.State.None));
        htlc.claim(SWAP, S);
    }

    // ------------------------------------------------------------------ parâmetros

    function test_revert_invalidParams() public {
        vm.startPrank(buyer);
        vm.expectRevert(HTLCPayment.InvalidParams.selector);
        htlc.lock(bytes32(0), seller, AMOUNT, H, T2);
        vm.expectRevert(HTLCPayment.InvalidParams.selector);
        htlc.lock(SWAP, address(0), AMOUNT, H, T2);
        vm.expectRevert(HTLCPayment.InvalidParams.selector);
        htlc.lock(SWAP, seller, 0, H, T2);
        vm.expectRevert(HTLCPayment.InvalidParams.selector);
        htlc.lock(SWAP, seller, AMOUNT, bytes32(0), T2);
        vm.expectRevert(HTLCPayment.InvalidParams.selector);
        htlc.lock(SWAP, buyer, AMOUNT, H, T2);
        vm.expectRevert(
            abi.encodeWithSelector(HTLCPayment.TimelockInPast.selector, uint64(block.timestamp), block.timestamp)
        );
        htlc.lock(SWAP, seller, AMOUNT, H, uint64(block.timestamp));
        vm.stopPrank();
    }

    function test_revert_contractNotAllowlisted() public {
        HTLCPayment other = new HTLCPayment(drex);
        vm.startPrank(buyer);
        drex.approve(address(other), AMOUNT);
        vm.expectRevert();
        other.lock(SWAP, seller, AMOUNT, H, T2);
        vm.stopPrank();
    }

    function test_sha256MatchesGoConvention() public pure {
        // sha256 de 32 bytes zero (vetor clássico)
        assertEq(
            sha256(abi.encodePacked(bytes32(0))),
            0x66687aadf862bd776c8fc18b8e9f8e20089714856ee233b3902a591d0d5f2925
        );
    }
}
