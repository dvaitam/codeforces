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

func solve(input string) string {
	reader := strings.NewReader(input)
	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return ""
	}
	a := make([]int, n)
	for i := range a {
		fmt.Fscan(reader, &a[i])
	}
	changes := 0
	for i := 0; i < k; i++ {
		cnt1, cnt2 := 0, 0
		for j := i; j < n; j += k {
			if a[j] == 1 {
				cnt1++
			} else {
				cnt2++
			}
		}
		groupSize := cnt1 + cnt2
		if cnt1 > cnt2 {
			changes += groupSize - cnt1
		} else {
			changes += groupSize - cnt2
		}
	}
	return fmt.Sprintf("%d", changes)
}

func generateTests() []test {
	rand.Seed(42)
	var tests []test
	fixed := []struct {
		n   int
		k   int
		arr []int
	}{
		{1, 1, []int{1}},
		{3, 1, []int{1, 2, 1}},
		{4, 2, []int{1, 2, 1, 2}},
	}
	for _, f := range fixed {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", f.n, f.k)
		for i, v := range f.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		inp := sb.String()
		tests = append(tests, test{inp, solve(inp)})
	}
	for len(tests) < 100 {
		n := rand.Intn(20) + 1
		k := rand.Intn(n) + 1
		arr := make([]int, n)
		for i := range arr {
			arr[i] = rand.Intn(2) + 1
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, k)
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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
		if strings.TrimSpace(out) != strings.TrimSpace(t.expected) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%sexpected:%s\n got:%s\n", i+1, t.input, t.expected, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
