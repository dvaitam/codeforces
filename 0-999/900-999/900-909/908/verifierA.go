package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func runCommand(cmd []string, input string) (string, error) {
	c := exec.Command(cmd[0], cmd[1:]...)
	c.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	c.Stdout = &out
	err := c.Run()
	return strings.TrimSpace(out.String()), err
}

func expectedOutput(input string) string {
	s := strings.TrimSpace(input)
	count := 0
	for _, ch := range s {
		if strings.ContainsRune("aeiou13579", ch) {
			count++
		}
	}
	return fmt.Sprintf("%d", count)
}

func generateTests() []string {
	rand.Seed(42)
	tests := make([]string, 100)
	chars := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	for i := 0; i < 100; i++ {
		n := rand.Intn(50) + 1
		var sb strings.Builder
		for j := 0; j < n; j++ {
			sb.WriteRune(chars[rand.Intn(len(chars))])
		}
		tests[i] = sb.String()
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go <binary>")
		os.Exit(1)
	}
	candidate := os.Args[1]
	useGoRun := strings.HasSuffix(candidate, ".go")

	tests := generateTests()
	for idx, input := range tests {
		var got string
		var err error
		if useGoRun {
			got, err = runCommand([]string{"go", "run", candidate}, input+"\n")
		} else {
			err = os.Chmod(candidate, 0755)
			got, err = runCommand([]string{candidate}, input+"\n")
		}
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		expect := expectedOutput(input)
		if got != expect {
			fmt.Printf("test %d failed: input=%s expected=%s got=%s\n", idx+1, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
