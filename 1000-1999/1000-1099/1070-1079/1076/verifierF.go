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

func solve(n int, k int64, x, y []int64) string {
	var u, v int64
	for i := 0; i < n; i++ {
		a := x[i]
		b := y[i]
		if (b+1)*k-u < a || (a+1)*k-v < b {
			return "NO"
		}
		if b*k-u < a {
			u = a - b*k + u
		} else {
			u = 0
		}
		if a*k-v < b {
			v = b - a*k + v
		} else {
			v = 0
		}
	}
	return "YES"
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	k := int64(rng.Intn(5) + 1)
	x := make([]int64, n)
	y := make([]int64, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i < n; i++ {
		x[i] = int64(rng.Intn(10) + 1)
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", x[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		y[i] = int64(rng.Intn(10) + 1)
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", y[i]))
	}
	sb.WriteByte('\n')
	expect := solve(n, k, x, y)
	return sb.String(), expect
}

func runCase(bin, in, expect string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expect {
		return fmt.Errorf("expected %s got %s", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := generateCase(rng)
		if err := runCase(bin, in, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
