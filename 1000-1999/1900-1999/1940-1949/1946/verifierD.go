package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type test struct {
	input    string
	expected string
}

func solve(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var t int
	fmt.Fscan(reader, &t)

	var out strings.Builder

	for ; t > 0; t-- {
		var n, x int
		fmt.Fscan(reader, &n, &x)

		p := make([]int, n+1)
		currentXor := 0
		for i := 0; i < n; i++ {
			var v int
			fmt.Fscan(reader, &v)
			currentXor ^= v
			p[i+1] = currentXor
		}

		finalXor := p[n]
		maxK := -1

		check := func(mask int) {
			if (finalXor & ^mask) != 0 {
				return
			}

			cnt := 0
			for i := 1; i <= n; i++ {
				if (p[i] & ^mask) == 0 {
					cnt++
				}
			}
			if cnt > maxK {
				maxK = cnt
			}
		}

		check(x)
		for b := 0; b < 30; b++ {
			if (x>>b)&1 == 1 {
				mask := (x &^ ((1 << (b + 1)) - 1)) | ((1 << b) - 1)
				check(mask)
			}
		}

		out.WriteString(fmt.Sprintf("%d\n", maxK))
	}

	return strings.TrimSpace(out.String())
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(49))
	var tests []test
	fixed := []string{
		"1\n1 0\n0\n",
	}
	for _, f := range fixed {
		tests = append(tests, test{f, solve(f)})
	}
	for len(tests) < 100 {
		n := rng.Intn(5) + 1
		x := rng.Int63n(8)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", n, x))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", rng.Int63n(8))
		}
		sb.WriteByte('\n')
		inp := sb.String()
		tests = append(tests, test{inp, solve(inp)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
