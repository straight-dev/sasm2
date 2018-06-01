#!/usr/bin/env python3

with open("specification.csv", "r") as fp:
    lines = fp.readlines()

out = ""
for i, l in enumerate(lines):
    if l.startswith("#"):
        continue
    l = l.strip()
    opcode, name, regtype, immsext = l.split(",")
    immsext = immsext.replace("Y", "S").replace("N", "U").replace("x", "U")
    out += "{0}{1}, ".format(regtype, immsext)
    if i % 8 == 0:
        out += "\n"
print(out)
