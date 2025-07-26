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

var primes = []int64{3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73}

func expectedB(n int, p, k int64, arr []int64) int64 {
	freq := make(map[int64]int64, n)
	for _, a := range arr {
		a %= p
		a2 := (a * a) % p
		a4 := (a2 * a2) % p
		key := (a4 - (k*a)%p + p) % p
		freq[key]++
	}
	var ans int64
	for _, c := range freq {
		if c > 1 {
			ans += c * (c - 1) / 2
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	p := primes[rng.Intn(len(primes))]
	n := rng.Intn(int(p)) + 1
	if n > 10 {
		n = 10
	}
	k := rng.Int63n(p)
	used := make(map[int64]bool)
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		for {
			v := rng.Int63n(p)
			if !used[v] {
				arr[i] = v
				used[v] = true
				break
			}
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, p, k)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", arr[i])
	}
	sb.WriteByte('\n')
	expect := expectedB(n, p, k, arr)
	return sb.String(), fmt.Sprintf("%d", expect)
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
