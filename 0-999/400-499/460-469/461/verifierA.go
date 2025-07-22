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

type testCase struct {
	n   int
	arr []int64
}

func solveCase(a []int64) int64 {
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	var sum int64
	for _, v := range a {
		sum += v
	}
	ans := sum
	n := len(a)
	for i, v := range a {
		c := i + 1
		if c > n-1 {
			c = n - 1
		}
		ans += v * int64(c)
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(50) + 1
	arr := make([]int64, n)
	for i := range arr {
		arr[i] = rng.Int63n(1_000_000) + 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	ans := solveCase(append([]int64(nil), arr...))
	return sb.String(), fmt.Sprintf("%d", ans)
}

func runCase(bin string, input, expected string) error {
	cmd := exec.Command(bin)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
