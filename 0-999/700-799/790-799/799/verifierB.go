package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type tshirt struct{ price, idx int }

type testB struct {
	n       int
	prices  []int
	a       []int
	b       []int
	m       int
	queries []int
}

func solveB(tc testB) []int {
	lists := make([][]tshirt, 4)
	for i := 0; i < tc.n; i++ {
		if tc.a[i] == tc.b[i] {
			lists[tc.a[i]] = append(lists[tc.a[i]], tshirt{tc.prices[i], i})
		} else {
			lists[tc.a[i]] = append(lists[tc.a[i]], tshirt{tc.prices[i], i})
			lists[tc.b[i]] = append(lists[tc.b[i]], tshirt{tc.prices[i], i})
		}
	}
	for c := 1; c <= 3; c++ {
		sort.Slice(lists[c], func(i, j int) bool { return lists[c][i].price < lists[c][j].price })
	}
	ptr := make([]int, 4)
	sold := make([]bool, tc.n)
	res := make([]int, tc.m)
	for i := 0; i < tc.m; i++ {
		c := tc.queries[i]
		for ptr[c] < len(lists[c]) && sold[lists[c][ptr[c]].idx] {
			ptr[c]++
		}
		if ptr[c] == len(lists[c]) {
			res[i] = -1
		} else {
			res[i] = lists[c][ptr[c]].price
			sold[lists[c][ptr[c]].idx] = true
			ptr[c]++
		}
	}
	return res
}

func genB() (string, []int) {
	n := rand.Intn(8) + 1
	prices := rand.Perm(100)[:n]
	for i := range prices {
		prices[i]++
	}
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rand.Intn(3) + 1
		b[i] = rand.Intn(3) + 1
	}
	m := rand.Intn(n) + 1
	queries := make([]int, m)
	for i := 0; i < m; i++ {
		queries[i] = rand.Intn(3) + 1
	}
	tc := testB{n, prices, a, b, m, queries}
	outNums := solveB(tc)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range prices {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", m))
	for i, v := range queries {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	input := sb.String()
	return input, outNums
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return buf.String(), err
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go <binary>")
		return
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, expSlice := genB()
		got, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\nInput:\n%sOutput:\n%s", i+1, err, in, got)
			return
		}
		gotFields := strings.Fields(got)
		if len(gotFields) != len(expSlice) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected %d numbers, got %d\n%s", i+1, in, len(expSlice), len(gotFields), got)
			return
		}
		for j, g := range gotFields {
			var val int
			fmt.Sscan(g, &val)
			if val != expSlice[j] {
				fmt.Printf("Test %d failed at position %d\nInput:\n%sExpected:\n%v\nGot:\n%s", i+1, j+1, in, expSlice, got)
				return
			}
		}
	}
	fmt.Println("All tests passed")
}
