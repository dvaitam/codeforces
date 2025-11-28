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

func runCandidate(bin, input string) (string, error) {
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

func expected(n int, x, t []int) float64 {
	minDiff := int64(1<<63 - 1)
	maxSum := int64(-1 << 63)
	for i := 0; i < n; i++ {
		diff := int64(x[i] - t[i])
		sum := int64(x[i] + t[i])
		if diff < minDiff {
			minDiff = diff
		}
		if sum > maxSum {
			maxSum = sum
		}
	}
	return float64(minDiff+maxSum) * 0.5
}

func isClose(a, b float64) bool {
	const eps = 1e-6
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	if diff <= eps {
		return true
	}
	if b < 0 {
		b = -b
	}
	if a < 0 {
		a = -a
	}
	maxAB := a
	if b > maxAB {
		maxAB = b
	}
	return diff <= eps*maxAB
}

func generateCase(rng *rand.Rand) (string, int, []int, []int) {
	n := rng.Intn(100) + 1
	xs := make([]int, n)
	ts := make([]int, n)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(strconv.Itoa(n))
	sb.WriteString("\n")
	
	for i := 0; i < n; i++ {
		xs[i] = rng.Intn(100000000)
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(strconv.Itoa(xs[i]))
	}
	sb.WriteString("\n")
	
	for i := 0; i < n; i++ {
		ts[i] = rng.Intn(100000000)
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(strconv.Itoa(ts[i]))
	}
	sb.WriteString("\n")
	return sb.String(), n, xs, ts
}

func main() {
	arg := ""
	if len(os.Args) == 2 {
		arg = os.Args[1]
	} else if len(os.Args) == 3 && os.Args[1] == "--" {
		arg = os.Args[2]
	} else {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := arg

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	for idx := 1; idx <= 100; idx++ {
		input, n, xs, ts := generateCase(rng)
		expect := expected(n, xs, ts)
		
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx, err)
			os.Exit(1)
		}
		if got == "" {
			fmt.Printf("case %d failed: empty output\n", idx)
			os.Exit(1)
		}
		gotVal, err := strconv.ParseFloat(strings.Fields(got)[0], 64)
		if err != nil {
			fmt.Printf("case %d failed: unable to parse output %q\n", idx, got)
			os.Exit(1)
		}
		if !isClose(gotVal, expect) {
			fmt.Printf("case %d failed: expected %.6f got %.8f\n", idx, expect, gotVal)
			fmt.Printf("Input:\n%s\n", input)
			os.Exit(1)
		}
	}
	fmt.Printf("All 100 tests passed\n")
}