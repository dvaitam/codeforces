package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func add(S []int64, l, r int, x int64) {
	if l > r {
		return
	}
	S[l] += x
	if r+1 < len(S) {
		S[r+1] -= x
	}
}

type orderedSet struct{ arr []int }

func (s *orderedSet) Insert(x int) int {
	i := sort.SearchInts(s.arr, x)
	if i == len(s.arr) || s.arr[i] != x {
		s.arr = append(s.arr, 0)
		copy(s.arr[i+1:], s.arr[i:])
		s.arr[i] = x
	}
	return i
}

func (s *orderedSet) At(i int) int { return s.arr[i] }
func (s *orderedSet) Size() int    { return len(s.arr) }

func solve(n int, a []int) []int64 {
	g := make([][]int, n+2)
	S := make([]int64, n+2)
	for i := 1; i <= n; i++ {
		g[a[i-1]] = append(g[a[i-1]], i)
	}
	se := orderedSet{}
	se.Insert(0)
	se.Insert(n + 1)
	for i := 1; i <= n; i++ {
		for _, x := range g[i] {
			idx := se.Insert(x)
			t1, t2 := idx-1, idx+1
			L := se.At(t1)
			R := se.At(t2)
			add(S, 0, L-1, int64(x-L)*int64(R-x)*int64(a[x-1]))
			add(S, R+1, n, int64(x-L)*int64(R-x)*int64(a[x-1]))
			add(S, L+1, x-1, int64(x-L-1)*int64(R-x)*int64(a[x-1]))
			add(S, x+1, R-1, int64(x-L)*int64(R-x-1)*int64(a[x-1]))
			if L != 0 && t1-1 >= 0 {
				LL := se.At(t1 - 1)
				add(S, L, L, int64(x-LL-1)*int64(R-x)*int64(a[x-1]))
			}
			if R != n+1 && t2+1 < se.Size() {
				RR := se.At(t2 + 1)
				add(S, R, R, int64(x-L)*int64(RR-x-1)*int64(a[x-1]))
			}
		}
	}
	ans := make([]int64, n)
	for i := 1; i <= n; i++ {
		S[i] += S[i-1]
		ans[i-1] = S[i]
	}
	return ans
}

func generate() (string, string) {
	const T = 100
	var in strings.Builder
	var out strings.Builder
	fmt.Fprintf(&in, "%d\n", T)
	rand.Seed(5)
	for t := 0; t < T; t++ {
		n := rand.Intn(10) + 1
		fmt.Fprintf(&in, "%d\n", n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = rand.Intn(n) + 1
			fmt.Fprintf(&in, "%d ", arr[i])
		}
		fmt.Fprintf(&in, "\n")
		res := solve(n, arr)
		for i, v := range res {
			if i+1 == len(res) {
				fmt.Fprintf(&out, "%d\n", v)
			} else {
				fmt.Fprintf(&out, "%d ", v)
			}
		}
	}
	return in.String(), out.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	in, exp := generate()
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("error running binary:", err)
		os.Exit(1)
	}
	got := buf.String()
	if strings.TrimSpace(got) != strings.TrimSpace(exp) {
		fmt.Println("wrong answer")
		fmt.Println("expected:\n" + exp)
		fmt.Println("got:\n" + got)
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
