package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/eris-ltd/common/go/docs"
	commands "github.com/eris-ltd/eris-pm/cmd"
	"github.com/eris-ltd/eris-pm/version"

	"github.com/spf13/cobra"
)

// repos should customize the next two lines
var Description = "Package Management Tooling"
var RenderDir = fmt.Sprintf("./docs/documentation/pm/%s/", version.VERSION)

type Cmd struct {
	Command     *cobra.Command
	Entry       *docs.Entry
	Description string
}

var Specs = []*docs.Entry{}
var Examples = []*docs.Entry{}
var SpecsDir = "./docs/specs"
var ExamplesDir = "./docs/examples"

func RenderFiles(cmdRaw *cobra.Command, tmpl *template.Template) error {
	this_entry := &docs.Entry{
		Title:          cmdRaw.CommandPath(),
		Specifications: Specs,
		Examples:       Examples,
		BaseURL:        strings.Replace(RenderDir, ".", "", 1), // if master run twice with version as "latest"
		Template:       tmpl,
		FileName:       docs.GenerateFileName(RenderDir, cmdRaw.CommandPath()),
	}

	cmd := &Cmd{
		Command:     cmdRaw,
		Entry:       this_entry,
		Description: Description,
	}

	for _, command := range cmd.Command.Commands() {
		RenderFiles(command, tmpl)
	}

	if !cmd.Command.HasParent() {
		entries := append(cmd.Entry.Specifications, cmd.Entry.Examples...)
		for _, entry := range entries {
			entry.Specifications = cmd.Entry.Specifications
			entry.Examples = cmd.Entry.Examples
			entry.CmdEntryPoint = cmd.Entry.Title
			entry.BaseURL = cmd.Entry.BaseURL
			if err := docs.RenderEntry(entry); err != nil {
				return err
			}
		}
	}

	outFile, err := os.Create(cmd.Entry.FileName)
	if err != nil {
		return err
	}
	defer outFile.Close()

	err = cmd.Entry.Template.Execute(outFile, cmd)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// repos should populate the top level command object
	pm := commands.EPMCmd
	commands.InitEPM()
	commands.AddGlobalFlags()

	// make the proper directory
	var err error
	if _, err = os.Stat(RenderDir); os.IsNotExist(err) {
		err = os.MkdirAll(RenderDir, 0755)
		if err != nil {
			panic(err)
		}
	}

	// generate specs and examples files
	Specs, err = docs.GenerateEntries(SpecsDir, (RenderDir + "specifications/"), Description)
	if err != nil {
		panic(err)
	}
	Examples, err = docs.GenerateEntries(ExamplesDir, (RenderDir + "examples/"), Description)
	if err != nil {
		panic(err)
	}

	// Get template from docs generator
	tmpl, err := docs.GenerateCommandsTemplate()
	if err != nil {
		panic(err)
	}

	// render templates
	if err = RenderFiles(pm, tmpl); err != nil {
		panic(err)
	}
}
