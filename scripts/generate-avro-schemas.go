package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	schemaDir := "lib/go/kafka/topics"
	files, err := ioutil.ReadDir(schemaDir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		os.Exit(1)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".avsc" {
			fmt.Println("Generating struct for:", file.Name())
			outFileName := file.Name()[:len(file.Name())-5] + ".go" // Remove .avsc and add .go
			cmd := exec.Command("avrogen", "-pkg", "avro", "-o", filepath.Join(schemaDir, outFileName), "-tags", "json:camel", filepath.Join(schemaDir, file.Name()))
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				fmt.Println("Error generating struct for", file.Name(), ":", err)
			}
		}
	}
}
