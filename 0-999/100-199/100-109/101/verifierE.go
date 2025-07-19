package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type testCaseE struct {
	input string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := []testCaseE{
		{input: "1 1 10\n3\n4\n"},
		{input: "2 2 5\n1 2\n3 4\n"},
	}
	for i, t := range tests {
		expect := solveE(strings.NewReader(t.input))
		out, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("test %d: execution failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed: expected %q got %q\n", i+1, strings.TrimSpace(expect), strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

var (
	n, m, p int
	X, Y    []int
	F, G    []int
	ans     []byte
	sum     int64
)

func f(a, b int) int { return (X[a] + Y[b]) % p }

func solve(x0, y0, x1, y1 int) {
	if x0 == x1 {
		sum += int64(f(x0, y0))
		for j := y0 + 1; j <= y1; j++ {
			sum += int64(f(x0, j))
			ans = append(ans, 'S')
		}
		return
	}
	for j := y0; j <= y1; j++ {
		F[j] = 0
		G[j] = 0
	}
	mid := (x0 + x1) >> 1
	for i := x0; i <= mid; i++ {
		F[y0] += f(i, y0)
		for j := y0 + 1; j <= y1; j++ {
			if F[j-1] > F[j] {
				F[j] = F[j-1]
			}
			F[j] += f(i, j)
		}
	}
	for i := x1; i >= mid+1; i-- {
		G[y1] += f(i, y1)
		for j := y1 - 1; j >= y0; j-- {
			if G[j+1] > G[j] {
				G[j] = G[j+1]
			}
			G[j] += f(i, j)
		}
	}
	bst := y1
	for j := y0; j < y1; j++ {
		if F[j]+G[j] > F[bst]+G[bst] {
			bst = j
		}
	}
	solve(x0, y0, mid, bst)
	ans = append(ans, 'C')
	solve(mid+1, bst, x1, y1)
}

func solveE(r io.Reader) string {
	in := bufio.NewReader(r)
	fmt.Fscan(in, &n, &m, &p)
	X = make([]int, n)
	Y = make([]int, m)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &X[i])
	}
	for j := 0; j < m; j++ {
		fmt.Fscan(in, &Y[j])
	}
	F = make([]int, m)
	G = make([]int, m)
	ans = make([]byte, 0, n+m)
	sum = 0
	solve(0, 0, n-1, m-1)
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%d\n", sum)
	if len(ans) > 0 {
		buf.Write(ans)
	}
	buf.WriteByte('\n')
	return buf.String()
}
