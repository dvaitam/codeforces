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

func genTestsG1() []string {
	rand.Seed(1)
	tests := make([]string, 100)
	for i := 0; i < 100; i++ {
		k := rand.Intn(5) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", k)
		for j := 0; j < 6; j++ {
			fmt.Fprintf(&sb, "%d ", rand.Intn(10))
		}
		sb.WriteByte('\n')
		fmt.Fprintf(&sb, "1\n")
		fmt.Fprintf(&sb, "%d\n", rand.Intn(200))
		tests[i] = sb.String()
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref := "1428G1.go"
	tests := genTestsG1()
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
			fmt.Fprintf(os.Stderr, "test %d: expected %q got %q\n", i+1, expect, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
