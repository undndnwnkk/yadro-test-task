package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

const serverURL = "http://localhost:8080/dns"

func main() {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addIP := addCmd.String("ip", "", "IP address to add")
	removeCmd := flag.NewFlagSet("remove", flag.ExitOnError)
	removeIP := removeCmd.String("ip", "", "IP address to remove")

	if len(os.Args) < 2 {
		printHelp()
		return
	}

	switch os.Args[1] {
	case "list":
		listServers()
	case "add":
		_ = addCmd.Parse(os.Args[2:])
		if *addIP == "" {
			addCmd.Usage()
			return
		}
		sendRequest("POST", *addIP)
	case "remove":
		_ = removeCmd.Parse(os.Args[2:])
		if *removeIP == "" {
			removeCmd.Usage()
			return
		}
		sendRequest("DELETE", *removeIP)
	case "--help", "-h":
		printHelp()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printHelp()
	}
}

func printHelp() {
	fmt.Println("Usage: dns-client <command> [arguments]")
	fmt.Println("Commands:")
	fmt.Println("	list : Get list of DNS servers")
	fmt.Println("	add --ip IP : Add a DNS server")
	fmt.Println(" 	remove --ip IP : Remove a DNS server")
}

func listServers() {
	response, err := http.Get(serverURL)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	defer func() {
		closeErr := response.Body.Close()
		if err == nil {
			err = closeErr
		}
	}()

	_, _ = io.Copy(os.Stdout, response.Body)
	fmt.Println()
}

func sendRequest(method, ip string) {
	var request *http.Request
	if method == "POST" {
		body, _ := json.Marshal(map[string]string{"ip": ip})

		request, _ = http.NewRequest(method, serverURL, bytes.NewBuffer(body))
	} else {
		request, _ = http.NewRequest(method, serverURL+"?ip="+ip, nil)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if response.StatusCode >= 400 {
		b, _ := io.ReadAll(response.Body)
		fmt.Printf("Fail:%s\n", string(b))
	} else {
		fmt.Println("Success")
	}
}
