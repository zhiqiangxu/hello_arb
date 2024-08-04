// SPDX-License-Identifier: MIT
// OpenZeppelin Contracts v4.4.1 (utils/introspection/ERC165Checker.sol)

pragma solidity ^0.8.0;

struct Path {
    Exchange exchange;
    address pool;
    address to;
    uint256 fee;
}

enum Exchange {
    UNISWAP,
    BAKERY
}

interface IPair {
    function token0() external view returns (address);

    function token1() external view returns (address);

    function getReserves()
        external
        view
        returns (
            uint112 reserve0,
            uint112 reserve1,
            uint32 blockTimestampLast
        );

    function swap(
        uint256 amount0Out,
        uint256 amount1Out,
        address to,
        bytes calldata data
    ) external;
}

interface IBPair {
    function token0() external view returns (address);

    function token1() external view returns (address);

    function getReserves()
        external
        view
        returns (
            uint112 reserve0,
            uint112 reserve1,
            uint32 blockTimestampLast
        );

    function swap(
        uint256 amount0Out,
        uint256 amount1Out,
        address to
    ) external;
}

interface ISwapExecutor {
    function addOperator(address[] calldata operators) external;

    function transfer(
        address token,
        address to,
        uint256 value
    ) external;

    function swap(
        address tokenIn,
        uint256 amountIn,
        uint256 amountOutMin,
        Path[] calldata paths
    ) external;

    function swapFlashloan(
        address tokenIn,
        uint256 amountIn,
        uint256 amountOutMin,
        Path[] memory paths
    ) external;
}
