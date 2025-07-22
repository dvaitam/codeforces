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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runCase(bin string, n int) error {
	input := fmt.Sprintf("%d\n", n)
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	expected := strconv.FormatInt(int64(n), 2)
	if out != expected {
		return fmt.Errorf("expected %s got %s", expected, out)
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
	edge := []int{1, 2, 3, 7, 8, 15, 16, 255, 256, 1000000}
	idx := 0
	for ; idx < len(edge); idx++ {
		n := edge[idx]
		if err := runCase(bin, n); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v (n=%d)\n", idx+1, err, n)
			os.Exit(1)
		}
	}
	for ; idx < 100; idx++ {
		n := rng.Intn(1000000) + 1
		if err := runCase(bin, n); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v (n=%d)\n", idx+1, err, n)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
