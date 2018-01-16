package main

import (
	"context"
	_ "net/http/pprof" // include pprop
	"os"

	sparta "github.com/mweagle/Sparta"
	spartaDecorator "github.com/mweagle/Sparta/decorator"
)

// Standard AWS λ function
func helloWorld(ctx context.Context) (string, error) {
	return "Hello World - 🌏", nil
}

////////////////////////////////////////////////////////////////////////////////
// Main
func main() {
	lambdaFn := sparta.HandleAWSLambda("HelloSafeLambda",
		helloWorld,
		sparta.IAMRoleDefinition{})
	lambdaFn.Options.MemorySize = 128

	// Sanitize the name so that it doesn't have any spaces
	var lambdaFunctions []*sparta.LambdaAWSInfo
	lambdaFunctions = append(lambdaFunctions, lambdaFn)

	// Add the decorator for the Canary10Percent5Minutes deploy
	workflowHooks := &sparta.WorkflowHooks{}
	workflowHooks.ServiceDecorators = append(workflowHooks.ServiceDecorators,
		spartaDecorator.CodeDeployServiceUpdateDecorator("Canary10Percent5Minutes",
			lambdaFunctions,
			nil,
			nil),
	)

	err := sparta.MainEx("SafeUpdateStack",
		"Simple Sparta application that demonstrates gated updates",
		lambdaFunctions,
		nil,
		nil,
		workflowHooks,
		false)
	if err != nil {
		os.Exit(1)
	}
}
