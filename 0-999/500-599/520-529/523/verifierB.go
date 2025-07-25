package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strings"
)

func runCandidate(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expectedOutput(n, t int, c float64, arr []int, queries []int) string {
	realAvg := make([]float64, n+1)
	approx := make([]float64, n+1)
	var sum int64
	ft := float64(t)
	for i := 1; i <= n; i++ {
		sum += int64(arr[i-1])
		if i == t {
			realAvg[i] = float64(sum) / ft
		}
		if i > t {
			sum -= int64(arr[i-1-t])
			realAvg[i] = float64(sum) / ft
		}
		approx[i] = (approx[i-1] + float64(arr[i-1])/ft) / c
	}
	var b strings.Builder
	for j, x := range queries {
		r := realAvg[x]
		ap := approx[x]
		errRel := math.Abs(ap-r) / r
		fmt.Fprintf(&b, "%.5f %.5f %.5f", r, ap, errRel)
		if j+1 < len(queries) {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var tcases int
	if _, err := fmt.Fscan(reader, &tcases); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read test count: %v\n", err)
		os.Exit(1)
	}
	for caseNum := 1; caseNum <= tcases; caseNum++ {
		var n, t int
		var c float64
		fmt.Fscan(reader, &n, &t, &c)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		var m int
		fmt.Fscan(reader, &m)
		queries := make([]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(reader, &queries[i])
		}
		var input strings.Builder
		fmt.Fprintf(&input, "%d %d %.6f\n", n, t, c)
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprintf(&input, "%d", arr[i])
		}
		input.WriteByte('\n')
		fmt.Fprintf(&input, "%d\n", m)
		for i := 0; i < m; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprintf(&input, "%d", queries[i])
		}
		input.WriteByte('\n')
		want := expectedOutput(n, t, c, arr, queries)
		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", caseNum, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(want) {
			fmt.Printf("case %d failed\nexpected:\n%s\n\ngot:\n%s\n", caseNum, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", tcases)
}
