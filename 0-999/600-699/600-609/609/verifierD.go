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

type item struct {
	t int
	c int64
}

func expectedDay(n, m, k int, s int64, a, b []int64, items []item) int {
	c1 := make([]int64, n+1)
	c2 := make([]int64, n+1)
	c1day := make([]int, n+1)
	c2day := make([]int, n+1)
	c1[1] = a[0]
	c1day[1] = 1
	for i := 2; i <= n; i++ {
		if c1[i-1] > a[i-1] {
			c1[i] = a[i-1]
			c1day[i] = i
		} else {
			c1[i] = c1[i-1]
			c1day[i] = c1day[i-1]
		}
	}
	c2[1] = b[0]
	c2day[1] = 1
	for i := 2; i <= n; i++ {
		if c2[i-1] > b[i-1] {
			c2[i] = b[i-1]
			c2day[i] = i
		} else {
			c2[i] = c2[i-1]
			c2day[i] = c2day[i-1]
		}
	}
	type node struct {
		x int
		c int64
	}
	q1 := make([]node, 0, m)
	q2 := make([]node, 0, m)
	for i, it := range items {
		if it.t == 1 {
			q1 = append(q1, node{i + 1, it.c})
		} else {
			q2 = append(q2, node{i + 1, it.c})
		}
	}
	tot1 := len(q1)
	tot2 := len(q2)
	sort.Slice(q1, func(i, j int) bool { return q1[i].c < q1[j].c })
	sort.Slice(q2, func(i, j int) bool { return q2[i].c < q2[j].c })
	sum1 := make([]int64, tot1+1)
	sum2 := make([]int64, tot2+1)
	for i := 1; i <= tot1; i++ {
		sum1[i] = sum1[i-1] + q1[i-1].c
	}
	for i := 1; i <= tot2; i++ {
		sum2[i] = sum2[i-1] + q2[i-1].c
	}
	jub := func(day int) bool {
		for i := 0; i <= k; i++ {
			if i <= tot1 && k-i <= tot2 {
				cost := sum1[i]*c1[day] + sum2[k-i]*c2[day]
				if cost <= s {
					return true
				}
			}
		}
		return false
	}
	if !jub(n) {
		return -1
	}
	l, r := 1, n
	for l < r {
		mid := (l + r) >> 1
		if jub(mid) {
			r = mid
		} else {
			l = mid + 1
		}
	}
	return l
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesD.txt")
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
		parts := strings.Fields(line)
		if len(parts) < 4 {
			fmt.Printf("test %d invalid\n", idx)
			os.Exit(1)
		}
		pos := 0
		n, _ := strconv.Atoi(parts[pos])
		pos++
		m, _ := strconv.Atoi(parts[pos])
		pos++
		k, _ := strconv.Atoi(parts[pos])
		pos++
		sVal, _ := strconv.ParseInt(parts[pos], 10, 64)
		pos++
		if len(parts) < pos+n+n+2*m {
			fmt.Printf("test %d invalid length\n", idx)
			os.Exit(1)
		}
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.ParseInt(parts[pos+i], 10, 64)
			a[i] = v
		}
		pos += n
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.ParseInt(parts[pos+i], 10, 64)
			b[i] = v
		}
		pos += n
		items := make([]item, m)
		for i := 0; i < m; i++ {
			t, _ := strconv.Atoi(parts[pos])
			pos++
			c, _ := strconv.ParseInt(parts[pos], 10, 64)
			pos++
			items[i] = item{t: t, c: c}
		}
		expDay := expectedDay(n, m, k, sVal, a, b, items)
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "%d %d %d %d", n, m, k, sVal)
		for i := 0; i < n; i++ {
			fmt.Fprintf(&buf, " %d", a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fprintf(&buf, " %d", b[i])
		}
		for i := 0; i < m; i++ {
			fmt.Fprintf(&buf, " %d %d", items[i].t, items[i].c)
		}
		buf.WriteByte('\n')
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(buf.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		fields := strings.Fields(string(out))
		if len(fields) == 0 {
			fmt.Printf("Test %d no output\n", idx)
			os.Exit(1)
		}
		gotDay, _ := strconv.Atoi(fields[0])
		if gotDay != expDay {
			fmt.Printf("Test %d failed: expected day %d got %d\n", idx, expDay, gotDay)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
