package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const solutionG = "0-999/900-999/900-909/908/908G.go"

func runCmd(cmd []string, input string) (string, error) {
	c := exec.Command(cmd[0], cmd[1:]...)
	c.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	c.Stdout = &out
	err := c.Run()
	return strings.TrimSpace(out.String()), err
}

func generateTests() []string {
	rand.Seed(42)
	var tests []string
	for i := 0; i < 100; i++ {
		length := rand.Intn(10) + 1
		var sb strings.Builder
		for j := 0; j < length; j++ {
			sb.WriteByte(byte('0' + rand.Intn(10)))
		}
		sb.WriteByte('\n')
		tests = append(tests, sb.String())
	}
	return tests
}

func runCandidate(path string, input string) (string, error) {
	if strings.HasSuffix(path, ".go") {
		return runCmd([]string{"go", "run", path}, input)
	}
	os.Chmod(path, 0755)
	return runCmd([]string{path}, input)
}

func expectedOutput(input string) (string, error) {
	return runCmd([]string{"go", "run", solutionG}, input)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierG.go <binary>")
		os.Exit(1)
	}
	cand := os.Args[1]
	tests := generateTests()
	for idx, t := range tests {
		exp, err := expectedOutput(t)
		if err != nil {
			fmt.Printf("reference failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := runCandidate(cand, t)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("test %d failed\ninput:%sexpected:%s got:%s\n", idx+1, t, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
