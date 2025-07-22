package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type test struct {
	input    string
	expected string
}

func solve(input string) string {
	reader := strings.NewReader(input)
	var n, x int
	if _, err := fmt.Fscan(reader, &n, &x); err != nil {
		return ""
	}
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &b[i])
	}
	sort.Ints(a)
	sort.Ints(b)
	best := 1
	i, j := 0, n-1
	worst := 0
	for i < n && j >= 0 {
		if a[i]+b[j] >= x {
			worst++
			i++
			j--
		} else {
			i++
		}
	}
	return fmt.Sprintf("%d %d", best, worst)
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(45))
	var tests []test
	fixed := []string{
		"1 1\n1\n0\n",
		"3 2\n1 1 1\n1 1 1\n",
	}
	for _, f := range fixed {
		tests = append(tests, test{f, solve(f)})
	}
	for len(tests) < 100 {
		n := rng.Intn(5) + 1
		x := rng.Intn(10)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, x))
		for i := 0; i < n; i++ {
			fmt.Fprintf(&sb, "%d ", rng.Intn(10))
		}
		sb.WriteByte('\n')
		for i := 0; i < n; i++ {
			fmt.Fprintf(&sb, "%d ", rng.Intn(10))
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
