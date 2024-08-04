// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import "./Utils.sol";

import "forge-std/console2.sol";

library ZeroCopySource {
    function NextBool(bytes memory buff, uint256 offset)
        internal
        pure
        returns (bool, uint256)
    {
        require(
            offset + 1 <= buff.length && offset < offset + 1,
            "Offset exceeds limit"
        );

        // 这儿不能是uint8，否则truncate方式不一样
        bytes1 v;
        assembly {
            v := mload(add(add(buff, 0x20), offset))
        }

        if (v == 0x01) {
            return (true, offset + 1);
        } else if (v == 0x00) {
            return (false, offset + 1);
        } else {
            revert("NextBool value error");
        }
    }

    function NextByte(bytes memory buff, uint256 offset)
        internal
        pure
        returns (bytes1, uint256)
    {
        require(
            offset + 1 <= buff.length && offset < offset + 1,
            "NextByte, Offset exceeds maximum"
        );
        bytes1 v;
        assembly {
            v := mload(add(add(buff, 0x20), offset))
        }
        return (v, offset + 1);
    }

    function NextUint8(bytes memory buff, uint256 offset)
        internal
        pure
        returns (uint8, uint256)
    {
        (bytes1 v, uint256 newOffset) = NextByte(buff, offset);
        return (uint8(v), newOffset);
    }

    function NextUint16(bytes memory buff, uint256 offset)
        internal
        pure
        returns (uint16, uint256)
    {
        require(
            offset + 2 <= buff.length && offset < offset + 2,
            "NextUint16, offset exceeds maximum"
        );

        bytes2 v;
        assembly {
            v := mload(add(add(buff, 0x20), offset))
        }
        return (uint16(v), offset + 2);
    }

    function NextUint32(bytes memory buff, uint256 offset)
        internal
        pure
        returns (uint32, uint256)
    {
        require(
            offset + 4 <= buff.length && offset < offset + 4,
            "NextUint32, offset exceeds maximum"
        );

        bytes4 v;
        assembly {
            v := mload(add(add(buff, 0x20), offset))
        }
        return (uint32(v), offset + 4);
    }

    function NextUint64(bytes memory buff, uint256 offset)
        internal
        pure
        returns (uint64, uint256)
    {
        require(
            offset + 8 <= buff.length && offset < offset + 8,
            "NextUint64, offset exceeds maximum"
        );

        bytes8 v;
        assembly {
            v := mload(add(add(buff, 0x20), offset))
        }
        return (uint64(v), offset + 8);
    }

    function NextUint256(bytes memory buff, uint256 offset)
        internal
        pure
        returns (uint256, uint256)
    {
        require(
            offset + 32 <= buff.length && offset < offset + 32,
            "NextUint256, offset exceeds maximum"
        );

        uint256 v;
        assembly {
            v := mload(add(add(buff, 0x20), offset))
        }

        return (v, offset + 32);
    }

    function NextBytes32(bytes memory buff, uint256 offset)
        internal
        pure
        returns (bytes32, uint256)
    {
        require(
            offset + 32 <= buff.length && offset < offset + 32,
            "NextUint256, offset exceeds maximum"
        );

        bytes32 v;
        assembly {
            v := mload(add(add(buff, 0x20), offset))
        }

        return (v, offset + 32);
    }

    function NextBytes20(bytes memory buff, uint256 offset)
        internal
        pure
        returns (bytes20, uint256)
    {
        require(
            offset + 20 <= buff.length && offset < offset + 20,
            "NextBytes20, offset exceeds maximum"
        );
        bytes20 v;
        assembly {
            v := mload(add(buff, add(offset, 0x20)))
        }
        return (v, offset + 20);
    }

    function NextVarUint64(bytes memory buff, uint256 offset)
        internal
        pure
        returns (uint64, uint256)
    {
        bytes1 v;
        (v, offset) = NextByte(buff, offset);

        uint64 value;
        if (v == 0xFD) {
            // return NextUint16(buff, offset);
            (value, offset) = NextUint16(buff, offset);
            require(
                value >= 0xFD && value <= 0xFFFF,
                "NextUint16, value outside range"
            );
            return (value, offset);
        } else if (v == 0xFE) {
            // return NextUint32(buff, offset);
            (value, offset) = NextUint32(buff, offset);
            require(
                value > 0xFFFF && value <= 0xFFFFFFFF,
                "NextVarUint, value outside range"
            );
            return (value, offset);
        } else if (v == 0xFF) {
            // return NextUint64(buff, offset);
            (value, offset) = NextUint64(buff, offset);
            require(value > 0xFFFFFFFF, "NextVarUint, value outside range");
            return (value, offset);
        } else {
            // return (uint8(v), offset);
            value = uint8(v);
            require(value < 0xFD, "NextVarUint, value outside range");
            return (value, offset);
        }
    }

    function NextVarBytes(bytes memory buff, uint256 offset)
        internal
        pure
        returns (bytes memory, uint256)
    {
        uint64 len;
        (len, offset) = NextVarUint64(buff, offset);

        require(
            offset + len <= buff.length && offset <= offset + len,
            "NextVarBytes, offset exceeds maximum"
        );

        return (Utils.slice(buff, offset, len), offset + len);
    }
}
