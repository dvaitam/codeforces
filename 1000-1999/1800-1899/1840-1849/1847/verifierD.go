package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type bit struct {
	n    int
	tree []int
}

func newBIT(n int) *bit { return &bit{n: n, tree: make([]int, n+2)} }
func (b *bit) add(i, val int) {
	for i <= b.n {
		b.tree[i] += val
		i += i & -i
	}
}
func (b *bit) sum(i int) int {
	res := 0
	for i > 0 {
		res += b.tree[i]
		i -= i & -i
	}
	return res
}

func solveD(n, m, q int, s string, intervals [][2]int, queries []int) []int {
	parent := make([]int, n+2)
	for i := 1; i <= n+1; i++ {
		parent[i] = i
	}
	var find func(int) int
	find = func(x int) int {
		if parent[x] == x {
			return x
		}
		parent[x] = find(parent[x])
		return parent[x]
	}
	order := make([]int, 0, n)
	for _, iv := range intervals {
		l, r := iv[0], iv[1]
		x := find(l)
		for x <= r {
			order = append(order, x)
			parent[x] = x + 1
			x = find(x)
		}
	}
	pos := make([]int, n+1)
	for i := range pos {
		pos[i] = -1
	}
	for idx, val := range order {
		pos[val] = idx + 1
	}
	b := newBIT(len(order))
	ones := 0
	arr := []byte(s)
	for i := 1; i <= n; i++ {
		if arr[i-1] == '1' {
			ones++
		}
		if p := pos[i]; p > 0 {
			b.add(p, int(arr[i-1]-'0'))
		}
	}
	result := make([]int, q)
	for qi, x := range queries {
		if arr[x-1] == '1' {
			arr[x-1] = '0'
			ones--
			if p := pos[x]; p > 0 {
				b.add(p, -1)
			}
		} else {
			arr[x-1] = '1'
			ones++
			if p := pos[x]; p > 0 {
				b.add(p, 1)
			}
		}
		k := ones
		if k > len(order) {
			k = len(order)
		}
		onesIn := b.sum(k)
		result[qi] = k - onesIn
	}
	return result
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(4)
	const t = 100
	var input bytes.Buffer
	fmt.Fprintln(&input, t)
	expectedLines := make([]string, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(20) + 1
		m := rand.Intn(n) + 1
		q := rand.Intn(n) + 1
		var sb strings.Builder
		sb.Grow(n)
		for j := 0; j < n; j++ {
			if rand.Intn(2) == 0 {
				sb.WriteByte('0')
			} else {
				sb.WriteByte('1')
			}
		}
		s := sb.String()
		intervals := make([][2]int, m)
		for j := 0; j < m; j++ {
			l := rand.Intn(n) + 1
			r := rand.Intn(n-l+1) + l
			intervals[j] = [2]int{l, r}
		}
		queries := make([]int, q)
		for j := 0; j < q; j++ {
			queries[j] = rand.Intn(n) + 1
		}
		fmt.Fprintf(&input, "%d %d %d\n", n, m, q)
		fmt.Fprintln(&input, s)
		for _, iv := range intervals {
			fmt.Fprintf(&input, "%d %d\n", iv[0], iv[1])
		}
		for _, x := range queries {
			fmt.Fprintln(&input, x)
		}
		res := solveD(n, m, q, s, intervals, queries)
		expectedLines[i] = strings.TrimSpace(strings.Join(intSliceToStr(res), " "))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error executing binary: %v\n", err)
		os.Exit(1)
	}
	outputs := splitNonEmptyLines(string(out))
	if len(outputs) != t {
		fmt.Fprintf(os.Stderr, "expected %d lines of output, got %d\n", t, len(outputs))
		fmt.Fprint(os.Stderr, string(out))
		os.Exit(1)
	}
	for i := 0; i < t; i++ {
		if outputs[i] != expectedLines[i] {
			fmt.Fprintf(os.Stderr, "mismatch on case %d: expected %q got %q\n", i+1, expectedLines[i], outputs[i])
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}

func intSliceToStr(a []int) []string {
	res := make([]string, len(a))
	for i, v := range a {
		res[i] = strconv.Itoa(v)
	}
	return res
}

func splitNonEmptyLines(s string) []string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	res := lines[:0]
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			res = append(res, line)
		}
	}
	return res
}
