#!/bin/sh

pwclip=$(realpath pwclip.py)
expected=$(realpath tests/expected)
actual=$(realpath tests/actual)

key="secret key"

# for python 2
export PYTHONIOENCODING=utf-8

runtests () {
cd tests
echo -n $key | $pwclip -p -k /dev/stdin test01
$pwclip -p -c "echo -n $key" test01
echo -n $key | $pwclip -p -k /dev/stdin test02
echo -n $key | $pwclip -p -k /dev/stdin test02 -q1
$pwclip -p -k test02 test02
echo -n $key | $pwclip -p -k /dev/stdin test03
echo -n $key | $pwclip -p -k /dev/stdin test04
echo -n $key | $pwclip -p -k /dev/stdin test05
}

runtests 2>/dev/null > "$actual"

if cmp -s "$actual" "$expected"; then
    echo "ok"
    rm "$actual"
    exit 0
else
    echo "test failed: see tests/actual"
    exit 1
fi
