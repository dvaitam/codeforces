package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func modPow(a, e, mod int64) int64 {
	res := int64(1)
	a %= mod
	for e > 0 {
		if e&1 != 0 {
			res = (res * a) % mod
		}
		a = (a * a) % mod
		e >>= 1
	}
	return res
}

func solve(n, m, k int, cells [][3]int, p int64) int64 {
	known := make([][]bool, n)
	for i := 0; i < n; i++ {
		known[i] = make([]bool, m)
	}
	bRow := make([]int, n)
	bCol := make([]int, m)
	for i := 0; i < n; i++ {
		bRow[i] = 1
	}
	for j := 0; j < m; j++ {
		bCol[j] = 1
	}
	knownR := make([]int, n)
	knownC := make([]int, m)
	for _, c := range cells {
		a := c[0] - 1
		b := c[1] - 1
		y := 0
		if c[2] < 0 {
			y = 1
		}
		known[a][b] = true
		knownR[a]++
		knownC[b]++
		bRow[a] ^= y
		bCol[b] ^= y
	}
	for i := 0; i < n; i++ {
		if knownR[i] == m && bRow[i] != 0 {
			return 0
		}
	}
	for j := 0; j < m; j++ {
		if knownC[j] == n && bCol[j] != 0 {
			return 0
		}
	}
	sumR, sumC := 0, 0
	for i := 0; i < n; i++ {
		sumR ^= bRow[i] & 1
	}
	for j := 0; j < m; j++ {
		sumC ^= bCol[j] & 1
	}
	if sumR != sumC {
		return 0
	}
	visitedR := make([]bool, n)
	visitedC := make([]bool, m)
	comps := 0
	queue := []int{}
	for i := 0; i < n; i++ {
		if !visitedR[i] && knownR[i] < m {
			comps++
			visitedR[i] = true
			queue = queue[:0]
			queue = append(queue, i)
			for qi := 0; qi < len(queue); qi++ {
				u := queue[qi]
				if u < n {
					r := u
					for cj := 0; cj < m; cj++ {
						if known[r][cj] || visitedC[cj] {
							continue
						}
						visitedC[cj] = true
						queue = append(queue, n+cj)
					}
				} else {
					cj := u - n
					for rr := 0; rr < n; rr++ {
						if known[rr][cj] || visitedR[rr] {
							continue
						}
						visitedR[rr] = true
						queue = append(queue, rr)
					}
				}
			}
		}
	}
	for j := 0; j < m; j++ {
		if knownC[j] == n {
			comps++
		}
	}
	U := int64(n)*int64(m) - int64(k)
	F := U - int64(n+m) + int64(comps)
	return modPow(2, F, p)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesE.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 4 {
			fmt.Printf("invalid testcase %d\n", idx)
			os.Exit(1)
		}
		n := atoi(fields[0])
		m := atoi(fields[1])
		k := atoi(fields[2])
		pos := 3
		cells := make([][3]int, k)
		for i := 0; i < k; i++ {
			a := atoi(fields[pos])
			b := atoi(fields[pos+1])
			c := atoi(fields[pos+2])
			cells[i] = [3]int{a, b, c}
			pos += 3
		}
		p := int64(atoi(fields[pos]))
		exp := solve(n, m, k, cells, p)
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "%d %d\n", n, m)
		fmt.Fprintf(&buf, "%d\n", k)
		for _, c := range cells {
			fmt.Fprintf(&buf, "%d %d %d\n", c[0], c[1], c[2])
		}
		fmt.Fprintf(&buf, "%d\n", p)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(buf.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		var got int64
		fmt.Sscan(strings.TrimSpace(string(out)), &got)
		if got != exp {
			fmt.Printf("Test %d failed: expected %d got %s\n", idx, exp, strings.TrimSpace(string(out)))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}

func atoi(s string) int {
	var v int
	fmt.Sscan(s, &v)
	return v
}
