package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

func expectedA(n int, k int64, arr []int64) int {
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	used := make(map[int64]bool, n)
	cnt := 0
	for _, v := range arr {
		if k != 0 && v%k == 0 {
			if used[v/k] {
				continue
			}
		}
		used[v] = true
		cnt++
	}
	return cnt
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(15) + 1
	k := int64(rng.Intn(5) + 1)
	vals := make([]int64, 0, n)
	used := map[int64]bool{}
	for len(vals) < n {
		v := int64(rng.Intn(30) + 1)
		if !used[v] {
			used[v] = true
			vals = append(vals, v)
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i, v := range vals {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	exp := expectedA(n, k, append([]int64(nil), vals...))
	return sb.String(), fmt.Sprintf("%d\n", exp)
}

func runCase(exe string, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
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
