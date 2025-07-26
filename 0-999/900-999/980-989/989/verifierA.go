package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	in  string
	out string
}

func expected(s string) string {
	n := len(s)
	for i := 0; i+2 < n; i++ {
		hasA, hasB, hasC := false, false, false
		for j := 0; j < 3; j++ {
			switch s[i+j] {
			case 'A':
				hasA = true
			case 'B':
				hasB = true
			case 'C':
				hasC = true
			}
		}
		if hasA && hasB && hasC {
			return "Yes"
		}
	}
	return "No"
}

func generate() []testCase {
	const T = 100
	rand.Seed(1)
	cases := make([]testCase, T)
	for i := 0; i < T; i++ {
		n := rand.Intn(100) + 1
		var sb strings.Builder
		for j := 0; j < n; j++ {
			ch := "ABC."[rand.Intn(4)]
			sb.WriteByte(ch)
		}
		s := sb.String()
		cases[i] = testCase{
			in:  s + "\n",
			out: expected(s),
		}
	}
	return cases
}

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
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generate()
	for idx, tc := range cases {
		got, err := run(bin, tc.in)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != strings.TrimSpace(tc.out) {
			fmt.Printf("case %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx+1, tc.out, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
