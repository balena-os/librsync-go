#!/bin/sh

#
# Generates files to use as input for `run-benchmark.sh`.
#
# The contents of these files are random, but there is some structure on how
# they differ from each other. For example, file F2 might be just like file F2
# but with some additional data appended to it. See the docs for
# `generateRandomFiles()` below for more details.
#


# Generates a random "reference file" measuring $1 KB, and variations of it
# (e.g., with data appended to it, or some data removed from its start).
#
# The reference file has an `abc` suffix on its name, meaning it is composed
# of three parts with random data: 10% on part A, 80% on part B and 10% on part
# C.
#
# An additional part (X) is also created, with size equals to 10% of the `abc`
# file. So, an `abx` file is like the "reference file", but had its last part
# (C) replaced with a completely new part (X). An `axc` is an interesting case,
# because it had its large middle section (B) replaced with a small one (X).
#
# This same naming pattern is used in every other case.
function generateRandomFiles() {
    dd if=/dev/urandom of=mid.tmp bs=1K count=$((8 * $1/10))
    dd if=/dev/urandom of=pre.tmp bs=1K count=$(($1/10))
    dd if=/dev/urandom of=post.tmp bs=1K count=$(($1/10))
    dd if=/dev/urandom of=alt.tmp bs=1K count=$(($1/10))

    cat pre.tmp mid.tmp post.tmp > $1-abc.data
    cat mid.tmp post.tmp > $1-bc.data
    cat mid.tmp > $1-b.data
    cat pre.tmp mid.tmp post.tmp alt.tmp > $1-abcx.data
    cat pre.tmp mid.tmp alt.tmp > $1-abx.data
    cat alt.tmp mid.tmp post.tmp > $1-xbc.data
    cat pre.tmp alt.tmp post.tmp > $1-axc.data

    rm -f pre.tmp mid.tmp post.tmp alt.tmp
}


# Main script body, generates all the benchmark data we need.
generateRandomFiles 10        # 10 KB
generateRandomFiles 100       # 100 KB
generateRandomFiles 1000      # 1 MB
generateRandomFiles 10000     # 10 MB
generateRandomFiles 100000    # 100 MB
generateRandomFiles 1000000   # 1 GB
generateRandomFiles 10000000  # 10 GB
