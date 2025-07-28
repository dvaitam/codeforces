package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func solveCase(x, y []int) int {
	n := len(x)
	diff := make([]int, n)
	for i := 0; i < n; i++ {
		diff[i] = y[i] - x[i]
	}
	sort.Ints(diff)
	l, r := 0, n-1
	ans := 0
	for l < r && diff[l] < 0 {
		if diff[l]+diff[r] >= 0 {
			ans++
			l++
			r--
		} else {
			l++
		}
	}
	if l < r {
		ans += (r - l + 1) / 2
	}
	return ans
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	type test struct {
		n    int
		x, y []int
	}
	tests := make([]test, 0, 120)
	// deterministic small tests
	tests = append(tests, test{n: 2, x: []int{1, 2}, y: []int{2, 2}})
	for len(tests) < 120 {
		n := rng.Intn(8) + 2
		x := make([]int, n)
		y := make([]int, n)
		for i := 0; i < n; i++ {
			x[i] = rng.Intn(20) + 1
			y[i] = rng.Intn(20) + 1
		}
		tests = append(tests, test{n: n, x: x, y: y})
	}

	for i, tc := range tests {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for idx, v := range tc.x {
			if idx > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		for idx, v := range tc.y {
			if idx > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expected := fmt.Sprintf("%d", solveCase(tc.x, tc.y))
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if out != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, out, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
