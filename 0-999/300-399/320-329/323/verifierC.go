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

func generateCase(seed int64) (string, []int) {
	rand.Seed(seed)
	n := 5 + rand.Intn(5) // 5..9
	p := rand.Perm(n)
	q := rand.Perm(n)
	for i := 0; i < n; i++ {
		p[i]++
		q[i]++
	}
	m := 3 + rand.Intn(3) // 3..5
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", p[i])
	}
	b.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", q[i])
	}
	b.WriteByte('\n')
	fmt.Fprintf(&b, "%d\n", m)
	queries := make([][4]int, m)
	for i := 0; i < m; i++ {
		for j := 0; j < 4; j++ {
			queries[i][j] = 1 + rand.Intn(n)
		}
		fmt.Fprintf(&b, "%d %d %d %d\n", queries[i][0], queries[i][1], queries[i][2], queries[i][3])
	}
	ans := solveCase(n, p, q, queries)
	return b.String(), ans
}

func solveCase(n int, p, q []int, queries [][4]int) []int {
	posP := make([]int, n+1)
	posQ := make([]int, n+1)
	for i := 0; i < n; i++ {
		posP[p[i]] = i + 1
		posQ[q[i]] = i + 1
	}
	x := 0
	res := make([]int, len(queries))
	for qi, qu := range queries {
		a, b, c, d := qu[0], qu[1], qu[2], qu[3]
		f := func(z int) int { return ((z - 1 + x) % n) + 1 }
		l1 := f(a)
		r1 := f(b)
		if l1 > r1 {
			l1, r1 = r1, l1
		}
		l2 := f(c)
		r2 := f(d)
		if l2 > r2 {
			l2, r2 = r2, l2
		}
		count := 0
		for v := 1; v <= n; v++ {
			if posP[v] >= l1 && posP[v] <= r1 && posQ[v] >= l2 && posQ[v] <= r2 {
				count++
			}
		}
		res[qi] = count
		x = count + 1
	}
	return res
}

func runCase(bin string, seed int64) error {
	input, ans := generateCase(seed)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	outLines := strings.Fields(buf.String())
	if len(outLines) != len(ans) {
		return fmt.Errorf("wrong answer")
	}
	for i, v := range outLines {
		x, err := strconv.Atoi(v)
		if err != nil || x != ans[i] {
			return fmt.Errorf("wrong answer")
		}
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		if err := runCase(bin, int64(i)+time.Now().UnixNano()); err != nil {
			fmt.Printf("test %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
