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

// compute expected answer for a single query
func steps(n, k uint64) uint64 {
	var cnt uint64
	if k <= 1 {
		return n
	}
	for n > 0 {
		if n < k {
			cnt += n
			break
		}
		r := n % k
		if r != 0 {
			cnt += r
			n -= r
		} else {
			n /= k
			cnt++
		}
	}
	return cnt
}

func generateCase(rng *rand.Rand) (string, string) {
	t := rng.Intn(5) + 1
	var input strings.Builder
	fmt.Fprintf(&input, "%d\n", t)
	var output strings.Builder
	for i := 0; i < t; i++ {
		n := rng.Uint64()%1_000_000_000 + 1
		k := rng.Uint64()%10 + 2
		fmt.Fprintf(&input, "%d %d\n", n, k)
		fmt.Fprintf(&output, "%d\n", steps(n, k))
	}
	return input.String(), output.String()
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

func runCase(bin string, in, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(exp) {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", exp, out.String())
	}
	return nil
}
