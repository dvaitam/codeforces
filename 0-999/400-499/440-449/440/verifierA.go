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

func generate() []testCase {
	const T = 100
	rand.Seed(1)
	cases := make([]testCase, T)
	for i := 0; i < T; i++ {
		n := rand.Intn(100) + 2 // at least 2
		missing := rand.Intn(n) + 1
		var in strings.Builder
		fmt.Fprintf(&in, "%d\n", n)
		first := true
		for j := 1; j <= n; j++ {
			if j == missing {
				continue
			}
			if !first {
				in.WriteByte(' ')
			}
			fmt.Fprintf(&in, "%d", j)
			first = false
		}
		in.WriteByte('\n')
		cases[i] = testCase{
			in:  in.String(),
			out: fmt.Sprintf("%d\n", missing),
		}
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generate()
	for idx, tc := range cases {
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(tc.in)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("case %d: runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != strings.TrimSpace(tc.out) {
			fmt.Printf("case %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx+1, tc.out, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
