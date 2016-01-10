package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	flag.Parse()
	contentB, _ := ioutil.ReadFile("genesis.json.example")
	contents := string(contentB)
	contents = strings.Replace(contents, "1040E6521541DAB4E7EE57F21226DD17CE9F0FB7", flag.Args()[0], 2)
	contents = strings.Replace(contents, "6A3AFFB16BFB95AA547930572D71C460EFBCD857", flag.Args()[1], 1)
	fmt.Println(contents)
}
