#!/usr/bin/env python3

import argparse
import sys

def xorme(buf, k):
    res = b''
    for ch in buf:
        res += chr(ch ^ k).encode('latin1')
    return res

def print_output(buf, t):
    if t == 'go':
        fmt = r'0x{:02x}'
        end=','
        print('buf := []byte {')
    elif t == 'c#':
        fmt = r'0x{:02x}'
        end=','
        print('byte[] buf = {')
    elif t == 'py':
        fmt = r'\x{:02x}'
        end=''
        print('buf="""\\')
    for i, ch in enumerate(buf[:-1]):
        if i and not i % 16:
            if t == 'py': print('\\',end='')
            print('')
        print(fmt.format(ch), end=end)

    print(fmt.format(buf[-1]), end='')
    if t == 'go':
        print(' }')
    elif t == 'c#':
        print(' };')
    elif t == 'py':
        print('"""')

if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument(
        '-t', choices=['py', 'go', 'c#'],
        default='py',
        help='Language type to generate for'
    )
    parser.add_argument(
        '-x', type=int, default=55,
        help='XOR integer (defaults to 55)'
    )
    args = parser.parse_args()
    buf = sys.stdin.buffer.read()
    xorbuf = xorme(buf, args.x)
    print_output(xorbuf, args.t)
