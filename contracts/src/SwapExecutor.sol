// SPDX-License-Identifier: MIT
// OpenZeppelin Contracts v4.4.1 (utils/introspection/ERC165Checker.sol)

pragma solidity ^0.8.13;

import "openzeppelin-contracts/contracts/access/Ownable.sol";
import "openzeppelin-contracts/contracts/access/AccessControl.sol";
import "openzeppelin-contracts/contracts/security/ReentrancyGuard.sol";
import "openzeppelin-contracts/contracts/token/ERC20/ERC20.sol";
import "./interfaces/ISwapExecutor.sol";
import "./interfaces/IPancakeCallee.sol";
import "./interfaces/IBiswapCallee.sol";
import "./interfaces/INomiswapCallee.sol";
import "./interfaces/IswapV2Callee.sol";
import "forge-std/console2.sol";

contract SwapExecutor is
    ISwapExecutor,
    IPancakeCallee,
    IBiswapCallee,
    INomiswapCallee,
    IswapV2Callee,
    Ownable,
    AccessControl,
    ReentrancyGuard
{
    bytes32 public constant OPERATOR_ROLE = keccak256("OPERATOR_ROLE");
    address flashloanPair;

    constructor(address newOwner) {
        if (newOwner != address(0)) {
            _grantRole(OPERATOR_ROLE, msg.sender);
            transferOwnership(newOwner);
        }
    }

    function addOperator(address[] calldata operators)
        external
        override
        onlyOwner
    {
        for (uint256 i = 0; i < operators.length; ++i) {
            _grantRole(OPERATOR_ROLE, operators[i]);
        }
    }

    function transfer(
        address token,
        address to,
        uint256 value
    ) public onlyOwner {
        if (token == address(0)) {
            payable(to).transfer(value);
        } else {
            IERC20(token).transfer(to, value);
        }
    }

    function _isInOrder(address tokenA, address tokenB)
        internal
        pure
        returns (bool)
    {
        return tokenA < tokenB;
    }

    function swap(
        address tokenIn,
        uint256 amountIn,
        uint256 amountOutMin, // 仅call时>0
        Path[] memory paths
    ) external override nonReentrant {
        require(
            owner() == msg.sender || hasRole(OPERATOR_ROLE, msg.sender),
            "O"
        );

        IERC20(tokenIn).transfer(paths[0].pool, amountIn);

        uint256 balanceBefore = IERC20(paths[paths.length - 1].to).balanceOf(
            address(this)
        );

        for (uint256 i; i < paths.length; i++) {
            uint256 amountInput;
            uint256 amountOutput;
            {
                (uint256 reserve0, uint256 reserve1, ) = IBPair(paths[i].pool)
                    .getReserves();
                (uint256 reserveInput, uint256 reserveOutput) = _isInOrder(
                    tokenIn,
                    paths[i].to
                )
                    ? (reserve0, reserve1)
                    : (reserve1, reserve0);
                amountInput =
                    IERC20(tokenIn).balanceOf(address(paths[i].pool)) -
                    reserveInput;
                amountOutput = _getAmountOut(
                    amountInput,
                    reserveInput,
                    reserveOutput,
                    paths[i].fee
                );
            }

            (uint256 amount0Out, uint256 amount1Out) = _isInOrder(
                tokenIn,
                paths[i].to
            )
                ? (uint256(0), amountOutput)
                : (amountOutput, uint256(0));
            address to = i < paths.length - 1
                ? paths[i + 1].pool
                : address(this);

            if (paths[i].exchange == Exchange.BAKERY) {
                IBPair(paths[i].pool).swap(amount0Out, amount1Out, to);
            } else {
                IPair(paths[i].pool).swap(
                    amount0Out,
                    amount1Out,
                    to,
                    new bytes(0)
                );
            }
            tokenIn = paths[i].to;
        }

        require(
            IERC20(paths[paths.length - 1].to).balanceOf(address(this)) -
                balanceBefore >=
                amountOutMin,
            "3"
        );
    }

    function _getAmountOut(
        uint256 amountIn,
        uint256 reserveIn,
        uint256 reserveOut,
        uint256 fee
    ) internal pure returns (uint256 amountOut) {
        uint256 amountInWithFee = amountIn * (10000 - fee);
        uint256 numerator = amountInWithFee * reserveOut;
        uint256 denominator = reserveIn * 10000 + amountInWithFee;
        amountOut = numerator / denominator;
    }

    // ----------- swapFlashloan start -----------
    function swapFlashloan(
        address tokenIn,
        uint256 amountIn,
        uint256 amountOutMin,
        Path[] memory paths
    ) external override nonReentrant {
        require(
            owner() == msg.sender || hasRole(OPERATOR_ROLE, msg.sender),
            "O"
        );
        uint256[] memory amounts = _getAmountsOut(tokenIn, amountIn, paths);
        require(amounts[amounts.length - 1] >= amountOutMin, "A");

        uint256 lastIdx = paths.length - 1;
        uint256 balanceBefore = IERC20(paths[lastIdx].to).balanceOf(
            address(this)
        );

        (uint256 amount0Out, uint256 amount1Out) = _isInOrder(
            paths[paths.length - 2].to,
            tokenIn
        )
            ? (uint256(0), amounts[lastIdx])
            : (amounts[lastIdx], uint256(0));

        require(paths[lastIdx].exchange == Exchange.UNISWAP, "E");

        bytes memory data = abi.encode(tokenIn, amountIn, paths, amounts);
        flashloanPair = paths[lastIdx].pool;
        IPair(paths[lastIdx].pool).swap(
            amount0Out,
            amount1Out,
            address(this),
            data
        );
        flashloanPair = address(0);

        require(
            IERC20(paths[lastIdx].to).balanceOf(address(this)) > balanceBefore,
            "3"
        );
    }

    function _afterFlashloan(bytes calldata data) internal {
        (
            address tokenIn,
            uint256 amountIn,
            Path[] memory paths,
            uint256[] memory amounts
        ) = abi.decode(data, (address, uint256, Path[], uint256[]));

        IERC20(tokenIn).transfer(paths[0].pool, amountIn);
        for (uint256 i = 0; i < paths.length - 1; i++) {
            (uint256 amount0Out, uint256 amount1Out) = _isInOrder(
                tokenIn,
                paths[i].to
            )
                ? (uint256(0), amounts[i])
                : (amounts[i], uint256(0));
            address to = paths[i + 1].pool;
            if (paths[i].exchange == Exchange.BAKERY) {
                IBPair(paths[i].pool).swap(amount0Out, amount1Out, to);
            } else {
                IPair(paths[i].pool).swap(
                    amount0Out,
                    amount1Out,
                    to,
                    new bytes(0)
                );
            }
            tokenIn = paths[i].to;
        }
    }

    function _onPairCallback(
        address,
        uint256,
        uint256,
        bytes calldata data
    ) internal {
        require(msg.sender == flashloanPair, "F");
        _afterFlashloan(data);
    }

    function pancakeCall(
        address sender,
        uint256 amount0,
        uint256 amount1,
        bytes calldata data
    ) external override {
        _onPairCallback(sender, amount0, amount1, data);
    }

    function BiswapCall(
        address sender,
        uint256 amount0,
        uint256 amount1,
        bytes calldata data
    ) external override {
        _onPairCallback(sender, amount0, amount1, data);
    }

    function nomiswapCall(
        address sender,
        uint256 amount0,
        uint256 amount1,
        bytes calldata data
    ) external override {
        _onPairCallback(sender, amount0, amount1, data);
    }

    function swapV2Call(
        address sender,
        uint256 amount0,
        uint256 amount1,
        bytes calldata data
    ) external override {
        _onPairCallback(sender, amount0, amount1, data);
    }

    function _getAmountsOut(
        address tokenIn,
        uint256 amountIn,
        Path[] memory paths
    ) internal view returns (uint256[] memory amounts) {
        amounts = new uint256[](paths.length);
        uint256 iAmountIn = amountIn;
        address iTokenIn = tokenIn;
        for (uint32 i = 0; i < paths.length; i++) {
            (amounts[i]) = _getSingleAmountOut(
                paths[i].pool,
                iAmountIn,
                paths[i].fee,
                _isInOrder(iTokenIn, paths[i].to)
            );
            iAmountIn = amounts[i];
            iTokenIn = paths[i].to;
        }
    }

    function _getSingleAmountOut(
        address pool,
        uint256 amountIn,
        uint256 fee,
        bool dir
    ) internal view returns (uint256 amountOut) {
        IPair pair = IPair(pool);

        (uint256 reserve0, uint256 reserve1, ) = pair.getReserves();
        uint256 reserveIn;
        uint256 reserveOut;
        if (dir) {
            (reserveIn, reserveOut) = (reserve0, reserve1);
        } else {
            (reserveIn, reserveOut) = (reserve1, reserve0);
        }

        amountOut = _getAmountOut(amountIn, reserveIn, reserveOut, fee);
        return amountOut;
    }

    // ----------- swapFlashloan end -----------
}
