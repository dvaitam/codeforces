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

func expected(n, m int, A, B []int) string {
	sort.Ints(A)
	sort.Ints(B)
	type block struct{ a, c, length int }
	blocks := make([]block, 0, n)
	for i := 0; i < n; i++ {
		start := A[i]
		j := i
		for j+1 < n && A[j+1] == A[j]+1 {
			j++
		}
		blocks = append(blocks, block{a: start, c: A[j], length: j - i + 1})
		i = j
	}
	total := 0
	infL := -2000000000
	infR := 2000000000
	for i, b := range blocks {
		L := infL
		R := infR
		if i > 0 {
			L = blocks[i-1].c + 1
		}
		if i+1 < len(blocks) {
			R = blocks[i+1].a - 1
		}
		l := sort.Search(len(B), func(j int) bool { return B[j] >= L })
		r := sort.Search(len(B), func(j int) bool { return B[j] > R }) - 1
		if l >= len(B) || r < l {
			continue
		}
		best := 0
		right := l
		maxDist := b.length - 1
		for left := l; left <= r; left++ {
			for right <= r && B[right]-B[left] <= maxDist {
				right++
			}
			cnt := right - left
			if cnt > best {
				best = cnt
			}
		}
		total += best
	}
	return fmt.Sprintf("%d", total)
}

func runCase(bin, input, expected string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(5) + 1
		m := rng.Intn(5) + 1
		setA := make(map[int]struct{})
		for len(setA) < n {
			setA[rng.Intn(20)] = struct{}{}
		}
		setB := make(map[int]struct{})
		for len(setB) < m {
			setB[rng.Intn(20)] = struct{}{}
		}
		A := make([]int, 0, n)
		for v := range setA {
			A = append(A, v)
		}
		B := make([]int, 0, m)
		for v := range setB {
			B = append(B, v)
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i2, v := range A {
			if i2 > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		for i2, v := range B {
			if i2 > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		input := sb.String()
		exp := expected(len(A), len(B), A, B)
		if err := runCase(bin, input, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
