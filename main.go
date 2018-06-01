package main

import (
	"flag"
)

func main() {
	var fileName = flag.String("file", "", "アセンブリファイルを指定する")
	var outputFileName = flag.String("output", "", "出力ファイルを指定する")

	flag.Parse()

	err := assemble(*fileName, *outputFileName)
	if err != nil {
		println(err.Error())
	}
}
