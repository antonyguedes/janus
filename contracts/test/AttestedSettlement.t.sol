// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.30;

import {Base} from "./Base.t.sol";
import {FabricProofVerifier, FabricLockAttestation, OrgSignature} from "../src/FabricProofVerifier.sol";
import {AttestedSettlement} from "../src/AttestedSettlement.sol";

contract AttestedSettlementTest is Base {
    FabricProofVerifier internal verifier;
    AttestedSettlement internal settlement;

    uint256 internal pkReg = uint256(keccak256("registradora-attestor")) % P256_N;
    uint256 internal pkBanco = uint256(keccak256("banco-attestor")) % P256_N;
    uint256 internal pkVend = uint256(keccak256("vendedora-attestor")) % P256_N;
    bytes32 internal mspReg;
    bytes32 internal mspBanco;
    bytes32 internal mspVend;

    bytes32 internal constant SWAP = keccak256("swap-b");
    bytes32 internal constant S = keccak256("segredo-b");
    bytes32 internal H;
    bytes32 internal ASSET = sha256("recebivel-0001");
    bytes32 internal RECIPIENT = sha256("BancoMSP:banco-owner");
    bytes32 internal CHANNEL = sha256("recebiveis");
    bytes32 internal NAMESPACE = sha256("tcc");
    uint256 internal constant AMOUNT = 500_000_00;
    uint64 internal T2;
    uint64 internal T1;

    function setUp() public {
        _deployDrex();
        verifier = new FabricProofVerifier(admin);
        settlement = new AttestedSettlement(drex, verifier, CHANNEL, NAMESPACE);
        _allow(address(settlement));

        vm.startPrank(admin);
        mspReg = _registerOrg("RegistradoraMSP", pkReg);
        mspBanco = _registerOrg("BancoMSP", pkBanco);
        mspVend = _registerOrg("VendedoraMSP", pkVend);
        // política 2-de-n, com a Vendedora desativada como atestadora (ADR-0003)
        verifier.setOrgActive(mspVend, false);
        verifier.setThreshold(2);
        vm.stopPrank();

        H = sha256(abi.encodePacked(S));
        T2 = uint64(block.timestamp + 2 hours);
        T1 = uint64(block.timestamp + 4 hours);

        vm.prank(buyer);
        drex.approve(address(settlement), type(uint256).max);
    }

    function _registerOrg(string memory msp, uint256 pk) internal returns (bytes32 h) {
        h = verifier.registerOrg(msp);
        (uint256 x, uint256 y) = vm.publicKeyP256(pk);
        verifier.addAttestorKey(h, bytes32(x), bytes32(y), keccak256(abi.encode("cert", msp)));
    }

    function _deposit() internal {
        vm.prank(buyer);
        settlement.deposit(SWAP, seller, AMOUNT, H, ASSET, RECIPIENT, T2, T1 - 30);
    }

    function _att() internal view returns (FabricLockAttestation memory) {
        return FabricLockAttestation({
            swapId: SWAP,
            hashlock: H,
            assetIdHash: ASSET,
            recipientIdHash: RECIPIENT,
            deadline: T1,
            channelHash: CHANNEL,
            namespaceHash: NAMESPACE,
            fabricTxId: keccak256("txid"),
            blockNumber: 7
        });
    }

    function _sig(bytes32 msp, uint256 pk, FabricLockAttestation memory a)
        internal
        view
        returns (OrgSignature memory)
    {
        (bytes32 r, bytes32 s) = _signP256(pk, verifier.digest(a));
        return OrgSignature({mspIdHash: msp, keyIndex: 0, r: r, s: s});
    }

    function _sorted2(OrgSignature memory x, OrgSignature memory y) internal pure returns (OrgSignature[] memory o) {
        o = new OrgSignature[](2);
        (o[0], o[1]) = x.mspIdHash < y.mspIdHash ? (x, y) : (y, x);
    }

    function _goodSigs(FabricLockAttestation memory a) internal view returns (OrgSignature[] memory) {
        return _sorted2(_sig(mspReg, pkReg, a), _sig(mspBanco, pkBanco, a));
    }

    // ------------------------------------------------------------------ caminho feliz

    function test_depositAndSettle() public {
        _deposit();
        FabricLockAttestation memory a = _att();
        vm.prank(stranger); // relay / qualquer um
        settlement.settle(SWAP, S, a, _goodSigs(a));
        assertEq(drex.balanceOf(seller), AMOUNT);
        assertEq(settlement.revealedPreimage(SWAP), S);
        assertEq(uint8(settlement.getDeposit(SWAP).state), uint8(AttestedSettlement.State.Settled));
        assertEq(settlement.totalDeposited(), 0);
    }

    function test_refundAfterT2() public {
        _deposit();
        vm.warp(T2);
        vm.prank(stranger);
        settlement.refund(SWAP);
        assertEq(drex.balanceOf(buyer), 10_000_000_00);
    }

    function test_revert_refundBeforeT2() public {
        _deposit();
        vm.warp(T2 - 1);
        vm.expectRevert(abi.encodeWithSelector(AttestedSettlement.NotExpired.selector, T2, T2 - 1));
        settlement.refund(SWAP);
    }

    // ------------------------------------------------------------------ prova adulterada (cenário 7)

    function test_revert_insufficientEndorsers() public {
        _deposit();
        FabricLockAttestation memory a = _att();
        OrgSignature[] memory one = new OrgSignature[](1);
        one[0] = _sig(mspReg, pkReg, a);
        vm.expectRevert(abi.encodeWithSelector(FabricProofVerifier.ThresholdNotMet.selector, 1, 2));
        settlement.settle(SWAP, S, a, one);
    }

    function test_revert_duplicateOrgSignature() public {
        _deposit();
        FabricLockAttestation memory a = _att();
        OrgSignature[] memory dup = new OrgSignature[](2);
        dup[0] = _sig(mspReg, pkReg, a);
        dup[1] = dup[0];
        vm.expectRevert(FabricProofVerifier.SignaturesNotSorted.selector);
        settlement.settle(SWAP, S, a, dup);
    }

    function test_revert_inactiveOrgCannotAttest() public {
        _deposit();
        FabricLockAttestation memory a = _att();
        OrgSignature[] memory sigs = _sorted2(_sig(mspReg, pkReg, a), _sig(mspVend, pkVend, a));
        vm.expectRevert(abi.encodeWithSelector(FabricProofVerifier.InactiveSigner.selector, mspVend, 0));
        settlement.settle(SWAP, S, a, sigs);
    }

    function test_revert_forgedSignature() public {
        _deposit();
        FabricLockAttestation memory a = _att();
        // a Vendedora assina com a própria chave, fingindo ser o Banco
        OrgSignature memory forged = _sig(mspBanco, pkVend, a);
        OrgSignature[] memory sigs = _sorted2(_sig(mspReg, pkReg, a), forged);
        vm.expectRevert(abi.encodeWithSelector(FabricProofVerifier.InvalidSignature.selector, mspBanco, 0));
        settlement.settle(SWAP, S, a, sigs);
    }

    function test_revert_highS() public {
        _deposit();
        FabricLockAttestation memory a = _att();
        OrgSignature memory x = _sig(mspReg, pkReg, a);
        x.s = bytes32(P256_N - uint256(x.s)); // forma maleável
        OrgSignature[] memory sigs = _sorted2(x, _sig(mspBanco, pkBanco, a));
        vm.expectRevert(abi.encodeWithSelector(FabricProofVerifier.InvalidSignature.selector, mspReg, 0));
        settlement.settle(SWAP, S, a, sigs);
    }

    function test_revert_tamperedAfterSigning() public {
        _deposit();
        FabricLockAttestation memory a = _att();
        OrgSignature[] memory sigs = _goodSigs(a);
        a.blockNumber = 8;
        vm.expectRevert();
        settlement.settle(SWAP, S, a, sigs);
    }

    function test_revert_revokedKey() public {
        _deposit();
        vm.prank(admin);
        verifier.revokeAttestorKey(mspBanco, 0);
        FabricLockAttestation memory a = _att();
        OrgSignature[] memory sigs = _goodSigs(a);
        vm.expectRevert(abi.encodeWithSelector(FabricProofVerifier.InactiveSigner.selector, mspBanco, 0));
        settlement.settle(SWAP, S, a, sigs);
    }

    function test_digestBoundToDeployment() public {
        FabricLockAttestation memory a = _att();
        FabricProofVerifier other = new FabricProofVerifier(admin);
        assertTrue(other.digest(a) != verifier.digest(a), "digest precisa depender do endereco");
        uint256 snap = vm.snapshotState();
        bytes32 d1 = verifier.digest(a);
        vm.chainId(999);
        assertTrue(verifier.digest(a) != d1, "digest precisa depender do chainId");
        vm.revertToState(snap);
    }

    // ------------------------------------------------------------------ termos divergentes

    function _expectMismatch(FabricLockAttestation memory a, string memory field) internal {
        OrgSignature[] memory sigs = _goodSigs(a);
        vm.expectRevert(abi.encodeWithSelector(AttestedSettlement.AttestationMismatch.selector, field));
        settlement.settle(SWAP, S, a, sigs);
    }

    function test_revert_attestationMismatches() public {
        _deposit();
        FabricLockAttestation memory a;

        a = _att();
        a.swapId = keccak256("outro-swap");
        _expectMismatch(a, "swapId");

        a = _att();
        a.assetIdHash = sha256("outro-recebivel");
        _expectMismatch(a, "assetId");

        a = _att();
        a.recipientIdHash = sha256("VendedoraMSP:laranja");
        _expectMismatch(a, "recipient");

        a = _att();
        a.deadline = T1 - 31; // abaixo do mínimo aceito
        _expectMismatch(a, "deadline");

        a = _att();
        a.channelHash = sha256("outro-canal");
        _expectMismatch(a, "channel");

        a = _att();
        a.namespaceHash = sha256("outro-cc");
        _expectMismatch(a, "namespace");
    }

    function test_revert_wrongPreimage() public {
        _deposit();
        FabricLockAttestation memory a = _att();
        OrgSignature[] memory sigs = _goodSigs(a);
        vm.expectRevert(abi.encodeWithSelector(AttestedSettlement.WrongPreimage.selector, SWAP));
        settlement.settle(SWAP, keccak256("x"), a, sigs);
    }

    function test_revert_settleAtT2() public {
        _deposit();
        FabricLockAttestation memory a = _att();
        OrgSignature[] memory sigs = _goodSigs(a);
        vm.warp(T2);
        vm.expectRevert(abi.encodeWithSelector(AttestedSettlement.Expired.selector, T2, T2));
        settlement.settle(SWAP, S, a, sigs);
    }

    function test_revert_depositRequiresT1AfterT2() public {
        vm.prank(buyer);
        vm.expectRevert(AttestedSettlement.InvalidParams.selector);
        settlement.deposit(SWAP, seller, AMOUNT, H, ASSET, RECIPIENT, T2, T2);
    }

    function test_revert_hashlockReuse() public {
        _deposit();
        vm.prank(buyer);
        vm.expectRevert(abi.encodeWithSelector(AttestedSettlement.HashlockReused.selector, H));
        settlement.deposit(keccak256("swap-c"), seller, 1, H, ASSET, RECIPIENT, T2, T1);
    }

    function test_revert_doubleSettle() public {
        _deposit();
        FabricLockAttestation memory a = _att();
        OrgSignature[] memory sigs = _goodSigs(a);
        settlement.settle(SWAP, S, a, sigs);
        vm.expectRevert(
            abi.encodeWithSelector(AttestedSettlement.NotDeposited.selector, SWAP, AttestedSettlement.State.Settled)
        );
        settlement.settle(SWAP, S, a, sigs);
    }

    // ------------------------------------------------------------------ governança

    function test_revert_thresholdAboveActiveOrgs() public {
        vm.prank(admin);
        vm.expectRevert(abi.encodeWithSelector(FabricProofVerifier.InvalidThreshold.selector, 3, 2));
        verifier.setThreshold(3);
    }

    function test_revert_deactivatingBelowThreshold() public {
        vm.prank(admin);
        vm.expectRevert(abi.encodeWithSelector(FabricProofVerifier.InvalidThreshold.selector, 2, 1));
        verifier.setOrgActive(mspBanco, false);
    }

    function test_revert_invalidPublicKey() public {
        vm.prank(admin);
        vm.expectRevert(FabricProofVerifier.InvalidPublicKey.selector);
        verifier.addAttestorKey(mspReg, bytes32(uint256(1)), bytes32(uint256(2)), bytes32(0));
    }

    function test_revert_nonGovernor() public {
        vm.prank(seller);
        vm.expectRevert();
        verifier.setThreshold(1);
    }

    function test_usesNativePrecompile() public view {
        // O Foundry (EVM osaka) expõe P256VERIFY em 0x100, como o Besu 26.11 com osakaTime=0.
        FabricLockAttestation memory a = _att();
        OrgSignature memory x = _sig(mspReg, pkReg, a);
        (uint256 qx, uint256 qy) = vm.publicKeyP256(pkReg);
        (bool ok, bytes memory ret) =
            address(0x100).staticcall(abi.encode(verifier.digest(a), x.r, x.s, qx, qy));
        assertTrue(ok);
        assertEq(abi.decode(ret, (uint256)), 1);
    }
}
