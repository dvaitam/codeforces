package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, _ := filepath.Glob("tests/B/*.in")
	total := 0
	for _, inFile := range cases {
		outFile := strings.TrimSuffix(inFile, ".in") + ".out"
		inData, _ := ioutil.ReadFile(inFile)
		want, _ := ioutil.ReadFile(outFile)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(inData)
		got, err := cmd.Output()
		if err != nil {
			fmt.Printf("%s: runtime error: %v\n", inFile, err)
			os.Exit(1)
		}
		if !bytes.Equal(bytes.TrimSpace(got), bytes.TrimSpace(want)) {
			fmt.Printf("%s: wrong answer\nwant: %q\ngot:  %q\n", inFile, want, got)
			os.Exit(1)
		}
		total++
	}
	fmt.Printf("OK %d cases\n", total)
}
