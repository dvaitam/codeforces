package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type BIT struct {
	n    int
	tree []int64
}

func newBIT(n int) *BIT {
	return &BIT{n: n, tree: make([]int64, n+2)}
}

func (b *BIT) add(idx int, val int64) {
	for idx <= b.n {
		b.tree[idx] += val
		idx += idx & -idx
	}
}

func (b *BIT) sum(idx int) int64 {
	var s int64
	for idx > 0 {
		s += b.tree[idx]
		idx -= idx & -idx
	}
	return s
}

func expected(arr []int, k int) int64 {
	if k == 0 {
		return int64(len(arr))
	}
	bits := make([]*BIT, k+1)
	for i := range bits {
		bits[i] = newBIT(len(arr) + 2)
	}
	for _, x := range arr {
		bits[0].add(x, 1)
		for j := k; j >= 1; j-- {
			val := bits[j-1].sum(x - 1)
			if val != 0 {
				bits[j].add(x, val)
			}
		}
	}
	return bits[k].sum(len(arr) + 1)
}

func generateCase(rng *rand.Rand) (string, int64) {
	n := rng.Intn(50) + 1
	k := rng.Intn(5)
	if k > n-1 {
		k = n - 1
	}
	arr := rng.Perm(n)
	// permutation returns 0..n-1; shift by 1
	for i := range arr {
		arr[i]++
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String(), expected(arr, k)
}

func runCase(exe, input string, expected int64) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	got, err := strconv.ParseInt(outStr, 10, 64)
	if err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
