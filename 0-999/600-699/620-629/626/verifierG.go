package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Test struct {
	N       int      `json:"n"`
	T       int      `json:"t"`
	Q       int      `json:"q"`
	Pi      []int    `json:"pi"`
	Li      []int    `json:"li"`
	Updates [][]int  `json:"updates"`
	Out     []string `json:"out"`
}

func runCase(bin string, t Test) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", t.N, t.T, t.Q))
	for i, v := range t.Pi {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	for i, v := range t.Li {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	for _, u := range t.Updates {
		sb.WriteString(fmt.Sprintf("%d %d\n", u[0], u[1]))
	}
	input := sb.String()
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), err
	}
	return string(out), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testsG.json")
	if err != nil {
		fmt.Println("cannot read testsG.json:", err)
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
		gotLines := strings.Fields(got)
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n", i+1, err)
			continue
		}
		ok := len(gotLines) == len(t.Out)
		if ok {
			for j, exp := range t.Out {
				if strings.TrimSpace(gotLines[j]) != strings.TrimSpace(exp) {
					ok = false
					break
				}
			}
		}
		if ok {
			passed++
		} else {
			fmt.Printf("case %d failed\n", i+1)
		}
	}
	fmt.Printf("passed %d/%d\n", passed, len(tests))
}
