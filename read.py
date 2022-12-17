from os import sys

def main():
    table, err = read_tsv(sys.stdin)
    if err != None:
        sys.exit(err)
    print(table[:10])

def read_tsv(f):
    if False:
        return None, Exception("read_tsv: cannot read first line")
    head = f.readline()[:-1].split("\t")
    table = [
        {head[i]: v for i, v in enumerate(l[:-1].split("\t"))} for l in f
    ]
    return table, None

if __name__ == "__main__":
    main()
