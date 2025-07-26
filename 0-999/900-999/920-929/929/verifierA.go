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

func expectedA(n, k int, x []int) int {
	cnt := 0
	i := 0
	for i < n-1 {
		j := i
		for j+1 < n && x[j+1]-x[i] <= k {
			j++
		}
		if j == i {
			return -1
		}
		cnt++
		i = j
	}
	return cnt
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(9) + 2 // 2..10
	k := rng.Intn(20) + 1
	x := make([]int, n)
	cur := rng.Intn(5)
	x[0] = cur
	for i := 1; i < n; i++ {
		cur += rng.Intn(5) + 1
		x[i] = cur
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i, v := range x {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	expect := expectedA(n, k, x)
	return sb.String(), fmt.Sprintf("%d\n", expect)
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(expected), got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
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
