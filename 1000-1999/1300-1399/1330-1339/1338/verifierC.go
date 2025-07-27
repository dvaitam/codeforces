package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Test struct {
	N   int64  `json:"n"`
	Out string `json:"out"`
}

func runCase(bin string, t Test) (string, error) {
	input := fmt.Sprintf("1\n%d\n", t.N)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testsC.json")
	if err != nil {
		fmt.Println("cannot read testsC.json:", err)
		return
	}
	var tests []Test
	if err := json.Unmarshal(data, &tests); err != nil {
		fmt.Println("bad tests json:", err)
		return
	}
	passed := 0
	for i, t := range tests {
		got, err := runCase(bin, t)
		got = strings.TrimSpace(got)
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n", i+1, err)
			continue
		}
		if got == strings.TrimSpace(t.Out) {
			passed++
		} else {
			fmt.Printf("case %d failed: expected %s got %s\n", i+1, t.Out, got)
		}
	}
	fmt.Printf("passed %d/%d\n", passed, len(tests))
}
