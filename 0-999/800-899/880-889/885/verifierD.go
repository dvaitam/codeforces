package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func generateTests() []string {
	r := rand.New(rand.NewSource(45))
	tests := make([]string, 100)
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	for i := 0; i < 100; i++ {
		length := r.Intn(10) + 1 // 1..10
		var sb strings.Builder
		for j := 0; j < length; j++ {
			sb.WriteRune(letters[r.Intn(len(letters))])
		}
		tests[i] = sb.String()
	}
	return tests
}

func expected(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	tests := generateTests()
	for i, s := range tests {
		input := fmt.Sprintf("%s\n", s)
		exp := expected(s)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("Test %d failed. Input: %sExpected %s got %s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
