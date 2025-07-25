package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	input  string
	expect []int
}

func brute(n int, arr []int) []int {
	ans := make([]int, n+1)
	for l := 0; l < n; l++ {
		freq := make([]int, n+1)
		maxColor, maxCount := 0, 0
		for r := l; r < n; r++ {
			c := arr[r]
			freq[c]++
			if freq[c] > maxCount || (freq[c] == maxCount && c < maxColor) {
				maxCount = freq[c]
				maxColor = c
			}
			ans[maxColor]++
		}
	}
	return ans[1:]
}

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(7) + 1 // 1..7
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(n) + 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", arr[i])
	}
	sb.WriteByte('\n')
	expect := brute(n, arr)
	return testCase{sb.String(), expect}
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseOutput(out string, n int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != n {
		return nil, fmt.Errorf("expected %d numbers, got %d", n, len(fields))
	}
	res := make([]int, n)
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("invalid integer: %v", err)
		}
		res[i] = v
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const cases = 100
	for i := 0; i < cases; i++ {
		c := genCase(rng)
		out, err := run(bin, c.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, c.input)
			os.Exit(1)
		}
		got, err := parseOutput(out, len(c.expect))
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%soutput:\n%s", i+1, err, c.input, out)
			os.Exit(1)
		}
		for j, v := range got {
			if v != c.expect[j] {
				fmt.Fprintf(os.Stderr, "case %d failed at position %d: expected %d got %d\ninput:\n%s", i+1, j+1, c.expect[j], v, c.input)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
