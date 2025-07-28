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

func expectedA(input string) string {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) == 0 {
		return ""
	}
	var t int
	fmt.Sscanf(lines[0], "%d", &t)
	var out strings.Builder
	for i := 0; i < t; i++ {
		var a, b int
		fmt.Sscanf(lines[1+i], "%d+%d", &a, &b)
		out.WriteString(fmt.Sprintf("%d\n", a+b))
	}
	return strings.TrimSpace(out.String())
}

func genTestsA() []string {
	rand.Seed(1)
	tests := make([]string, 0, 100)
	for len(tests) < 100 {
		t := rand.Intn(100) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", t))
		for i := 0; i < t; i++ {
			a := rand.Intn(10)
			b := rand.Intn(10)
			sb.WriteString(fmt.Sprintf("%d+%d\n", a, b))
		}
		tests = append(tests, sb.String())
	}
	return tests
}

func runBinary(bin string, input string) (string, error) {
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
		fmt.Fprintf(os.Stderr, "Usage: go run verifierA.go <binary>\n")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsA()
	for i, t := range tests {
		want := expectedA(t)
		got, err := runBinary(bin, t)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("Test %d failed.\nInput:\n%s\nExpected:\n%s\nGot:\n%s\n", i+1, t, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
	_ = time.Now()
}
