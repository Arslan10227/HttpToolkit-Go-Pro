package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/mcp/ctlclient"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: htk-ctl <status|operations|execute> [args-json]")
		os.Exit(2)
	}
	client, err := ctlclient.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	switch os.Args[1] {
	case "status":
		out, err := client.Status()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		printJSON(out)
	case "operations":
		out, err := client.Operations()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		printJSON(out)
	case "execute":
		if len(os.Args) < 4 {
			fmt.Fprintln(os.Stderr, "usage: htk-ctl execute <name> <args-json>")
			os.Exit(2)
		}
		var args map[string]any
		if err := json.Unmarshal([]byte(os.Args[3]), &args); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		out, err := client.Execute(os.Args[2], args)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		printJSON(out)
	default:
		fmt.Fprintln(os.Stderr, "unknown command:", os.Args[1])
		os.Exit(2)
	}
}

func printJSON(v any) {
	data, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(data))
}
