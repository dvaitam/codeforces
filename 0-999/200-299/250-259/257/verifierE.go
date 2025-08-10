package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type testCase struct {
	input  string
	output string
}

type BIT struct {
	n int
	t []int
}

func NewBIT(n int) *BIT {
	return &BIT{n, make([]int, n+1)}
}

func (b *BIT) Add(i, v int) {
	for x := i; x <= b.n; x += x & -x {
		b.t[x] += v
	}
}

func (b *BIT) Sum1(i int) int {
	s := 0
	for x := i; x > 0; x -= x & -x {
		s += b.t[x]
	}
	return s
}

func (b *BIT) Sum(l, r int) int {
	if r < l {
		return 0
	}
	return b.Sum1(r) - b.Sum1(l-1)
}

func (b *BIT) Select(k int) int {
	idx := 0
	bitMask := 1
	for bitMask<<1 <= b.n {
		bitMask <<= 1
	}
	for d := bitMask; d > 0; d >>= 1 {
		nxt := idx + d
		if nxt <= b.n && b.t[nxt] < k {
			idx = nxt
			k -= b.t[nxt]
		}
	}
	return idx + 1
}

type Person struct {
	t    int64
	s, f int
	idx  int
}

func solve(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n, m int
	fmt.Fscan(in, &n, &m)

	ti := make([]int64, n)
	si := make([]int, n)
	fi := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &ti[i], &si[i], &fi[i])
	}

	ord := make([]int, n)
	for i := range ord {
		ord[i] = i
	}
	sort.Slice(ord, func(i, j int) bool {
		if ti[ord[i]] == ti[ord[j]] {
			return ord[i] < ord[j]
		}
		return ti[ord[i]] < ti[ord[j]]
	})

	waiters := make([][]int, m+2)
	onDest := make([][]int, m+2)
	fen := NewBIT(m)

	active := 0
	x := 1
	var t int64
	p := 0
	ans := make([]int64, n)
	done := 0

	for done < n {
		if active == 0 && p < n && t < ti[ord[p]] {
			t = ti[ord[p]]
		}

		for p < n && ti[ord[p]] == t {
			id := ord[p]
			s := si[id]
			waiters[s] = append(waiters[s], id)
			fen.Add(s, 1)
			active++
			p++
		}

		if len(onDest[x]) > 0 {
			k := len(onDest[x])
			fen.Add(x, -k)
			for _, id := range onDest[x] {
				ans[id] = t
				done++
			}
			onDest[x] = onDest[x][:0]
			active -= k
		}

		if len(waiters[x]) > 0 {
			k := len(waiters[x])
			fen.Add(x, -k)
			for _, id := range waiters[x] {
				dest := fi[id]
				onDest[dest] = append(onDest[dest], id)
				fen.Add(dest, 1)
			}
			waiters[x] = waiters[x][:0]
		}

		if done == n {
			break
		}
		if active == 0 {
			if p >= n {
				break
			}
			continue
		}

		less := fen.Sum(1, x-1)
		greater := active - less
		dir := 1
		if greater < less {
			dir = -1
		}

		var dtArr int64
		if p < n {
			dtArr = ti[ord[p]] - t
		} else {
			dtArr = int64(1 << 62)
		}

		var y int
		if dir == 1 {
			k := fen.Sum1(x) + 1
			y = fen.Select(k)
		} else {
			k := fen.Sum1(x - 1)
			y = fen.Select(k)
		}
		dist := y - x
		if dist < 0 {
			dist = -dist
		}
		dist64 := int64(dist)

		if dtArr < dist64 {
			x += dir * int(dtArr)
			t += dtArr
		} else {
			x = y
			t += dist64
		}
	}

	var sb strings.Builder
	for i := 0; i < n; i++ {
		fmt.Fprintln(&sb, ans[i])
	}
	return sb.String()
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func generateTests() []testCase {
	rand.Seed(46)
	var tests []testCase
	fixed := []string{
		"1 2\n1 1 2\n",
		"2 3\n1 1 3\n2 3 1\n",
	}
	for _, f := range fixed {
		tests = append(tests, testCase{f, solve(f)})
	}
	for len(tests) < 100 {
		n := rand.Intn(20) + 1
		m := rand.Intn(10) + 2
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		for i := 0; i < n; i++ {
			t := rand.Intn(20) + 1
			s := rand.Intn(m) + 1
			f := rand.Intn(m) + 1
			for f == s {
				f = rand.Intn(m) + 1
			}
			fmt.Fprintf(&sb, "%d %d %d\n", t, s, f)
		}
		inp := sb.String()
		tests = append(tests, testCase{inp, solve(inp)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(t.output) {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected: %sGot: %s\n", i+1, t.input, strings.TrimSpace(t.output), strings.TrimSpace(got))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
