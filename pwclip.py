#!/usr/bin/env python
# pwclip - a hash-based password manager <https://github.com/davidlazar/pwclip>
import argparse
import hashlib
import hmac
import os
import subprocess
import sys
import time

import yaml

version = '0.1.0'
envkey = 'PWCLIP_KEYFILE'
pwm_defaults = {
    'length': 32,
    'prefix': '',
    'charset': 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
}


def genpass(key, pwm, secret=None):
    h = hmac.new(key, digestmod=hashlib.sha512)
    h.update(pwm['url'].encode('UTF-8'))
    h.update(pwm['username'].encode('UTF-8'))
    if secret:
        h.update(pwm['s' + str(secret)].encode('UTF-8'))

    # Run the hash function many times to get a digest of desired length.
    # A sponge would be more elegant.
    i = 0
    d = b''
    n = pwm['length']
    while len(d) < n:
        d += h.digest()
        h.update(h.digest() + bytes([i % 256]))
        i += 1

    p = baseX(d, pwm['charset'])
    p = pwm['prefix'] + p
    p = p[:n]

    return p


def genpass_file(keyfile, pwmfile, secret=None):
    with open(pwmfile, 'r') as f:
        yam = yaml.load(f)
    
    with open(keyfile, 'rb') as f:
        key = f.read()

    pwm = pwm_defaults.copy()
    pwm.update(yam)

    return genpass(key, pwm, secret)


def pwclip(pw):
    cb = get_clipboard() 
    set_clipboard(pw)
    print('Password copied to clipboard for 10 seconds.')
    time.sleep(10)
    set_clipboard(cb)


def argparser():
    desc = 'pwclip {}\nhttps://github.com/davidlazar/pwclip'.format(version)
    parser = argparse.ArgumentParser(description=desc,
        formatter_class=argparse.RawDescriptionHelpFormatter)
    parser.add_argument('FILE', help='password settings in YAML format')
    parser.add_argument('-p', action='store_true',
        help='print password instead of copying it to the clipboard')
    parser.add_argument('-s', type=int, metavar='N',
        help='print answer to secret question N')
    return parser


def main():
    parser = argparser()
    args = parser.parse_args()

    if not envkey in os.environ:
        sys.stderr.write('error: no value for environment variable {}\n'.format(envkey))
        sys.exit(1)

    keyfile = os.environ[envkey]
    pwmfile = args.FILE

    pw = genpass_file(keyfile, pwmfile, secret=args.s)
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


if __name__ == "__main__":
    main()
