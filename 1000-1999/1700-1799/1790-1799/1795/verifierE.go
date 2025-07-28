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

type Fenwick struct {
	n   int
	bit []int64
}

func newFenwick(n int) *Fenwick {
	return &Fenwick{n: n, bit: make([]int64, n+2)}
}

func (f *Fenwick) add(i int, delta int64) {
	for i <= f.n {
		f.bit[i] += delta
		i += i & -i
	}
}

func (f *Fenwick) sum(i int) int64 {
	s := int64(0)
	for i > 0 {
		s += f.bit[i]
		i -= i & -i
	}
	return s
}

func (f *Fenwick) rangeSum(l, r int) int64 {
	if r < l {
		return 0
	}
	if l <= 1 {
		return f.sum(r)
	}
	return f.sum(r) - f.sum(l-1)
}

func compress(values []int64) ([]int64, map[int64]int) {
	uniq := make([]int64, len(values))
	copy(uniq, values)
	sort.Slice(uniq, func(i, j int) bool { return uniq[i] < uniq[j] })
	m := 1
	for i := 1; i < len(uniq); i++ {
		if uniq[i] != uniq[m-1] {
			uniq[m] = uniq[i]
			m++
		}
	}
	uniq = uniq[:m]
	mp := make(map[int64]int, len(uniq))
	for idx, v := range uniq {
		mp[v] = idx + 1
	}
	return uniq, mp
}

func expected(h []int64) string {
	n := len(h)
	if n == 1 {
		return fmt.Sprintf("%d", h[0])
	}
	A := make([]int64, n)
	B := make([]int64, n)
	for i := 0; i < n; i++ {
		idx := int64(i + 1)
		A[i] = h[i] + idx
		B[i] = h[i] - idx
	}
	valsA, mapA := compress(A)
	valsB, mapB := compress(B)
	bitCntRight := newFenwick(len(valsA))
	bitSumRight := newFenwick(len(valsA))
	right := make([]int64, n)
	for i := n - 1; i >= 0; i-- {
		if i < n-1 {
			v := A[i+1]
			pos := mapA[v]
			bitCntRight.add(pos, 1)
			bitSumRight.add(pos, v)
		}
		t := A[i]
		idx := sort.Search(len(valsA), func(j int) bool { return valsA[j] > t })
		if idx < len(valsA) {
			sum := bitSumRight.rangeSum(idx+1, len(valsA))
			cnt := bitCntRight.rangeSum(idx+1, len(valsA))
			right[i] = sum - int64(t)*cnt
		}
	}
	bitCntLeft := newFenwick(len(valsB))
	bitSumLeft := newFenwick(len(valsB))
	left := make([]int64, n)
	for i := 0; i < n; i++ {
		if i > 0 {
			v := B[i-1]
			pos := mapB[v]
			bitCntLeft.add(pos, 1)
			bitSumLeft.add(pos, v)
		}
		t := B[i]
		idx := sort.Search(len(valsB), func(j int) bool { return valsB[j] > t })
		if idx < len(valsB) {
			sum := bitSumLeft.rangeSum(idx+1, len(valsB))
			cnt := bitCntLeft.rangeSum(idx+1, len(valsB))
			left[i] = sum - int64(t)*cnt
		}
	}
	ans := int64(1 << 60)
	for i := 0; i < n; i++ {
		cost := h[i] + left[i] + right[i]
		if cost < ans {
			ans = cost
		}
	}
	return fmt.Sprintf("%d", ans)
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for tc := 0; tc < 100; tc++ {
		n := rng.Intn(8) + 1
		h := make([]int64, n)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			h[i] = int64(rng.Intn(20) + 1)
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", h[i]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expectedOut := expected(h)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", tc+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expectedOut {
			fmt.Printf("case %d failed: expected %s got %s\ninput:\n%s", tc+1, expectedOut, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
