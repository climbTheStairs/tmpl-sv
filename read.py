# This program does not actually do anything.
# I just made it for fun to see
# how short I could rewrite read.go in Python.

def read_tsv(f):
    head = f.readline()[:-1].split("\t")
    table = [{head[i]: v for i, v in enumerate(l[:-1].split("\t"))} for l in f]
    return table

