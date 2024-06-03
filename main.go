package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

const UrlFlag = "url"
const OutputFlag = "out"
const OutputProtoFlag = "out-gen"

var urlPtr *string
var outputPtr *string
var outputGeneratePtr *string

func init() {
	urlPtr = flag.String(UrlFlag, "", "enter the link to the proto file")
	outputPtr = flag.String(OutputFlag, "", "enter the path to the file where the proto file will be created")
	outputGeneratePtr = flag.String(OutputProtoFlag, "", "enter the path where the file will be generated")

	flag.Parse()
	if urlPtr == nil || *urlPtr == "" {
		os.Exit(-1)
	}
	if outputPtr == nil || *outputPtr == "" {
		os.Exit(-1)
	}
	if outputGeneratePtr == nil || *outputGeneratePtr == "" {
		os.Exit(-1)
	}
}

func main() {
	f, err := os.Create(*outputPtr)
	if err != nil {
		fmt.Println("ERROR (create file): ", err)
		return
	}
	defer f.Close()

	resp, err := http.Get(*urlPtr)
	if err != nil {
		fmt.Println("ERROR (get from url): ", err)
		return
	}
	defer resp.Body.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		fmt.Println("ERROR (copy data to file): ", err)
		return
	}

	cmd := exec.Command("protoc",
		fmt.Sprintf("--go_out=%s", *outputGeneratePtr),
		fmt.Sprintf("--go_opt=paths=source_relative"),
		fmt.Sprintf("--go-grpc_out=%s", *outputGeneratePtr),
		fmt.Sprintf("--go-grpc_opt=paths=source_relative"),
		fmt.Sprintf(*outputPtr),
	)
	err = cmd.Run()

	if err != nil {
		fmt.Println("ERROR: (exec protoc) ", err)
		return
	}
	fmt.Println("files generated")
}
