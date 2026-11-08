// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.30;

import {Base} from "./Base.t.sol";
import {DrexToken} from "../src/DrexToken.sol";
import {IAccessControl} from "@openzeppelin/contracts/access/IAccessControl.sol";

contract DrexTokenTest is Base {
    function setUp() public {
        _deployDrex();
    }

    function test_decimals() public view {
        assertEq(drex.decimals(), 2);
    }

    function test_transferBetweenAllowed() public {
        vm.prank(buyer);
        drex.transfer(seller, 100);
        assertEq(drex.balanceOf(seller), 100);
    }

    function test_revert_transferToNotAllowed() public {
        address outsider = makeAddr("fora");
        vm.prank(buyer);
        vm.expectRevert(abi.encodeWithSelector(DrexToken.NotAllowlisted.selector, outsider));
        drex.transfer(outsider, 1);
    }

    function test_revert_transferFromRevokedSender() public {
        vm.prank(admin);
        drex.setAllowed(buyer, false);
        vm.prank(buyer);
        vm.expectRevert(abi.encodeWithSelector(DrexToken.NotAllowlisted.selector, buyer));
        drex.transfer(seller, 1);
    }

    function test_revert_mintToNotAllowed() public {
        address outsider = makeAddr("fora");
        vm.prank(admin);
        vm.expectRevert(abi.encodeWithSelector(DrexToken.NotAllowlisted.selector, outsider));
        drex.mint(outsider, 1);
    }

    function test_revert_onlyIssuerMints() public {
        bytes32 role = drex.ISSUER_ROLE();
        vm.prank(buyer);
        vm.expectRevert(abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, buyer, role));
        drex.mint(buyer, 1);
    }

    function test_revert_onlyComplianceAllowlists() public {
        vm.prank(buyer);
        vm.expectRevert();
        drex.setAllowed(stranger, true);
    }

    function test_burn() public {
        vm.prank(admin);
        drex.burn(buyer, 50);
        assertEq(drex.balanceOf(buyer), 10_000_000_00 - 50);
    }
}
