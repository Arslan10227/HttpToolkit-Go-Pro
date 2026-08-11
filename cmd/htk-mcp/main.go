package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/mcp"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/mcp/ctlclient"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "stdio" || os.Args[1] == "--stdio") {
		if err := mcp.RunStdio(os.Stderr); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	client, err := ctlclient.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	ready, err := client.Status()
	if err != nil || ready["ready"] != true {
		fmt.Fprintln(os.Stderr, "UI bridge not ready")
		os.Exit(1)
	}
	ops, err := client.Operations()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	out := map[string]any{
		"protocol": "httptoolkit-mcp",
		"version":  "1.0.0",
		"tools":    ops,
	}
	data, _ := json.MarshalIndent(out, "", "  ")
	fmt.Println(string(data))
}
