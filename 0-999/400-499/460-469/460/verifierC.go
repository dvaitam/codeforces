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

func canReach(a []int64, m, w int, h int64) bool {
	n := len(a)
	ops := make([]int64, n)
	var added, used int64
	for i := 0; i < n; i++ {
		if i >= w {
			added -= ops[i-w]
		}
		curr := a[i] + added
		if curr < h {
			need := h - curr
			used += need
			if used > int64(m) {
				return false
			}
			added += need
			ops[i] = need
		}
	}
	return true
}

func expectedC(a []int64, m, w int) int64 {
	minA := a[0]
	for _, v := range a {
		if v < minA {
			minA = v
		}
	}
	low := minA
	high := minA + int64(m) + 1
	ans := minA
	for low <= high {
		mid := (low + high) / 2
		if canReach(a, m, w, mid) {
			ans = mid
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return ans
}

func runCase(bin string, n, m, w int, arr []int64) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, w))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteString("\n")
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	expect := expectedC(arr, m, w)
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
}

func generateCase(rng *rand.Rand) (int, int, int, []int64) {
	n := rng.Intn(20) + 1
	w := rng.Intn(n) + 1
	m := rng.Intn(30) + 1
	arr := make([]int64, n)
	for i := range arr {
		arr[i] = int64(rng.Intn(20) + 1)
	}
	return n, m, w, arr
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// some edge cases
	edges := []struct {
		n, m, w int
		arr     []int64
	}{
		{1, 1, 1, []int64{1}},
		{3, 5, 2, []int64{1, 2, 3}},
		{5, 10, 3, []int64{5, 5, 5, 5, 5}},
	}
	for i, e := range edges {
		if err := runCase(bin, e.n, e.m, e.w, e.arr); err != nil {
			fmt.Fprintf(os.Stderr, "edge case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}

	for i := 0; i < 100; i++ {
		n, m, w, arr := generateCase(rng)
		if err := runCase(bin, n, m, w, arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
