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

func expected(arr []int) []int {
	n := len(arr)
	ans := make([]int, n+1)
	for l := 0; l < n; l++ {
		cnt := make([]int, n+1)
		best := 0
		for r := l; r < n; r++ {
			c := arr[r]
			cnt[c]++
			if cnt[c] > cnt[best] || (cnt[c] == cnt[best] && c < best) {
				best = c
			}
			ans[best]++
		}
	}
	return ans[1:]
}

func runCase(exe string, input string, exp []int) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) != len(exp) {
		return fmt.Errorf("expected %d numbers, got %d", len(exp), len(fields))
	}
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("bad integer %q", f)
		}
		if v != exp[i] {
			return fmt.Errorf("mismatch at %d: expected %d got %d", i, exp[i], v)
		}
	}
	return nil
}

func generateCase(rng *rand.Rand) (string, []int) {
	n := rng.Intn(9) + 1 // up to 10
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(n) + 1
	}
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	exp := expected(arr)
	return sb.String(), exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
