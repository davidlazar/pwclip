pwclip
------

pwclip is a hash-based, command-line password manager.  pwclip does not store
passwords.  Instead, it computes an account's password by hashing a secret key
together with account-specific information stored in a YAML file.

The previous Python implementation (in the `python/` directory) is not
compatible with the Go implementation.  The pwclip algorithm is defined in the
[go-crypto repository](https://github.com/davidlazar/go-crypto).

Usage
-----

0.  `go get github.com/davidlazar/pwclip`

1.  Pick and remember a passphrase. Alternatively, you can use a key file
    with the `-k` flag.

2.  Create a separate YAML file containing the password settings for each
    account.  Here is a minimal example:

        $ cat github
        url: github.com
        username: davidlazar

    Here is an example that gives a value to every recognized field:

        $ cat example
        url: example.com
        username: example@example.com
        length: 48
        prefix: foobar
        charset: ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789`~!@#$%^&*()_-+={}|[]\:";'<>?,./
        q1: frequent flier number
        q2: first car

    The `q1`..`qN` fields are used to give unique answers to secret questions
    used for password recovery.

3.  Copy the password to the clipboard:

        $ pwclip github
        Passphrase: ...
        Password copied to clipboard for 10 seconds.

    Print the password to the screen:

        $ pwclip -p example
        Passpharse: ...
        foobarGLKyG"Cd,,Yv2w:S5Z[*p`]z3jQp^X2};nyYf<.dNK

    Use a key file instead of prompting for a passphrase:

        $ pwclip -k ~/secret/key -p example
        foobarb)+H69iq<{[%/V'8bFVRN@l2&-$iGr0PB#zK1T`*CL

    The `-q N` flag generates the answer to secret question N:

        $ pwclip -k ~/secret/key -p example -q1
        foobarhqL3d,!zsyYOrko1%I`L@&Q-mE1`%K|soR0>%BN,^D

        $ pwclip -k ~/secret/key -p example -q2
        foobarL;)eL!Ij+{E&8++*F0XO5'4APtE{>INFb4sF,d[):7

    Note that the answers use the same settings as the password.
