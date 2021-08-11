#!/bin/sh

# Measures CPU time, memory usage and resulting file sizes for the various
# operations supported by librsync-go (signature, delta, patch) when invoked
# with different parameters and different sized inputs.
#
# IMPORTANT: You will need to adjust `benchmarkForSize()` to get the
# combinations of parameters you want. The implementation is cursed by
# combinatorial explosion and the defaults may take more than a lifetime to run.
#
# See `create-benchmark-data.sh` for the script used to generate the input data
# required by this script.
#
# This scripts expects to find GNU time at the location indicated by the `TIME`
# variable defined below. We rely on GNU extensions to `time`, so the shell
# built-in is not enough.
#
# This script's output is a CSV file with the following fields:
#
# - FileSize: The "reference file size", in KB (like 10 for 10 KB files or
#   10000000 for 10GB files)
# - BasisFile: The file name used as the basis (old) file.
# - TargetFile: The file name used as the target (new) file.
#
# - BlockSize: The "block size" argument passed to rdiff, in bytes
# - StrongSumSize: The "strong sum size" argument passed to rdiff, in bytes
#
# - SigSize: The resulting signature size, in bytes
# - SigTime: The time taken to generate the signature in seconds
# - SigMem: The maximum amount of memory used to generate the signature, in KB
# - DeltaSize: The resulting delta size, in byes
# - DeltaTime: The time taken to generate the delta in seconds
# - DeltaMem: The maximum amount of memory used to generate the delta, in KB
# - PatchTime: The time taken to apply the delta in seconds
# - PatchMem: The maximum amount of memory used to apply the delta, in KB

# TODO: Some of the values above are in bytes, others are in KB. It's annoying
# to analyze data in different units. We should probably just use bytes
# everywhere.

# GNU time is here.
TIME="/usr/bin/time"

# And this is librsync-go's `rdiff`. (This script should probably also work with
# the "original" `rdiff` written in C, but it was created to benchmark this Go
# implementation.)
RDIFF="./rdiff"


# Checks if `rdiff` and `time` are present on the expected places. Exits with
# failure if any of them is missing.
function checkForDependencies() {
    $RDIFF &> /dev/null
    if [ $? != 0 ]; then
        echo "Expected to have the librsync-go rdiff binary on the current directory!"
        exit 1
    fi

    $TIME &> /dev/null
    if [ $? != 0 ]; then
        echo "Expected to have the `time` binary at $TIME"
        exit 1
    fi
}


# Converts $1 from "minutes:seconds.subseconds" to "seconds.subseconds"
function toSeconds() {
    mins=$(echo $1 | cut -f 1 -d ':')
    secsDecimal=$(echo $1 | cut -f 2 -d ':')
    secs=$(echo $secsDecimal | cut -f 1 -d '.')
    subSecs=$(echo $secsDecimal | cut -f 2 -d '.')
    echo $((10#$mins * 60 + 10#$secs)).$subSecs
}


# Benchmarks one case, echoes the result as a line of our final CSV results
# file.
#
# Parameters:
#
# $1: The basis file name
# $2: The target file name
# $3: The block size in bytes
# $4: The strong sum size in bytes
function benchmarkOneCase() {
    outFileSize=`echo $1 | cut -d '-' -f 1`
    outBasisFile="$1"
    outTargetFile="$2"
    outBlockSize="$3"
    outStrongSumSize="$4"

    $TIME -f "%M\t%E" -o the-sig-data $RDIFF signature --block-size $3 --sum-size $4 "$1" the-sig
    outSigSize=$(stat -c %s the-sig)
    outSigTime=$(toSeconds $(cut -f 2 the-sig-data))
    outSigMem=$(cut -f 1 the-sig-data)

    $TIME -f "%M\t%E" -o the-delta-data $RDIFF delta the-sig "$2" the-delta
    outDeltaSize=$(stat -c %s the-delta)
    outDeltaTime=$(toSeconds $(cut -f 2 the-delta-data))
    outDeltaMem=$(cut -f 1 the-delta-data)

    $TIME -f "%M\t%E" -o the-patch-data $RDIFF patch "$1" the-delta the-target
    outPatchTime=$(toSeconds $(cut -f 2 the-patch-data))
    outPatchMem=$(cut -f 1 the-patch-data)

    rm the-sig the-sig-data the-delta the-delta-data the-patch-data the-target

    echo "$outFileSize,$outBasisFile,$outTargetFile,$outBlockSize,$outStrongSumSize,$outSigSize,$outSigTime,$outSigMem,$outDeltaSize,$outDeltaTime,$outDeltaMem,$outPatchTime,$outPatchMem"
}


# Benchmarks an assortment of cases for files with "reference file size" equals
# to $1. Echoes CSV rows with the results.
function benchmarkForSize() {
    for basisFile in $1-abc.data; do
        for targetFile in $1-abcx.data $1-abx.data $1-axc.data $1-b.data $1-bc.data $1-xbc.data; do
            for blockSize in 256 512 1024 2048 4096 8194 16384 32768 65536 257 513 1025 2049 4097 8195 16385 32769 65537; do
                for strongSumSize in 16 20 24 28 32; do
                    benchmarkOneCase $basisFile $targetFile $blockSize $strongSumSize
                done
            done
        done
    done
}


function printHeader() {
    echo "fileSize,basisFile,targetFile,blockSize,strongSumSize,sigSize,sigTime,sigMem,deltaSize,deltaTime,deltaMem,patchTime,patchMem"
}


#
# Main script body
#

if [ ! checkForDependencies ]; then
    exit 1
fi

printHeader

# Run benchmarks for various input file sizes
benchmarkForSize 10
benchmarkForSize 100
benchmarkForSize 1000
benchmarkForSize 10000
benchmarkForSize 100000
benchmarkForSize 1000000
benchmarkForSize 10000000

# Alternatively, if you are just playing around, you can do this to measure one
# single set of parameters:
# benchmarkOneCase 10000000-abc.data 10000000-abcx.data 2048 32
