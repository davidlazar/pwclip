pwclip is a hash-based, command-line password manager.  pwclip does not store
passwords.  Instead, it computes an account's password by hashing a secret key
together with account-specific information stored in a YAML file.

I've tested the pwclip on Python 2.7.5 and Python 3.3.2.

Usage
-----

0.  Read the code (it's short!) to understand how passwords are generated.


1.  Generate a random key:

        $ dd if=/dev/random of=secret_key bs=1 count=128
        $ export PWCLIP_KEYFILE=`realpath secret_key`

    It's your responsibility to protect the key, for example using GPG.
    I keep mine in an encrypted directory.


2.  Create a separate YAML file containing the password settings for each
    account.  Here is a minimal example:

        $ cat github
        url: github.com
        username: davidlazar

    Here is an example that gives a value to every recognized field:

        $ cat amazon
        url: amazon.com
        username: example@example.com
        length: 48
        prefix: foobar
        charset: ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789`~!@#$%^&*()_-+={}|[]\:";'<>?,./
        s1: favorite band
        s2: city born

    The `s1`..`sN` fields are used to give unique answers to secret questions
    used for password recovery.


3.  Without any flags, the password is copied to the clipboard using `xclip`:

        $ pwclip.py github
        Password copied to clipboard for 10 seconds.

    The `-p` flag simply prints the password to the screen:

        $ pwclip.py -p amazon
        foobarGQXlD=sw|~U1JC-fd.dFS$Gio);o)Txm5s{zL~FPvr

    The `-s N` flag generates the answer to secret question N:

        $ pwclip.py -p amazon -s1
        foobarB7/*Oc-vP+55s[K9@Duqw<s4L]!q7I%dx:N{)_CW~0

        $ pwclip.py -p amazon -s2
        foobar^W_b$agBVU!m56C~NnKcke3G?#i&@pchg'OSWuW,#L

    Note that the answers use the same settings as the password.
