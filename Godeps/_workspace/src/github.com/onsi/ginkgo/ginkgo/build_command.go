package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cloudfoundry-incubator/docker-circus/Godeps/_workspace/src/github.com/onsi/ginkgo/ginkgo/interrupthandler"
	"github.com/cloudfoundry-incubator/docker-circus/Godeps/_workspace/src/github.com/onsi/ginkgo/ginkgo/testrunner"
)

func BuildBuildCommand() *Command {
	commandFlags := NewBuildCommandFlags(flag.NewFlagSet("build", flag.ExitOnError))
	interruptHandler := interrupthandler.NewInterruptHandler()
	builder := &SpecBuilder{
		commandFlags:     commandFlags,
		interruptHandler: interruptHandler,
	}

	return &Command{
		Name:         "build",
		FlagSet:      commandFlags.FlagSet,
		UsageCommand: "ginkgo build <FLAGS> <PACKAGES>",
		Usage: []string{
			"Build the passed in <PACKAGES> (or the package in the current directory if left blank).",
			"Accepts the following flags:",
		},
		Command: builder.BuildSpecs,
	}
}

type SpecBuilder struct {
	commandFlags     *RunWatchAndBuildCommandFlags
	interruptHandler *interrupthandler.InterruptHandler
}

func (r *SpecBuilder) BuildSpecs(args []string, additionalArgs []string) {
	r.commandFlags.computeNodes()

	suites, _ := findSuites(args, r.commandFlags.Recurse, r.commandFlags.SkipPackage, false)

	if len(suites) == 0 {
		complainAndQuit("Found no test suites")
	}

	passed := true
	for _, suite := range suites {
		runner := testrunner.New(suite, 1, false, r.commandFlags.Race, r.commandFlags.Cover, r.commandFlags.Tags, nil)
		fmt.Printf("Compiling %s...\n", suite.PackageName)
		err := runner.Compile()
		if err != nil {
			fmt.Println(err.Error())
			passed = false
		} else {
			fmt.Printf("    compiled %s.test\n", filepath.Join(suite.Path, suite.PackageName))
		}
	}

	if passed {
		os.Exit(0)
	}
	os.Exit(1)
}