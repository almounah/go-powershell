package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/almounah/go-powershell/pkg/logger/kloghelper"
	"github.com/almounah/go-powershell/pkg/powershell"
	"k8s.io/klog"
)

type callbackTest struct{}

func (c callbackTest) Callback(_ powershell.Runspace, str string, input []powershell.Object, results powershell.CallbackResultsWriter) {
	fmt.Println("\tIn callback:", str)
	results.WriteString(str)
	for i, object := range input {
		if object.IsNull() {
			fmt.Println("\tIn callback: index", i, "Object Is Null") // ToString and Type are still valid
		}
		fmt.Println("\tIn callback: index", i, "type:", object.Type(), "with value:", object.ToString())
		results.Write(object, false)
	}
}

// PrintAndExecuteCommand executes a command in powershell and prints the results
func PrintAndExecuteCommand(runspace powershell.Runspace, command string, useLocalScope bool) {
	fmt.Println("Executing powershell command:", command)

	// determine if executing just a .ps1 file, if so use command, otherwise script
	var results *powershell.InvokeResults
	if strings.HasSuffix(command, ".ps1") {
		results = runspace.ExecCommand(command, useLocalScope, nil)
	} else {
		results = runspace.ExecScript(command, useLocalScope, nil)
	}
	defer results.Close()

	fmt.Println("Completed Executing powershell command:", command)
	if !results.Success() {
		fmt.Println("\tCommand threw exception of type", results.Exception.Type(), "and ToString", results.Exception.ToString())
	} else {
		fmt.Println("Command returned", len(results.Objects), "objects")
		for i, object := range results.Objects {
			fmt.Println("\tObject", i, "is of type", object.Type(), "and ToString", object.ToString())
		}
	}
}

// Example on how to use powershell wrappers
func Example() {
	runspace := powershell.CreateRunspace(kloghelper.Klog{VerboseLevel: 1, DebugLevel: 2}, callbackTest{})
	defer runspace.Close()

	if len(commandFlags) == 0 {
		klog.Exit("Did not specify a \"-command\" to run")
	}
	for i := 0; i < len(commandFlags); i++ {
		command := strings.ReplaceAll(commandFlags[i], "\\", "\\\\")
		PrintAndExecuteCommand(runspace, command, *useLocalScope)

	}
}

type arrayCommandFlags []string

func (i *arrayCommandFlags) String() string {
	return "my string representation"
}

func (i *arrayCommandFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var commandFlags arrayCommandFlags
var useLocalScope = flag.Bool("useLocalScope", false, "True if should execute scripts in the local scope")

func main() {
	klog.InitFlags(nil)
	flag.Var(&commandFlags, "command", "Command to run in powershell")
	flag.Set("logtostderr", "true")
	flag.Parse()
	Example()
	klog.Flush()
}
