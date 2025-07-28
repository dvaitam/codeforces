package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Node struct {
	cnt int
	val int
}

type MaxHeap struct{ data []Node }

func (h *MaxHeap) greater(i, j int) bool {
	a, b := h.data[i], h.data[j]
	if a.cnt == b.cnt {
		return a.val > b.val
	}
	return a.cnt > b.cnt
}

func (h *MaxHeap) Push(x Node) {
	h.data = append(h.data, x)
	i := len(h.data) - 1
	for i > 0 {
		p := (i - 1) / 2
		if h.greater(i, p) {
			h.data[i], h.data[p] = h.data[p], h.data[i]
			i = p
		} else {
			break
		}
	}
}

func (h *MaxHeap) Pop() Node {
	n := len(h.data)
	top := h.data[0]
	if n == 1 {
		h.data = h.data[:0]
		return top
	}
	h.data[0] = h.data[n-1]
	h.data = h.data[:n-1]
	i := 0
	for {
		l := 2*i + 1
		r := 2*i + 2
		largest := i
		if l < len(h.data) && h.greater(l, largest) {
			largest = l
		}
		if r < len(h.data) && h.greater(r, largest) {
			largest = r
		}
		if largest == i {
			break
		}
		h.data[i], h.data[largest] = h.data[largest], h.data[i]
		i = largest
	}
	return top
}

type testCaseD struct {
	n int
	m int
	a []int
	s []int
	d []int
}

func expectedD(tc testCaseD) (string, bool) {
	n, m := tc.n, tc.m
	cnt := make([]int, n+1)
	for _, x := range tc.a {
		if x >= 1 && x <= n {
			cnt[x]++
		}
	}
	h := &MaxHeap{}
	for x := 1; x <= n; x++ {
		if cnt[x] > 0 {
			h.Push(Node{cnt: cnt[x], val: x})
		}
	}
	ans := make([][]int, m)
	for i := 0; i < m; i++ {
		ans[i] = make([]int, tc.s[i])
		for j := 0; j < tc.s[i]; j++ {
			if j >= tc.d[i] {
				prev := ans[i][j-tc.d[i]]
				if cnt[prev] > 0 {
					h.Push(Node{cnt: cnt[prev], val: prev})
				}
			}
			if len(h.data) == 0 {
				return "-1", false
			}
			node := h.Pop()
			ans[i][j] = node.val
			cnt[node.val]--
		}
		for j := tc.s[i]; j < tc.s[i]+tc.d[i]; j++ {
			prev := ans[i][j-tc.d[i]]
			if cnt[prev] > 0 {
				h.Push(Node{cnt: cnt[prev], val: prev})
			}
		}
	}
	var sb strings.Builder
	for i := 0; i < m; i++ {
		for j := 0; j < len(ans[i]); j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(ans[i][j]))
		}
		if i+1 < m {
			sb.WriteByte('\n')
		}
	}
	return sb.String(), true
}

func genTestsD() []testCaseD {
	rand.Seed(4)
	tests := make([]testCaseD, 0, 100)
	for len(tests) < 100 {
		n := rand.Intn(6) + 1
		m := rand.Intn(3) + 1
		a := make([]int, n)
		for i := range a {
			a[i] = rand.Intn(n) + 1
		}
		s := make([]int, m)
		remaining := n
		for i := 0; i < m; i++ {
			if i == m-1 {
				s[i] = remaining
			} else {
				max := remaining - (m - i - 1)
				s[i] = rand.Intn(max) + 1
				remaining -= s[i]
			}
		}
		d := make([]int, m)
		for i := range d {
			d[i] = rand.Intn(s[i]) + 1
		}
		tests = append(tests, testCaseD{n: n, m: m, a: a, s: s, d: d})
	}
	return tests
}

func runCase(bin string, tc testCaseD) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	for i, v := range tc.s {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	for i, v := range tc.d {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')

	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	expect, ok := expectedD(tc)
	if !ok {
		if got != "-1" {
			return fmt.Errorf("expected -1 got %s", got)
		}
		return nil
	}
	if got != expect {
		return fmt.Errorf("expected %q got %q", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsD()
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
