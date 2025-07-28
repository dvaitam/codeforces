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

func solveF(n, m int, a []int, b []int, c [][]int) int64 {
	sumA, sumB := 0, 0
	for _, v := range a {
		sumA += v
	}
	for _, v := range b {
		sumB += v
	}
	if sumA > sumB {
		return -1
	}
	pow := make([]int, n)
	pow[0] = 1
	for i := 1; i < n; i++ {
		pow[i] = pow[i-1] * 5
	}
	finalNeed := 0
	for i := 0; i < n; i++ {
		finalNeed += a[i] * pow[i]
	}
	type state struct {
		need int
		v1   int
		v2   int
		rem  int
	}
	start := state{0, 0, 0, 0}
	dp := map[state]int{start: 0}
	best := int64(1<<63 - 1)
	for len(dp) > 0 {
		next := make(map[state]int)
		for st, cost := range dp {
			if int64(cost) >= best {
				continue
			}
			if st.v2 == m {
				if st.need == finalNeed && int64(cost) < best {
					best = int64(cost)
				}
				continue
			}
			i := st.v1
			j := st.v2
			cur := (st.need / pow[i]) % 5
			maxF := a[i] - cur
			if maxF > b[j]-st.rem {
				maxF = b[j] - st.rem
			}
			for f := 0; f <= maxF; f++ {
				nn := st.need + f*pow[i]
				nc := cost
				if f > 0 {
					nc += c[i][j]
				}
				nv1 := i + 1
				nv2 := j
				nr := st.rem + f
				if nv1 == n {
					nv1 = 0
					nv2 = j + 1
					nr = 0
				}
				nst := state{nn, nv1, nv2, nr}
				if nv2 == m {
					if nn == finalNeed && int64(nc) < best {
						best = int64(nc)
					}
					continue
				}
				if prev, ok := next[nst]; !ok || nc < prev {
					next[nst] = nc
				}
			}
		}
		dp = next
	}
	if best == int64(1<<63-1) {
		return -1
	}
	return best
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesF.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	in := bufio.NewReader(file)
	var T int
	fmt.Fscan(in, &T)
	for idx := 1; idx <= T; idx++ {
		var n, m int
		fmt.Fscan(in, &n, &m)
		a := make([]int, n)
		b := make([]int, m)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for j := 0; j < m; j++ {
			fmt.Fscan(in, &b[j])
		}
		c := make([][]int, n)
		for i := 0; i < n; i++ {
			c[i] = make([]int, m)
			for j := 0; j < m; j++ {
				fmt.Fscan(in, &c[i][j])
			}
		}
		expect := solveF(n, m, a, b, c)
		var input strings.Builder
		input.WriteString("1\n")
		fmt.Fprintf(&input, "%d %d\n", n, m)
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprintf(&input, "%d", a[i])
		}
		input.WriteByte('\n')
		for j := 0; j < m; j++ {
			if j > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprintf(&input, "%d", b[j])
		}
		input.WriteByte('\n')
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if j > 0 {
					input.WriteByte(' ')
				}
				fmt.Fprintf(&input, "%d", c[i][j])
			}
			input.WriteByte('\n')
		}
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		err := cmd.Run()
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n%s", idx, err, errBuf.String())
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(out.String())
		got, err2 := strconv.ParseInt(gotStr, 10, 64)
		if err2 != nil || got != expect {
			fmt.Printf("case %d failed: expected %d got %s\n", idx, expect, gotStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", T)
}
