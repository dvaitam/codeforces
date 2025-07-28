package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

type caseC struct {
	n int
	h []int64
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func solveC(n int, h []int64) int64 {
	if n%2 == 1 {
		var ans int64
		for i := 1; i < n-1; i += 2 {
			need := max64(h[i-1], h[i+1]) + 1
			if need > h[i] {
				ans += need - h[i]
			}
		}
		return ans
	}
	m := n - 2
	pairs := m / 2
	if pairs == 0 {
		return 0
	}
	cost := make([]int64, m+1)
	for k := 1; k <= m; k++ {
		idx := k
		need := max64(h[idx-1], h[idx+1]) + 1
		if need > h[idx] {
			cost[k] = need - h[idx]
		}
	}
	dp0 := make([]int64, pairs+1)
	dp1 := make([]int64, pairs+1)
	dp0[1] = cost[1]
	dp1[1] = cost[2]
	for j := 2; j <= pairs; j++ {
		dp0[j] = cost[2*j-1] + dp0[j-1]
		prev := dp0[j-1]
		if dp1[j-1] < prev {
			prev = dp1[j-1]
		}
		dp1[j] = cost[2*j] + prev
	}
	ans := dp0[pairs]
	if dp1[pairs] < ans {
		ans = dp1[pairs]
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const T = 100
	tests := make([]caseC, T)
	expected := make([]int64, T)
	var input bytes.Buffer
	fmt.Fprintln(&input, T)
	for i := 0; i < T; i++ {
		n := rng.Intn(8) + 3
		h := make([]int64, n)
		for j := range h {
			h[j] = rand.Int63n(10) + 1
		}
		tests[i] = caseC{n: n, h: h}
		fmt.Fprintln(&input, n)
		for j, v := range h {
			if j+1 == len(h) {
				fmt.Fprintf(&input, "%d\n", v)
			} else {
				fmt.Fprintf(&input, "%d ", v)
			}
		}
		expected[i] = solveC(n, h)
	}
	cmd := exec.Command(binary)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error running binary:", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	scanner.Split(bufio.ScanWords)
	for i := 0; i < T; i++ {
		if !scanner.Scan() {
			fmt.Fprintln(os.Stderr, "insufficient output")
			os.Exit(1)
		}
		var got int64
		fmt.Sscan(scanner.Text(), &got)
		if got != expected[i] {
			fmt.Fprintf(os.Stderr, "mismatch on test %d: expected %d got %d\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output after", T, "tests")
		os.Exit(1)
	}
	fmt.Println("all tests passed")
}
