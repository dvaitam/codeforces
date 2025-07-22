package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
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
	fmt.Fscan(reader, &n)
	stations := make([]struct {
		x   int
		idx int
	}, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &stations[i].x)
		stations[i].idx = i + 1
	}
	fmt.Fscan(reader, &k)
	sort.Slice(stations, func(i, j int) bool { return stations[i].x < stations[j].x })
	var sumX int64
	var S int64
	for t := 0; t < k; t++ {
		xi := int64(stations[t].x)
		sumX += xi
		S += xi * int64(2*t-k+1)
	}
	best := S
	bestL := 0
	for l := 0; l+k < n; l++ {
		leftX := int64(stations[l].x)
		rightX := int64(stations[l+k].x)
		S = S + int64(k+1)*leftX - 2*sumX + int64(k-1)*rightX
		sumX = sumX - leftX + rightX
		if S < best {
			best = S
			bestL = l + 1
		}
	}
	var out bytes.Buffer
	for t := 0; t < k; t++ {
		if t > 0 {
			out.WriteByte(' ')
		}
		fmt.Fprint(&out, stations[bestL+t].idx)
	}
	return out.String()
}

func generateTests() []test {
	rand.Seed(46)
	var tests []test
	fixed := []string{
		"3\n1 2 3\n2\n",
		"4\n10 20 30 40\n1\n",
	}
	for _, in := range fixed {
		tests = append(tests, test{in, solve(in)})
	}
	for len(tests) < 100 {
		n := rand.Intn(8) + 2
		k := rand.Intn(n-1) + 1
		xs := make([]int, n)
		for i := 0; i < n; i++ {
			xs[i] = rand.Intn(50)
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for i, v := range xs {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		fmt.Fprintf(&sb, "%d\n", k)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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
