package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type test struct {
	input    string
	expected string
}

func isLeap(year int) bool {
	if year%400 == 0 {
		return true
	}
	if year%100 == 0 {
		return false
	}
	return year%4 == 0
}

func solve(input string) string {
	reader := strings.NewReader(input)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return ""
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	months := make([]int, 0, 9600)
	for y := 2000; y < 2800; y++ {
		feb := 28
		if isLeap(y) {
			feb = 29
		}
		months = append(months, 31, feb, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31)
	}
	found := false
	for i := 0; i+len(a) <= len(months); i++ {
		match := true
		for j := 0; j < n; j++ {
			if months[i+j] != a[j] {
				match = false
				break
			}
		}
		if match {
			found = true
			break
		}
	}
	if found {
		return "YES"
	}
	return "NO"
}

func generateTests() []test {
	rand.Seed(8992)
	var tests []test
	fixed := []string{
		"1\n31\n",
		"2\n31 28\n",
		"3\n30 31 30\n",
	}
	for _, in := range fixed {
		tests = append(tests, test{in, solve(in)})
	}
	for len(tests) < 100 {
		n := rand.Intn(24) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			val := 28 + rand.Intn(4)
			if val == 29 {
				sb.WriteString("29")
			} else {
				sb.WriteString(fmt.Sprintf("%d", val))
			}
		}
		sb.WriteByte('\n')
		inp := sb.String()
		tests = append(tests, test{inp, solve(inp)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("time limit")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		out, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if out != strings.TrimSpace(t.expected) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%sexpected:%s\n got:%s\n", i+1, t.input, t.expected, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
