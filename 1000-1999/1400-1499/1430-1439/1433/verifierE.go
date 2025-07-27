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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(n int) uint64 {
	k := n / 2
	comb := uint64(1)
	for i := 1; i <= k; i++ {
		comb = comb * uint64(n-k+i) / uint64(i)
	}
	fact := uint64(1)
	for i := 1; i <= k-1; i++ {
		fact *= uint64(i)
	}
	return comb * fact * fact / 2
}

func genCase(rng *rand.Rand) int {
	return 2 * (rng.Intn(10) + 1)
}

func runCase(bin string, n int) error {
	input := fmt.Sprintf("%d\n", n)
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	got, err := strconv.ParseUint(strings.TrimSpace(out), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse output: %v", err)
	}
	want := expected(n)
	if got != want {
		return fmt.Errorf("expected %d got %d", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := genCase(rng)
		if err := runCase(bin, n); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
