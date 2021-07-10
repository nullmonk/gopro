import sys

# A basic script which dumps a file into a go slice
width = 15
with open(sys.argv[1], 'rb') as f:
    print("var data = []byte{")
    while True:
        bytes = f.read(width)
        print("\t" + ", ".join(["0x{:02x}".format(b) for b in bytes])+",")
        if len(bytes) < width:
            break
    print("}")