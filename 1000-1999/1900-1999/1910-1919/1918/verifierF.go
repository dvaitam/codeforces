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

func generateCase(rng *rand.Rand) (int, int, []int) {
	n := rng.Intn(10) + 2
	k := rng.Intn(10)
	parent := make([]int, n-1)
	for i := 2; i <= n; i++ {
		parent[i-2] = rng.Intn(i-1) + 1
	}
	return n, k, parent
}

func runCase(bin string, n, k int, parent []int) error {
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i, p := range parent {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(fmt.Sprintf("%d", p))
	}
	input.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	if got != "0" {
		return fmt.Errorf("expected 0 got %s", got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, k, parent := generateCase(rng)
		if err := runCase(bin, n, k, parent); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
