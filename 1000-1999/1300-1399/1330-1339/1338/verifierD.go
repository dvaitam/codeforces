package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Test struct {
	N     int     `json:"n"`
	Edges [][]int `json:"edges"`
	Out   string  `json:"out"`
}

func runCase(bin string, t Test) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t.N))
	for _, e := range t.Edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testsD.json")
	if err != nil {
		fmt.Println("cannot read testsD.json:", err)
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
