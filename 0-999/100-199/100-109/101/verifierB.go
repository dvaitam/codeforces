package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type testCaseB struct {
	input string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := []testCaseB{
		{input: "2 1\n0 2\n"},
		{input: "3 3\n0 2\n2 3\n0 3\n"},
	}
	for i, t := range tests {
		expect := solveB(strings.NewReader(t.input))
		gotOut, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("test %d: execution failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(gotOut) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed: expected %q got %q\n", i+1, strings.TrimSpace(expect), strings.TrimSpace(gotOut))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

const mod = 1000000007

type seg struct{ r, l int }

func solveB(r io.Reader) string {
	in := bufio.NewReader(r)
	var n, m int
	fmt.Fscan(in, &n, &m)
	segs := make([]seg, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &segs[i].l, &segs[i].r)
	}
	sort.Slice(segs, func(i, j int) bool { return segs[i].r < segs[j].r })
	b := []int{0}
	c := make([]int, m)
	last := -1
	for i, s := range segs {
		if s.r != last {
			last = s.r
			b = append(b, s.r)
		}
		c[i] = len(b) - 1
	}
	E := len(b) - 1
	if b[E] != n {
		return "0\n"
	}
	f := make([]int, E+1)
	ssum := make([]int, E+1)
	f[0], ssum[0] = 1, 1
	for i, seg := range segs {
		ci := c[i]
		lo, hi := 0, ci
		for lo < hi {
			mid := (lo + hi) / 2
			if b[mid] < seg.l {
				lo = mid + 1
			} else {
				hi = mid
			}
		}
		if lo < ci {
			add := ssum[ci-1]
			if lo > 0 {
				add = (add - ssum[lo-1] + mod) % mod
			}
			f[ci] = (f[ci] + add) % mod
		}
		if i+1 == m || c[i+1] != ci {
			ssum[ci] = (ssum[ci-1] + f[ci]) % mod
		}
	}
	return fmt.Sprintf("%d\n", f[E])
}
