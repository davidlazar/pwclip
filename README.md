pwclip is a hash-based, command-line password manager.  pwclip does not store
passwords.  Instead, it computes an account's password by hashing a secret key
together with account-specific information stored in a YAML file.

I've tested pwclip on Python 2.7.6 and Python 3.4.0.

Usage
-----

0.  Read the code (it's short!) to understand how passwords are generated.


1.  Generate a random key:

        $ dd if=/dev/random of=pwclip_key bs=1 count=128

2.  Secure the key.  I keep mine in an encrypted directory:

        $ alias pwclip='pwclip.py -k ~/enc/pwclip_key'

    Alternatively, you can use pwclip with a master passphrase by encrypting
    the key with a tool like [scrypt](https://www.tarsnap.com/scrypt.html):

        $ scrypt enc pwclip_key pwclip_key.enc
        $ rm pwclip_key
        $ alias pwclip='pwclip.py -c "scrypt dec ~/pwclip_key.enc"'

    You will be prompted for the passphrase whenever you run pwclip.


3.  Create a separate YAML file containing the password settings for each
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
        q1: favorite band
        q2: city born

    The `q1`..`qN` fields are used to give unique answers to secret questions
    used for password recovery.


4.  Without any flags, the password is copied to the clipboard using `xclip`:

        $ pwclip github
        Password copied to clipboard for 10 seconds.

    The `-p` flag simply prints the password to the screen:

        $ pwclip -p amazon
        foobarGQXlD=sw|~U1JC-fd.dFS$Gio);o)Txm5s{zL~FPvr

    The `-q N` flag generates the answer to secret question N:

        $ pwclip -p amazon -q1
        foobarB7/*Oc-vP+55s[K9@Duqw<s4L]!q7I%dx:N{)_CW~0

        $ pwclip -p amazon -q2
        foobar^W_b$agBVU!m56C~NnKcke3G?#i&@pchg'OSWuW,#L

    Note that the answers use the same settings as the password.
