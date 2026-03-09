#!/bin/env python3

import sys
from pathlib import Path

import matplotlib.pyplot as plt
import matplotlib.ticker as ticker
import pandas as pd


def timeFormatter(t, pos):
    if t < 10e-6:
        return "%dns" % (t * 1e9)
    elif t < 10e-3:
        return "%dμs" % (t * 1e6)
    elif t < 1:
        return "%dms" % (t * 1e3)
    elif t < 1e3:
        return "%.2fs" % (t)
    else:
        return "%.2es" % (t)


def bytesFormatter(n, pos):
    if n < 1e3:
        return "%db" % (n)
    elif n < 1e6:
        return "%.1fkb" % (n / 1e3)
    elif n < 1e9:
        return "%.1fMb" % (n / 1e6)
    elif n < 1e12:
        return "%.1fGb" % (n / 1e9)
    else:
        return "%.1fTb" % (n / 1e12)


def cmp(dfs, col):
    unit = COLS_2_UNIT[col]
    fmt = COLS_2_FMT[col]

    for (name, df) in dfs.items():
        mean = df.groupby(by="fileSize")[[col]].mean()
        plt.scatter(mean.index, mean.values, label=name)

        for (idx, row) in mean.iterrows():
            t = row.values[0]
            plt.gca().annotate(fmt(t), (idx, t), xytext=(
                10, 0), textcoords='offset points')

    plt.title("%s" % col)
    plt.legend()

    plt.xlabel("fileSize")
    plt.xscale("log")
    plt.gca().xaxis.set_major_formatter(ticker.FuncFormatter(bytesFormatter))

    plt.ylabel(unit)
    plt.yscale("log")
    plt.gca().yaxis.set_major_formatter(fmt)

    plt.show()


COLS = [
    "sigSize",
    "sigTime",
    "sigMem",
    "deltaSize",
    "deltaTime",
    "deltaMem",
    "patchTime",
    "patchMem",
]

COLS_2_UNIT = {
    "sigSize": "size (bytes)",
    "sigTime": "mean time",
    "sigMem": "size (bytes)",
    "deltaSize": "size (bytes)",
    "deltaTime": "mean time",
    "deltaMem": "size (bytes)",
    "patchTime": "mean time",
    "patchMem": "size (bytes)",
}

COLS_2_FMT = {
    "sigSize": ticker.FuncFormatter(bytesFormatter),
    "sigTime": ticker.FuncFormatter(timeFormatter),
    "sigMem": ticker.FuncFormatter(bytesFormatter),
    "deltaSize": ticker.FuncFormatter(bytesFormatter),
    "deltaTime": ticker.FuncFormatter(timeFormatter),
    "deltaMem": ticker.FuncFormatter(bytesFormatter),
    "patchTime": ticker.FuncFormatter(timeFormatter),
    "patchMem": ticker.FuncFormatter(bytesFormatter),
}


def printUsage():
    binName = sys.argv[0]
    print("""Compare a given stat from the provided csv files and displays a plot.

Usage: %s STAT CSV_FILE...

Arguments:
    STAT
        can be one of %s

    CSV_FILE...
        csv files containing the given stat""" % (binName, ", ".join(COLS)))


def main(argv):
    if len(argv) < 2:
        printUsage()
        exit(1)

    col = argv[0]
    if col not in COLS:
        print("column must be one off: %s" % ", ".join(COLS))
        printUsage()
        exit(1)

    dfs = {}
    for arg in argv[1:]:
        df = pd.read_csv(arg)

        name = Path(arg).stem
        dfs[name] = df

    cmp(dfs, col)


if __name__ == "__main__":
    main(sys.argv[1:])
