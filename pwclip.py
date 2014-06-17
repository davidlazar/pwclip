#!/usr/bin/env python
# pwclip - a hash-based password manager <https://github.com/davidlazar/pwclip>
from __future__ import print_function
import argparse
import hashlib
import hmac
import os
import subprocess
import sys
import time

import yaml

# Bump at least Y in version X.Y.Z whenever the password generation algorithm changes.
version = '0.3.0'
pwm_defaults = {
    'length': 32,
    'prefix': '',
    'charset': 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
}


def genpass(key, pwm, question=None):
    rng = DRBG(key)
    rng.reseed(pwm['url'].encode('UTF-8'))
    rng.reseed(pwm['username'].encode('UTF-8'))
    if question:
        rng.reseed(question.encode('UTF-8'))

    n = pwm['length']
    r = rng.generate(n)

    p = baseX(r, pwm['charset'])
    p = pwm['prefix'] + p

    return p[:n]


def readpwm(pwmfile):
    with open(pwmfile, 'r') as f:
        yam = yaml.safe_load(f)

    pwm = pwm_defaults.copy()
    pwm.update(yam)

    return pwm


def pwclip(pw):
    cb = get_clipboard()
    try:
        set_clipboard(pw)
        print('Password copied to clipboard for 10 seconds.')
        time.sleep(10)
    except KeyboardInterrupt:
        pass
    finally:
        set_clipboard(cb)


def argparser():
    desc = 'pwclip {}\nhttps://github.com/davidlazar/pwclip'.format(version)
    parser = argparse.ArgumentParser(description=desc,
        formatter_class=argparse.RawDescriptionHelpFormatter)
    group = parser.add_mutually_exclusive_group(required=True)
    group.add_argument('-k', metavar='KEYFILE', help='read key from file')
    group.add_argument('-c', metavar='COMMAND',
        help='read key from stdout of command')
    parser.add_argument('YAMLFILE', help='password settings in YAML format')
    parser.add_argument('-p', action='store_true',
        help='print password instead of copying it to the clipboard')
    parser.add_argument('-q', type=int, metavar='N',
        help='print answer to secret question N')
    return parser


def readkey(args):
    if args.c:
        return subprocess.check_output(args.c, shell=True)
    else:
        return open(args.k, 'rb').read()


def main():
    parser = argparser()
    args = parser.parse_args()

    key = readkey(args)
    pwm = readpwm(args.YAMLFILE)

    if args.q:
        question = pwm['q' + str(args.q)]
        print('question: ' + question, file=sys.stderr)
    else:
        question = None
        print('username: ' + pwm['username'], file=sys.stderr)

    pw = genpass(key, pwm, question)
    if args.p:
        print(pw)
    else:
        pwclip(pw)


### utility functions

def get_clipboard():
    if sys.platform == 'darwin':
        return subprocess.check_output(['pbpaste'], universal_newlines=True)
    else:
        return subprocess.check_output(['xclip', '-out'], universal_newlines=True)


def set_clipboard(s):
    if sys.platform == 'darwin':
        clip = 'pbcopy'
    else:
        clip = 'xclip'
    p = subprocess.Popen([clip], stdin=subprocess.PIPE, universal_newlines=True)
    p.communicate(s)
    if p.returncode != 0:
        raise subprocess.CalledProcessError(p.returncode, clip)


def baseX(bs, alphabet):
    n = from_bytes(bs)
    b = len(alphabet)

    if n == 0:
        return alphabet[0]

    r = ''
    while n:
        r = alphabet[n % b] + r
        n //= b

    return r


def from_bytes(bs):
    try:
        # python 3.2+
        return int.from_bytes(bs, byteorder='little')
    except AttributeError:
        # python 2.x
        bs = reversed(bs)
        r = 0
        for b in bs:
            r <<= 8
            r |= ord(b)
        return r


# HMAC_DRBG from https://github.com/davidlazar/python-drbg
class DRBG(object):
    def __init__(self, seed):
        self.key = b'\x00' * 64
        self.val = b'\x01' * 64
        self.reseed(seed)

    def hmac(self, key, val):
        return hmac.new(key, val, hashlib.sha512).digest()

    def reseed(self, data=b''):
        self.key = self.hmac(self.key, self.val + b'\x00' + data)
        self.val = self.hmac(self.key, self.val)

        if data:
            self.key = self.hmac(self.key, self.val + b'\x01' + data)
            self.val = self.hmac(self.key, self.val)

    def generate(self, n):
        xs = b''
        while len(xs) < n:
            self.val = self.hmac(self.key, self.val)
            xs += self.val

        self.reseed()

        return xs[:n]


if __name__ == "__main__":
    main()
