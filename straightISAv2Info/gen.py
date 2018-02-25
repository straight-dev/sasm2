#!/usr/bin/env python3

with open("specification.csv", "r") as fp:
  lines = fp.readlines()

with open("isainfo.go.skel", "r") as fp:
  skelton = fp.read()

opcodes = ""
cases = ""
for i, l in enumerate(lines):
  if l.startswith("#"):
    continue
  
  opcode, name, regtype, immsext, _ = l.split(",", 4)
  if i == 1:
    opcodes += "\tOp{0} OpCode = iota\n".format(name)
  else:
    opcodes += "\tOp{0}\n".format(name)

  regtype = regtype.replace("Z", "ZeroReg").replace("O", "OneReg").replace("T", "TwoReg")
  immsext = immsext.replace("Y", "true").replace("N", "false").replace("x", "false")
  cases += """\tcase Op{0}:
		return {1}, {2}, nil
""".format(name, regtype, immsext)

with open("isainfo.go", "w") as fp:
  fp.write(skelton.format(opcodes, cases))
