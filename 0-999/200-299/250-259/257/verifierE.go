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
	persons := make([]Person, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &persons[i].t, &persons[i].s, &persons[i].f)
		persons[i].idx = i
	}
	events := make([]Person, n)
	copy(events, persons)
	sort.Slice(events, func(i, j int) bool { return events[i].t < events[j].t })

	waitingBIT := NewBIT(m)
	insideBIT := NewBIT(m)
	waitingQ := make([][]int, m+1)
	insideQ := make([][]int, m+1)
	dest := make([]int, n)
	for i := range persons {
		dest[persons[i].idx] = persons[i].f
	}
	ans := make([]int64, n)

	t := int64(0)
	x := 1
	ptr := 0
	delivered := 0

	dir := func() int {
		up := waitingBIT.Sum(x+1, m) + insideBIT.Sum(x+1, m)
		down := waitingBIT.Sum(1, x-1) + insideBIT.Sum(1, x-1)
		if up+down == 0 {
			return 0
		}
		if up >= down {
			return 1
		}
		return -1
	}

	for delivered < n {
		d := dir()
		if d == 0 {
			if ptr < n {
				t = max64(t, events[ptr].t)
				for ptr < n && events[ptr].t == t {
					e := &events[ptr]
					waitingQ[e.s] = append(waitingQ[e.s], e.idx)
					waitingBIT.Add(e.s, 1)
					ptr++
				}
				continue
			} else {
				break
			}
		}
		tNextArr := int64(1 << 62)
		if ptr < n {
			tNextArr = events[ptr].t
		}
		var fExit, dExit int64
		if d == 1 && insideBIT.Sum(x+1, m) > 0 {
			s0 := insideBIT.Sum1(x)
			fExit = int64(insideBIT.Select(s0 + 1))
			dExit = fExit - int64(x)
		} else if d == -1 && insideBIT.Sum(1, x-1) > 0 {
			s1 := insideBIT.Sum1(x - 1)
			fExit = int64(insideBIT.Select(s1))
			dExit = int64(x) - fExit
		} else {
			dExit = int64(1 << 62)
		}
		var fEntry, dEntry int64
		if d == 1 && waitingBIT.Sum(x+1, m) > 0 {
			s0 := waitingBIT.Sum1(x)
			fEntry = int64(waitingBIT.Select(s0 + 1))
			dEntry = fEntry - int64(x)
		} else if d == -1 && waitingBIT.Sum(1, x-1) > 0 {
			s1 := waitingBIT.Sum1(x - 1)
			fEntry = int64(waitingBIT.Select(s1))
			dEntry = int64(x) - fEntry
		} else {
			dEntry = int64(1 << 62)
		}
		dEvent := dExit
		fEvent := fExit
		if dEntry < dEvent {
			dEvent = dEntry
			fEvent = fEntry
		}
		tEvent := t + dEvent
		if tNextArr <= tEvent {
			dt := tNextArr - t
			x += d * int(dt)
			t = tNextArr
			for ptr < n && events[ptr].t == t {
				e := &events[ptr]
				waitingQ[e.s] = append(waitingQ[e.s], e.idx)
				waitingBIT.Add(e.s, 1)
				ptr++
			}
			if len(insideQ[x]) > 0 {
				for _, id := range insideQ[x] {
					ans[id] = t
					delivered++
				}
				insideBIT.Add(x, -len(insideQ[x]))
				insideQ[x] = nil
			}
			if len(waitingQ[x]) > 0 {
				for _, id := range waitingQ[x] {
					f := dest[id]
					insideQ[f] = append(insideQ[f], id)
					insideBIT.Add(f, 1)
				}
				waitingBIT.Add(x, -len(waitingQ[x]))
				waitingQ[x] = nil
			}
		} else {
			t = tEvent
			x = int(fEvent)
			if len(insideQ[x]) > 0 {
				for _, id := range insideQ[x] {
					ans[id] = t
					delivered++
				}
				insideBIT.Add(x, -len(insideQ[x]))
				insideQ[x] = nil
			}
			if len(waitingQ[x]) > 0 {
				for _, id := range waitingQ[x] {
					f := dest[id]
					insideQ[f] = append(insideQ[f], id)
					insideBIT.Add(f, 1)
				}
				waitingBIT.Add(x, -len(waitingQ[x]))
				waitingQ[x] = nil
			}
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
