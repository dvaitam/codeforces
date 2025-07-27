package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genTestsB() []string {
	rand.Seed(1)
	tests := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(8) + 2
		var sb strings.Builder
		for j := 0; j < n; j++ {
			r := rand.Intn(3)
			if r == 0 {
				sb.WriteByte('<')
			} else if r == 1 {
				sb.WriteByte('>')
			} else {
				sb.WriteByte('-')
			}
		}
		s := sb.String()
		tests[i] = fmt.Sprintf("1\n%d\n%s\n", n, s)
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref := "1428B.go"
	tests := genTestsB()
	for i, input := range tests {
		expect, err := run(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "test %d: expected %s got %s\n", i+1, expect, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
