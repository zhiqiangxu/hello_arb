// SPDX-License-Identifier: MIT
// OpenZeppelin Contracts v4.4.1 (utils/introspection/ERC165Checker.sol)

pragma solidity ^0.8.13;

import "openzeppelin-contracts/contracts/token/ERC20/ERC20.sol";

interface IWERC20 {
    function deposit(uint256 amount) external payable;
}

contract Wrapper {
    constructor(address wtoken) payable {
        uint256 balanceBefore = IERC20(wtoken).balanceOf(address(this));
        IWERC20(wtoken).deposit{value: msg.value}(msg.value);

        uint256 diff = IERC20(wtoken).balanceOf(address(this)) - balanceBefore;
        bytes memory result = abi.encode(diff);
        assembly {
            return(add(result, 32), mload(result))
        }
    }
}
