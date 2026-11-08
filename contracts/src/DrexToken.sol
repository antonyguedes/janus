// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.30;

import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {AccessControl} from "@openzeppelin/contracts/access/AccessControl.sol";

/// @title DrexToken — simulação local do Drex (ERC-20 permissionado)
/// @notice 2 casas decimais (centavos). Só contas na allowlist podem enviar ou receber.
///         Emissão e queima são restritas ao `ISSUER_ROLE` (simula o BCB).
/// @dev A allowlist do token se soma ao allowlist de contas do Besu (defesa em camadas).
///      Os contratos de liquidação (HTLCPayment, AttestedSettlement) também precisam estar na allowlist.
contract DrexToken is ERC20, AccessControl {
    bytes32 public constant ISSUER_ROLE = keccak256("ISSUER_ROLE");
    bytes32 public constant COMPLIANCE_ROLE = keccak256("COMPLIANCE_ROLE");

    mapping(address account => bool) public isAllowed;

    event AllowlistUpdated(address indexed account, bool allowed);

    error NotAllowlisted(address account);
    error ZeroAddress();

    constructor(address admin) ERC20("Drex (Janus local)", "DREX") {
        if (admin == address(0)) revert ZeroAddress();
        _grantRole(DEFAULT_ADMIN_ROLE, admin);
        _grantRole(ISSUER_ROLE, admin);
        _grantRole(COMPLIANCE_ROLE, admin);
    }

    function decimals() public pure override returns (uint8) {
        return 2;
    }

    function setAllowed(address account, bool allowed) external onlyRole(COMPLIANCE_ROLE) {
        if (account == address(0)) revert ZeroAddress();
        isAllowed[account] = allowed;
        emit AllowlistUpdated(account, allowed);
    }

    function mint(address to, uint256 amount) external onlyRole(ISSUER_ROLE) {
        _mint(to, amount);
    }

    function burn(address from, uint256 amount) external onlyRole(ISSUER_ROLE) {
        _burn(from, amount);
    }

    /// @dev Toda movimentação exige remetente e destinatário na allowlist.
    ///      Na emissão (`from == 0`) e na queima (`to == 0`), só a contraparte é checada.
    function _update(address from, address to, uint256 value) internal override {
        if (from != address(0) && !isAllowed[from]) revert NotAllowlisted(from);
        if (to != address(0) && !isAllowed[to]) revert NotAllowlisted(to);
        super._update(from, to, value);
    }
}
