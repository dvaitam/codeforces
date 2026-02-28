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
	n       int
	prices  []int
	queries [][2]int
}

func runProg(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(50) + 1
	prices := make([]int, n)
	for i := range prices {
		prices[i] = rng.Intn(100) + 1
	}
	q := rng.Intn(50) + 1
	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		l := rng.Intn(n) + 1
		r := rng.Intn(n-l+1) + l
		queries[i] = [2]int{l, r}
	}
	return testCase{n: n, prices: prices, queries: queries}
}

func renderCase(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", tc.n)
	for i, v := range tc.prices {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	fmt.Fprintf(&sb, "%d\n", len(tc.queries))
	for _, qr := range tc.queries {
		fmt.Fprintf(&sb, "%d %d\n", qr[0], qr[1])
	}
	return sb.String()
}

func brute(tc testCase) string {
	ans := make([]int64, len(tc.queries))
	for i, qr := range tc.queries {
		l, r := qr[0], qr[1]
		var sum int64
		mx := 0
		for pos := r; pos >= l; pos-- {
			if tc.prices[pos-1] > mx {
				mx = tc.prices[pos-1]
			}
			sum += int64(mx)
		}
		ans[i] = sum
	}

	var sb strings.Builder
	for i, v := range ans {
		if i > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		input := renderCase(tc)
		exp := brute(tc)

		got, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			fmt.Printf("input:\n%s", input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Printf("case %d mismatch\ninput:\n%s\nexpected:\n%s\n got:\n%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
