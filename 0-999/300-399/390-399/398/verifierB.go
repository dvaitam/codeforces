package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesB = `2 4 1 2 1 2 2 2 2 1
1 1 1 1
4 4 1 4 3 2 1 3 1 1
1 0
4 5 2 4 1 2 4 4 2 3 2 2
4 2 1 4 1 2
3 0
3 5 3 3 2 3 3 1 2 2 3 2
5 3 5 1 4 2 4 4
2 2 2 1 2 1
2 4 2 2 2 1 2 1 2 2
2 1 1 1
2 4 1 2 2 2 2 2 1 2
5 1 5 5
2 3 1 2 2 1 2 2
3 3 2 1 3 3 3 3
3 3 3 1 1 3 1 3
5 1 1 5
3 0
1 0
1 1 1 1
2 2 1 1 2 2
1 0
2 2 1 2 2 2
3 3 2 1 1 2 2 2
4 1 3 1
3 5 3 1 3 2 1 1 1 2 1 1
2 3 2 1 2 1 1 2
5 2 4 1 3 2
2 0
3 0
1 1 1 1
4 4 3 2 1 1 2 4 2 1
4 1 3 1
2 4 2 1 2 1 2 2 2 1
3 4 2 2 1 1 1 2 3 1
3 3 1 2 3 1 2 3
3 5 3 2 3 1 1 3 1 1 1 1
2 4 1 2 2 2 2 2 2 1
3 1 3 3
4 1 1 3
1 1 1 1
2 1 2 1
5 4 4 1 5 5 2 5 1 3
3 2 3 3 1 2
3 0
1 1 1 1
1 1 1 1
2 1 2 1
1 1 1 1
2 0
4 3 3 3 4 3 1 2
3 0
1 0
3 5 3 2 2 2 2 2 1 1 2 3
4 0
3 1 3 3
4 5 3 3 2 2 3 2 2 3 1 3
1 1 1 1
2 3 2 1 2 1 2 2
2 2 1 1 1 1
1 0
4 0
3 4 1 3 1 1 3 1 2 2
4 3 2 1 3 1 2 2
2 1 2 2
1 1 1 1
2 4 1 2 1 1 2 2 1 1
2 2 1 2 2 2
5 3 5 4 1 4 3 2
3 3 1 3 2 3 1 1
3 4 1 3 1 1 2 2 2 3
4 1 1 2
4 0
2 4 2 2 1 1 2 2 2 1
4 2 3 2 1 1
5 5 3 2 5 2 3 3 3 5 3 2
4 4 1 1 4 2 2 3 4 2
5 5 1 4 4 3 4 5 2 5 1 5
1 1 1 1
1 0
5 5 1 4 2 4 4 4 2 3 4 2
5 3 2 1 4 5 5 4
1 1 1 1
4 5 1 2 4 1 1 2 3 2 2 3
2 4 1 2 2 2 2 1 2 2
4 0
2 4 2 1 2 1 1 1 1 2
2 0
5 2 5 3 4 5
3 4 2 1 1 2 3 2 2 2
5 3 3 5 4 1 4 4
2 4 1 2 1 2 2 2 1 2
5 5 5 2 3 5 1 4 5 4 4 3
5 4 1 4 2 3 1 4 2 4
3 1 1 3
1 1 1 1
5 2 2 4 3 4
2 3 1 2 1 2 1 2
1 1 1 1
3 1 1 3
1 1 1 1
5 2 2 4 3 4
2 3 1 2 1 2 1 2
1 1 1 1`

func expected(n int, pairs [][2]int) string {
	a := make([]bool, n+1)
	b := make([]bool, n+1)
	for _, p := range pairs {
		if p[0] <= n {
			a[p[0]] = true
		}
		if p[1] <= n {
			b[p[1]] = true
		}
	}
	A, B := 0, 0
	for i := 1; i <= n; i++ {
		if a[i] {
			A++
		}
		if b[i] {
			B++
		}
	}
	f := make([][]float64, n+2)
	for i := range f {
		f[i] = make([]float64, n+2)
	}
	nn := float64(n)
	nn2 := nn * nn
	for i := n; i >= A; i-- {
		for j := n; j >= B; j-- {
			if i < n || j < n {
				denom := nn2 - float64(i*j)
				term1 := float64(n-i) * float64(j) * f[i+1][j]
				term2 := float64(i) * float64(n-j) * f[i][j+1]
				term3 := float64(n-i) * float64(n-j) * f[i+1][j+1]
				f[i][j] = (term1 + term2 + term3 + nn2) / denom
			}
		}
	}
	return fmt.Sprintf("%.9f", f[A][B])
}

func runCandidate(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcasesB))
	scanner.Buffer(make([]byte, 0, 1024), 1<<20)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 2 {
			fmt.Fprintf(os.Stderr, "bad test line %d\n", idx)
			os.Exit(1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad test line %d\n", idx)
			os.Exit(1)
		}
		m, err := strconv.Atoi(fields[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad test line %d\n", idx)
			os.Exit(1)
		}
		if len(fields) != 2+2*m {
			fmt.Fprintf(os.Stderr, "bad test line %d\n", idx)
			os.Exit(1)
		}
		pairs := make([][2]int, m)
		for i := 0; i < 2*m; i++ {
			v, _ := strconv.Atoi(fields[2+i])
			pairs[i/2][i%2] = v
		}
		want := expected(n, pairs)
		input := fmt.Sprintf("%d %d\n%s\n", n, m, strings.Join(fields[2:], " "))
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\n", idx, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
