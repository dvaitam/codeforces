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

type testCaseH struct {
	n       int
	parent  []int
	hunger  []int64
	q       int
	queries [][2]interface{}
}

func generateTestsH(num int) []testCaseH {
	rand.Seed(time.Now().UnixNano())
	tests := make([]testCaseH, num)
	for i := 0; i < num; i++ {
		n := rand.Intn(10) + 1
		parent := make([]int, n+1)
		for v := 2; v <= n; v++ {
			parent[v] = rand.Intn(v-1) + 1
		}
		hunger := make([]int64, n+1)
		for v := 1; v <= n; v++ {
			hunger[v] = rand.Int63n(10)
		}
		qn := rand.Intn(10) + 1
		queries := make([][2]interface{}, qn)
		for j := 0; j < qn; j++ {
			v := rand.Intn(n) + 1
			x := rand.Int63n(10)
			queries[j] = [2]interface{}{v, x}
		}
		tests[i] = testCaseH{n: n, parent: parent, hunger: hunger, q: qn, queries: queries}
	}
	return tests
}

func solveH(tc testCaseH) []int64 {
	n := tc.n
	parent := tc.parent
	children := make([]int, n+1)
	for i := 2; i <= n; i++ {
		children[parent[i]]++
	}
	leaves := []int{}
	for i := 1; i <= n; i++ {
		if children[i] == 0 {
			leaves = append(leaves, i)
		}
	}
	leafIndex := make(map[int]int)
	for idx, v := range leaves {
		leafIndex[v] = idx
	}
	m := len(leaves)

	leafH := make([]int64, m)
	for idx, v := range leaves {
		leafH[idx] = tc.hunger[v]
	}

	calc := func() int64 {
		vals := make([]int64, m)
		copy(vals, leafH)
		sort.Slice(vals, func(i, j int) bool { return vals[i] > vals[j] })
		var ans int64
		for i, v := range vals {
			if v == 0 {
				continue
			}
			t := int64(i+1) + (v-1)*int64(m)
			if t > ans {
				ans = t
			}
		}
		return ans % 998244353
	}

	res := make([]int64, tc.q+1)
	res[0] = calc()
	for qi, q := range tc.queries {
		v := q[0].(int)
		x := q[1].(int64)
		if idx, ok := leafIndex[v]; ok {
			leafH[idx] = x
		}
		res[qi+1] = calc()
	}
	return res
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTestsH(100)
	var input bytes.Buffer
	for _, tc := range tests {
		fmt.Fprintln(&input, tc.n)
		for i := 2; i <= tc.n; i++ {
			if i > 2 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, tc.parent[i])
		}
		if tc.n >= 2 {
			input.WriteByte('\n')
		}
		for i := 1; i <= tc.n; i++ {
			if i > 1 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, tc.hunger[i])
		}
		input.WriteByte('\n')
		fmt.Fprintln(&input, tc.q)
		for _, q := range tc.queries {
			fmt.Fprintf(&input, "%d %d\n", q[0].(int), q[1].(int64))
		}
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("binary execution failed:", err)
		os.Exit(1)
	}
	outputs := strings.Split(strings.TrimSpace(string(out)), "\n")
	idx := 0
	for tnum, tc := range tests {
		expected := solveH(tc)
		for _, val := range expected {
			if idx >= len(outputs) || strings.TrimSpace(outputs[idx]) != fmt.Sprint(val) {
				fmt.Printf("test %d mismatch at line %d\n", tnum+1, idx)
				os.Exit(1)
			}
			idx++
		}
	}
	if idx != len(outputs) {
		fmt.Printf("expected %d output lines, got %d\n", idx, len(outputs))
		os.Exit(1)
	}
	fmt.Println("all tests passed")
}
