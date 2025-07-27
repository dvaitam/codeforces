package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	n         int
	row1      []int
	row2      []int
	input     string
	expectedK int
}

func solveCase(n int, a1, a2 []int) (int, []int) {
	a := make([]int, n+1)
	b := make([]int, n+1)
	copy(a[1:], a1)
	copy(b[1:], a2)
	c := make([]int, n+1)
	d := make([]int, n+1)
	num := make([]int, n+1)
	res := make([]int, n+1)
	vis := make([]bool, n+1)

	for i := 1; i <= n; i++ {
		if c[a[i]] == 0 {
			c[a[i]] = i
		} else {
			d[a[i]] = i
		}
	}
	for i := 1; i <= n; i++ {
		if c[b[i]] == 0 {
			c[b[i]] = i
		} else {
			d[b[i]] = i
		}
	}
	for i := 1; i <= n; i++ {
		num[a[i]]++
		num[b[i]]++
	}
	for i := 1; i <= n; i++ {
		if num[i] != 2 {
			return -1, nil
		}
	}
	for i := 1; i <= n; i++ {
		if vis[i] {
			continue
		}
		cur := i
		x := c[i]
		buf := []int{x}
		cnt := 0
		if a[x] != i {
			res[x] = 1
			cnt++
		}
		for {
			vis[cur] = true
			nxt := a[x] + b[x] - cur
			x = c[nxt] + d[nxt] - x
			cur = nxt
			if vis[cur] {
				break
			}
			buf = append(buf, x)
			if b[x] == nxt {
				res[x] = 1
				cnt++
			}
		}
		if cnt > len(buf)-cnt {
			for _, v := range buf {
				if res[v] == 1 {
					res[v] = 0
				} else {
					res[v] = 1
				}
			}
		}
	}
	var outIdx []int
	for i := 1; i <= n; i++ {
		if res[i] == 1 {
			outIdx = append(outIdx, i)
		}
	}
	return len(outIdx), outIdx
}

func performSwaps(r1, r2 []int, idx []int) ([]int, []int) {
	row1 := append([]int(nil), r1...)
	row2 := append([]int(nil), r2...)
	for _, p := range idx {
		row1[p], row2[p] = row2[p], row1[p]
	}
	return row1, row2
}

func isPermutation(row []int) bool {
	seen := make([]bool, len(row))
	for _, v := range row[1:] {
		if v <= 0 || v >= len(row) || seen[v] {
			return false
		}
		seen[v] = true
	}
	return true
}

func buildCase(n int, r1, r2 []int) testCase {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range r1 {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteString("\n")
	for i, v := range r2 {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteString("\n")
	k, _ := solveCase(n, append([]int(nil), r1...), append([]int(nil), r2...))
	return testCase{n: n, row1: r1, row2: r2, input: sb.String(), expectedK: k}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(10) + 1
	positions := make([]int, 2*n)
	for i := range positions {
		positions[i] = i
	}
	rng.Shuffle(len(positions), func(i, j int) { positions[i], positions[j] = positions[j], positions[i] })
	r1 := make([]int, n)
	r2 := make([]int, n)
	for val := 1; val <= n; val++ {
		p1 := positions[2*(val-1)]
		p2 := positions[2*(val-1)+1]
		if p1 < n {
			r1[p1] = val
		} else {
			r2[p1-n] = val
		}
		if p2 < n {
			r1[p2] = val
		} else {
			r2[p2-n] = val
		}
	}
	return buildCase(n, r1, r2)
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	if len(fields) == 0 {
		return fmt.Errorf("no output")
	}
	var k int
	if _, err := fmt.Sscan(fields[0], &k); err != nil {
		return fmt.Errorf("bad first line")
	}
	if k == -1 {
		if tc.expectedK != -1 {
			return fmt.Errorf("expected %d got -1", tc.expectedK)
		}
		return nil
	}
	if tc.expectedK == -1 {
		return fmt.Errorf("expected -1 got %d", k)
	}
	if len(fields)-1 != k {
		return fmt.Errorf("expected %d indices, got %d", k, len(fields)-1)
	}
	idx := make([]int, k)
	used := make(map[int]bool)
	for i := 0; i < k; i++ {
		var v int
		if _, err := fmt.Sscan(fields[i+1], &v); err != nil {
			return fmt.Errorf("bad index")
		}
		if v < 1 || v > tc.n || used[v] {
			return fmt.Errorf("invalid index %d", v)
		}
		used[v] = true
		idx[i] = v
	}
	r1, r2 := performSwaps(append([]int{0}, tc.row1...), append([]int{0}, tc.row2...), idx)
	if !isPermutation(r1) || !isPermutation(r2) {
		return fmt.Errorf("result not permutations")
	}
	if k != tc.expectedK {
		return fmt.Errorf("expected %d swaps got %d", tc.expectedK, k)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{buildCase(1, []int{1}, []int{1})}
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
