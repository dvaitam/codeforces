package main

import (
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
	reader := strings.NewReader(input)
	var t int
	fmt.Fscan(reader, &t)
	var out strings.Builder
	const mask30 int64 = (1 << 30) - 1
	for ; t > 0; t-- {
		var n int
		var x int64
		fmt.Fscan(reader, &n, &x)
		mask := (^x) & mask30
		px := int64(0)
		base := int64(0)
		segments := 0
		for i := 0; i < n; i++ {
			var v int64
			fmt.Fscan(reader, &v)
			px ^= v & mask
			if px == base {
				segments++
				base = px
			}
		}
		if px == base {
			out.WriteString(fmt.Sprintf("%d\n", segments))
		} else {
			out.WriteString("-1\n")
		}
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
