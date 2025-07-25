package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func solveCase(n, k int, p int64, a, b []int64) int64 {
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
	const inf int64 = 1 << 60
	dp := make([]int64, k+1)
	for i := range dp {
		dp[i] = 0
	}
	for i := 1; i <= n; i++ {
		ndp := make([]int64, k+1)
		for j := range ndp {
			ndp[j] = inf
		}
		for j := 1; j <= k; j++ {
			cost := abs(a[i-1]-b[j-1]) + abs(b[j-1]-p)
			val := max(dp[j-1], cost)
			if val < ndp[j] {
				ndp[j] = val
			}
			if ndp[j-1] < ndp[j] {
				ndp[j] = ndp[j-1]
			}
		}
		dp = ndp
	}
	return dp[k]
}

func run(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	data, err := os.ReadFile("testcasesA.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not read testcasesA.txt:", err)
		os.Exit(1)
	}

	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Fprintln(os.Stderr, "invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())

	type caseData struct {
		n int
		k int
		p int64
		a []int64
		b []int64
	}
	cases := make([]caseData, t)
	expected := make([]int64, t)

	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Fprintln(os.Stderr, "bad file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		k, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		p64, _ := strconv.ParseInt(scan.Text(), 10, 64)
		a := make([]int64, n)
		for j := 0; j < n; j++ {
			scan.Scan()
			val, _ := strconv.ParseInt(scan.Text(), 10, 64)
			a[j] = val
		}
		b := make([]int64, k)
		for j := 0; j < k; j++ {
			scan.Scan()
			val, _ := strconv.ParseInt(scan.Text(), 10, 64)
			b[j] = val
		}
		cases[i] = caseData{n, k, p64, a, b}
		expected[i] = solveCase(n, k, p64, a, b)
	}

	for i, c := range cases {
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d %d %d\n", c.n, c.k, c.p)
		for j, v := range c.a {
			if j > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
		for j, v := range c.b {
			if j > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')

		outStr, err := run(bin, input.Bytes())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := strconv.ParseInt(strings.TrimSpace(outStr), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: bad output\n", i+1)
			os.Exit(1)
		}
		if got != expected[i] {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", t)
}
