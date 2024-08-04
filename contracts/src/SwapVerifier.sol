// SPDX-License-Identifier: MIT
// OpenZeppelin Contracts v4.4.1 (utils/introspection/ERC165Checker.sol)

pragma solidity ^0.8.13;

import "./interfaces/ISwapExecutor.sol";
import "openzeppelin-contracts/contracts/token/ERC20/ERC20.sol";

import "forge-std/console2.sol";
import "./libs/Utils.sol";

interface EthRouter {
    function swapETHForExactTokens(
        uint256 amountOutMin,
        address[] calldata path,
        address to,
        uint256 deadline
    ) external payable returns (uint256[] memory amounts);
}

interface IWETH {
    function deposit() external payable;
}

contract SwapVerifier {
    constructor(
        address arbToken,
        address wrapToken,
        address router,
        address[] memory wrapPath,
        uint256[] memory arbAmountIn,
        uint256[] memory amountOutMin,
        Path[][] memory paths
    ) payable {
        console2.log(msg.sender);
        console2.log(msg.value);

        console2.log(arbToken == wrapToken);
        console2.log("test1");
        console2.log(arbAmountIn.length);
        console2.log(amountOutMin.length);
        console2.log(paths.length);
        console2.log(paths[0].length);
        require(
            arbAmountIn.length == amountOutMin.length &&
                amountOutMin.length == paths.length,
            "REQ_INVARIANT"
        );

        uint256 totalArbAmountIn = 0;
        for (uint256 i = 0; i < arbAmountIn.length; i++) {
            totalArbAmountIn += arbAmountIn[i];
        }
        // console2.log(totalArbAmountIn);
        if (arbToken == wrapToken) {
            IWETH(wrapToken).deposit{value: msg.value}();
        } else {
            if (wrapPath.length == 0) {
                // default wrapPath
                wrapPath = new address[](2);
                wrapPath[0] = wrapToken;
                wrapPath[1] = arbToken;
            }

            EthRouter(router).swapETHForExactTokens{value: msg.value}(
                totalArbAmountIn,
                wrapPath,
                address(this),
                block.timestamp + 100
            );
        }

        bool[] memory results = new bool[](arbAmountIn.length);
        bytes[] memory returnData = new bytes[](arbAmountIn.length);

        for (uint256 i = 0; i < arbAmountIn.length; i++) {
            IERC20(arbToken).transfer(paths[i][0].pool, arbAmountIn[i]);

            // swap(arbToken, arbAmountIn, amountOutMin, paths);

            (results[i], returnData[i]) = swapSupportingFee(
                arbToken,
                amountOutMin[i],
                paths[i]
            );
        }

        bytes memory result = abi.encode(results, returnData);

        assembly {
            return(add(result, 32), mload(result))
        }
    }

    function _isInOrder(address tokenA, address tokenB)
        internal
        pure
        returns (bool)
    {
        return tokenA < tokenB;
    }

    function swapSupportingFee(
        address tokenIn,
        uint256 amountOutMin,
        Path[] memory paths
    ) internal returns (bool, bytes memory) {
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
                console2.log("_getAmountOut");
                console2.log(reserveInput, reserveOutput);
                console2.log(amountInput, amountOutput);
                console2.log(paths[i].pool);
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

            bool success;
            bytes memory reason;
            if (paths[i].exchange == Exchange.BAKERY) {
                IBPair pair = IBPair(paths[i].pool);
                (success, reason) = address(pair).call(
                    abi.encodeWithSelector(
                        pair.swap.selector,
                        amount0Out,
                        amount1Out,
                        to
                    )
                );
            } else {
                IPair pair = IPair(paths[i].pool);
                (success, reason) = address(pair).call(
                    abi.encodeWithSelector(
                        pair.swap.selector,
                        amount0Out,
                        amount1Out,
                        to,
                        new bytes(0)
                    )
                );
            }
            if (!success) {
                console2.log("fee");
                console2.log(paths[i].fee);
                // console2.log(paths[i].pool);
                return (false, reason);
            }

            tokenIn = paths[i].to;
        }

        uint256 diff = IERC20(paths[paths.length - 1].to).balanceOf(
            address(this)
        ) - balanceBefore;
        if (diff < amountOutMin) {
            return (false, "3");
        }

        // console2.log(diff);
        return (true, Utils.uint256ToBytes(diff));
    }

    function swap(
        address tokenIn,
        uint256 amountIn,
        uint256 amountOutMin,
        Path[] memory paths
    ) internal {
        uint256[] memory amounts = _getAmountsOut(tokenIn, amountIn, paths);
        if (amountOutMin > amounts[amounts.length - 1]) {
            return;
        }

        IERC20(tokenIn).transfer(paths[0].pool, amountIn);
        for (uint32 i = 0; i < paths.length; i++) {
            (uint256 amount0Out, uint256 amount1Out) = _isInOrder(
                tokenIn,
                paths[i].to
            )
                ? (uint256(0), amounts[i])
                : (amounts[i], uint256(0));
            address to = i < paths.length - 1
                ? paths[i + 1].pool
                : address(this);

            if (paths[i].exchange == Exchange.BAKERY) {
                IBPair pair = IBPair(paths[i].pool);
                pair.swap(amount0Out, amount1Out, to);
            } else {
                IPair pair = IPair(paths[i].pool);
                pair.swap(amount0Out, amount1Out, to, new bytes(0));
            }
            tokenIn = paths[i].to;
        }
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
            amounts[i] = _getSingleAmountOut(
                paths[i].pool,
                iTokenIn,
                iAmountIn,
                paths[i].fee
            );
            iAmountIn = amounts[i];
            iTokenIn = paths[i].to;
        }
    }

    function _getSingleAmountOut(
        address pool,
        address tokenIn,
        uint256 amountIn,
        uint256 fee
    ) internal view returns (uint256 amountOut) {
        IPair pair = IPair(pool);
        require(tokenIn == pair.token0() || tokenIn == pair.token1(), "T");

        (uint256 reserve0, uint256 reserve1, ) = pair.getReserves();
        (uint256 reserveIn, uint256 reserveOut) = tokenIn == pair.token0()
            ? (reserve0, reserve1)
            : (reserve1, reserve0);

        return _getAmountOut(amountIn, reserveIn, reserveOut, fee);
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
}
