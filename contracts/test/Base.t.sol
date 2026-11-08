// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.30;

import {Test} from "forge-std/Test.sol";
import {DrexToken} from "../src/DrexToken.sol";

abstract contract Base is Test {
    uint256 internal constant P256_N = 0xFFFFFFFF00000000FFFFFFFFFFFFFFFFBCE6FAADA7179E84F3B9CAC2FC632551;

    address internal admin = makeAddr("bcb");
    address internal seller = makeAddr("vendedora");
    address internal buyer = makeAddr("banco");
    address internal stranger = makeAddr("terceiro");

    DrexToken internal drex;

    function _deployDrex() internal {
        drex = new DrexToken(admin);
        vm.startPrank(admin);
        drex.setAllowed(seller, true);
        drex.setAllowed(buyer, true);
        drex.setAllowed(stranger, true);
        drex.mint(buyer, 10_000_000_00); // R$ 10 mi
        vm.stopPrank();
    }

    function _allow(address a) internal {
        vm.prank(admin);
        drex.setAllowed(a, true);
    }

    /// @dev Assina com P-256 e normaliza para low-s (exigido pelo OZ P256 e pelo Fabric).
    function _signP256(uint256 pk, bytes32 digest) internal pure returns (bytes32 r, bytes32 s) {
        (r, s) = vm.signP256(pk, digest);
        if (uint256(s) > P256_N / 2) s = bytes32(P256_N - uint256(s));
    }
}
