# sasm2
An assembler for STRAIGHT

## Installation
    go get -u github.com/clkbug/sasm2

## Usage
    sasm2 -file input.s -output a.out

## Build
    cd straightISAv2Info
    python3 specification.csv
    go generate
    cd ../
    go build