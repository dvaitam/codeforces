package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

type Case struct {
	a, b, c []int
	ds      []int
	ans     []int
}

func solveOne(a, b, c []int, ds []int) []int {
	n := len(a)
	m := 1 << n
	type pair struct {
		cost int
		val  int
	}
	res := make([]pair, 0, m)
	for mask := 0; mask < m; mask++ {
		cost := 0
		A := make([]int, n)
		B := make([]int, n)
		for i := 0; i < n; i++ {
			if mask>>i&1 == 1 {
				cost += c[i]
				A[i] = b[i]
				B[i] = a[i]
			} else {
				A[i] = a[i]
				B[i] = b[i]
			}
		}
		gA := A[0]
		for i := 1; i < n; i++ {
			gA = gcd(gA, A[i])
		}
		gB := B[0]
		for i := 1; i < n; i++ {
			gB = gcd(gB, B[i])
		}
		val := gA + gB
		res = append(res, pair{cost, val})
	}
	sort.Slice(res, func(i, j int) bool { return res[i].cost < res[j].cost })
	best := make([]int, len(res))
	cur := 0
	for i, p := range res {
		if p.val > cur {
			cur = p.val
		}
		best[i] = cur
	}
	ans := make([]int, len(ds))
	for i, dv := range ds {
		j := sort.Search(len(res), func(k int) bool { return res[k].cost > dv }) - 1
		if j >= 0 {
			ans[i] = best[j]
		} else {
			ans[i] = 0
		}
	}
	return ans
}

func genCases(n int) []Case {
	rand.Seed(time.Now().UnixNano())
	cs := make([]Case, n)
	for i := 0; i < n; i++ {
		size := rand.Intn(4) + 2
		a := make([]int, size)
		b := make([]int, size)
		cst := make([]int, size)
		for j := 0; j < size; j++ {
			a[j] = rand.Intn(20) + 1
			b[j] = rand.Intn(20) + 1
			cst[j] = rand.Intn(10) + 1
		}
		q := rand.Intn(3) + 1
		ds := make([]int, q)
		for j := 0; j < q; j++ {
			ds[j] = rand.Intn(40)
		}
		ans := solveOne(a, b, cst, ds)
		cs[i] = Case{a, b, cst, ds, ans}
	}
	return cs
}

func buildInput(cs []Case) string {
	var sb strings.Builder
	fmt.Fprintln(&sb, len(cs))
	for _, c := range cs {
		fmt.Fprintf(&sb, "%d %d\n", len(c.a), len(c.ds))
		for j, v := range c.a {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprint(&sb, v)
		}
		sb.WriteByte('\n')
		for j, v := range c.b {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprint(&sb, v)
		}
		sb.WriteByte('\n')
		for j, v := range c.c {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprint(&sb, v)
		}
		sb.WriteByte('\n')
		for j, v := range c.ds {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprint(&sb, v)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cs := genCases(100)
	input := buildInput(cs)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("runtime error:", err)
		os.Exit(1)
	}
	outputs := strings.Fields(strings.TrimSpace(out.String()))
	totalLines := 0
	for _, c := range cs {
		totalLines += len(c.ds)
	}
	if len(outputs) != totalLines {
		fmt.Printf("expected %d numbers got %d\n", totalLines, len(outputs))
		os.Exit(1)
	}
	idx := 0
	for ci, c := range cs {
		for j := 0; j < len(c.ds); j++ {
			v, err := strconv.Atoi(outputs[idx])
			if err != nil || v != c.ans[j] {
				fmt.Printf("case %d query %d: expected %d got %s\n", ci+1, j+1, c.ans[j], outputs[idx])
				os.Exit(1)
			}
			idx++
		}
	}
	fmt.Println("All tests passed")
}
