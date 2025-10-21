package main

import (
	"fmt"
	"os"

	"github.com/uniquejava/ctx/internal/kubeconfig"
)

func main() {
	if len(os.Args) == 1 {
		// Default behavior: list contexts
		runList()
		return
	}

	command := os.Args[1]

	switch command {
	case "ls":
		runList()
	case "use":
		if len(os.Args) < 3 {
			fmt.Println("Error: context name required")
			fmt.Println("Usage: ctx use <context> [namespace]")
			os.Exit(1)
		}
		runUse(os.Args[2:])
	case "rm":
		if len(os.Args) < 3 {
			fmt.Println("Error: context name required")
			fmt.Println("Usage: ctx rm <context>")
			os.Exit(1)
		}
		runRemove(os.Args[2])
	case "--help", "-h":
		printUsage()
	default:
		fmt.Printf("Error: unknown command '%s'\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("ctx - A CLI tool for managing kubectl contexts and namespaces")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  ctx                              List all contexts (default)")
	fmt.Println("  ctx ls                           List all contexts")
	fmt.Println("  ctx use <context> [namespace]     Switch to a context and optionally set namespace")
	fmt.Println("  ctx rm <context>                 Remove a context")
	fmt.Println("  ctx --help                       Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  ctx                              # List all contexts")
	fmt.Println("  ctx ls                           # List all contexts")
	fmt.Println("  ctx use my-cluster               # Switch to context 'my-cluster'")
	fmt.Println("  ctx use my-cluster default        # Switch to context and set namespace to 'default'")
	fmt.Println("  ctx use \"complex context-name\" namespace  # Handle context names with spaces")
	fmt.Println("  ctx rm old-cluster               # Remove context 'old-cluster'")
	fmt.Println()
	fmt.Println("Note: For context names with spaces, use quotes:")
	fmt.Println("  ctx use \"complex context-name\" namespace")
}

func runList() {
	kc, err := kubeconfig.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	currentContext := kc.CurrentContext
	contexts := kc.GetContexts()

	if len(contexts) == 0 {
		fmt.Println("No contexts found")
		return
	}

	for _, ctx := range contexts {
		prefix := " "
		if ctx.Name == currentContext {
			prefix = "*"
		}
		fmt.Printf("%s %s", prefix, ctx.Name)

		if ctx.Name == currentContext && kc.CurrentNamespace() != "" {
			fmt.Printf(" (namespace: %s)", kc.CurrentNamespace())
		}
		fmt.Println()
	}
}

func runUse(args []string) {
	kc, err := kubeconfig.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	contextName := args[0]
	namespace := ""
	if len(args) > 1 {
		namespace = args[1]
	}

	if !kc.ContextExists(contextName) {
		fmt.Fprintf(os.Stderr, "Error: context '%s' does not exist\n", contextName)
		os.Exit(1)
	}

	if err := kc.SetCurrentContext(contextName, namespace); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := kc.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if namespace != "" {
		fmt.Printf("Switched to context '%s' with namespace '%s'\n", contextName, namespace)
	} else {
		fmt.Printf("Switched to context '%s'\n", contextName)
	}
}

func runRemove(contextName string) {
	kc, err := kubeconfig.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if !kc.ContextExists(contextName) {
		fmt.Fprintf(os.Stderr, "Error: context '%s' does not exist\n", contextName)
		os.Exit(1)
	}

	if kc.CurrentContext == contextName {
		fmt.Fprintf(os.Stderr, "Error: cannot remove currently active context '%s'\n", contextName)
		os.Exit(1)
	}

	if err := kc.RemoveContext(contextName); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := kc.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Removed context '%s'\n", contextName)
}