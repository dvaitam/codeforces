package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Query struct {
	l, r int
	x    int64
}

type Test struct {
	n, m, q int
	a       []int64
	b       []int64
	queries []Query
}

func generateTests() []Test {
	rand.Seed(46)
	tests := make([]Test, 0, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(3) + 2 // 2..4
		m := n + rand.Intn(3) // >=n
		q := rand.Intn(3) + 1 // 1..3
		a := make([]int64, n)
		b := make([]int64, m)
		for j := range a {
			a[j] = int64(rand.Intn(5))
		}
		for j := range b {
			b[j] = int64(rand.Intn(5))
		}
		queries := make([]Query, q)
		for j := 0; j < q; j++ {
			l := rand.Intn(n) + 1
			r := rand.Intn(n-l+1) + l
			x := int64(rand.Intn(5) - 2)
			queries[j] = Query{l, r, x}
		}
		tests = append(tests, Test{n: n, m: m, q: q, a: a, b: b, queries: queries})
	}
	// simple edge case
	tests = append(tests, Test{n: 2, m: 3, q: 1, a: []int64{1, 2}, b: []int64{2, 3, 4}, queries: []Query{{1, 2, 1}}})
	return tests
}

func minAbs(arr []int64, x int64) int64 {
	i := sort.Search(len(arr), func(i int) bool { return arr[i] >= x })
	best := int64(1<<63 - 1)
	if i < len(arr) {
		diff := arr[i] - x
		if diff < 0 {
			diff = -diff
		}
		if diff < best {
			best = diff
		}
	}
	if i > 0 {
		diff := x - arr[i-1]
		if diff < 0 {
			diff = -diff
		}
		if diff < best {
			best = diff
		}
	}
	return best
}

func solve(t Test) []int64 {
	A := int64(0)
	for i, v := range t.a {
		if (i+1)%2 == 1 {
			A += v
		} else {
			A -= v
		}
	}
	bPrefix := make([]int64, t.m+1)
	for i := 1; i <= t.m; i++ {
		v := t.b[i-1]
		if i%2 == 1 {
			bPrefix[i] = bPrefix[i-1] + v
		} else {
			bPrefix[i] = bPrefix[i-1] - v
		}
	}
	arr := make([]int64, t.m-t.n+1)
	for j := 0; j <= t.m-t.n; j++ {
		val := bPrefix[j+t.n] - bPrefix[j]
		if j%2 == 1 {
			val = -val
		}
		arr[j] = val
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	res := make([]int64, 0, t.q+1)
	res = append(res, minAbs(arr, A))
	for _, qu := range t.queries {
		if (qu.r-qu.l+1)%2 == 1 {
			if qu.l%2 == 1 {
				A += qu.x
			} else {
				A -= qu.x
			}
		}
		res = append(res, minAbs(arr, A))
	}
	return res
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: verifierE <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	passed := 0
	for i, t := range tests {
		input := fmt.Sprintf("%d %d %d\n", t.n, t.m, t.q)
		for _, v := range t.a {
			input += fmt.Sprintf("%d ", v)
		}
		input += "\n"
		for _, v := range t.b {
			input += fmt.Sprintf("%d ", v)
		}
		input += "\n"
		for _, qu := range t.queries {
			input += fmt.Sprintf("%d %d %d\n", qu.l, qu.r, qu.x)
		}
		wants := solve(t)
		output, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("Test %d exec err %v\n", i+1, err)
			continue
		}
		lines := strings.Fields(strings.TrimSpace(output))
		if len(lines) != len(wants) {
			fmt.Printf("Test %d wrong number lines got %v\n", i+1, output)
			continue
		}
		ok := true
		for k, w := range wants {
			g, err := strconv.ParseInt(lines[k], 10, 64)
			if err != nil || g != w {
				ok = false
				break
			}
		}
		if !ok {
			fmt.Printf("Test %d expected %v got %s\n", i+1, wants, output)
			continue
		}
		passed++
	}
	fmt.Printf("Passed %d/%d tests\n", passed, len(tests))
}
