package main

import (
	"fmt"

	commands "github.com/eris-ltd/eris-pm/cmd"
	"github.com/eris-ltd/eris-pm/version"

	"github.com/eris-ltd/common/go/common"
)

var RENDER_DIR = fmt.Sprintf("./docs/eris-pm/%s/", version.VERSION)

var SPECS_DIR = "./docs/"

var BASE_URL = fmt.Sprintf("https://docs.erisindustries.com/documentation/eris-pm/%s/", version.VERSION)

const FRONT_MATTER = `---

layout:     documentation
title:      "Documentation | eris:pm | {{}}"

---

`

func main() {
	epm := commands.EPMCmd
	commands.InitEPM()
	commands.AddGlobalFlags()
	specs := common.GenerateSpecs(SPECS_DIR, RENDER_DIR, FRONT_MATTER)
	common.GenerateTree(epm, RENDER_DIR, specs, FRONT_MATTER, BASE_URL)
}
