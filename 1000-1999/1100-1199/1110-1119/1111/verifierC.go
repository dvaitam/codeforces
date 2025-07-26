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

type test struct {
	n   int
	k   int
	A   int64
	B   int64
	pos []int64
}

func work(l, beg, end int, st int64, A, B int64, p []int64) int64 {
	if beg > end {
		return A
	}
	if l == 0 {
		return int64(end-beg+1) * B
	}
	half := int64(1) << (l - 1)
	m := st + half
	cnt := end - beg + 1
	idx := sort.Search(cnt, func(i int) bool { return p[beg+i] >= m })
	mid := beg + idx
	total := int64(cnt)
	costAll := total * (int64(1) << l) * B
	costSplit := work(l-1, beg, mid-1, st, A, B, p) + work(l-1, mid, end, m, A, B, p)
	if costAll < costSplit {
		return costAll
	}
	return costSplit
}

func solve(tc test) string {
	p := append([]int64(nil), tc.pos...)
	sort.Slice(p, func(i, j int) bool { return p[i] < p[j] })
	res := work(tc.n, 0, tc.k-1, 1, tc.A, tc.B, p)
	return fmt.Sprintf("%d", res)
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(44))
	tests := make([]test, 0, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(3) + 1
		maxPos := 1 << uint(n)
		k := rng.Intn(maxPos)
		A := int64(rng.Intn(5) + 1)
		B := int64(rng.Intn(5) + 1)
		pos := make([]int64, k)
		for j := 0; j < k; j++ {
			pos[j] = int64(rng.Intn(maxPos) + 1)
		}
		tests = append(tests, test{n, k, A, B, pos})
	}
	return tests
}

func run(binary string, input string) (string, error) {
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d %d\n", tc.n, tc.k, tc.A, tc.B)
		for j, v := range tc.pos {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		if tc.k > 0 {
			sb.WriteByte('\n')
		} else {
			sb.WriteString("\n")
		}
		input := sb.String()
		want := solve(tc)
		got, err := run(binary, input)
		if err != nil {
			fmt.Printf("Test %d: error running binary: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected: %s\nGot: %s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
