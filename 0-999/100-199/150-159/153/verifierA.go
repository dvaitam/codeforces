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

func runCase(bin string, a, b int) error {
	input := fmt.Sprintf("%d %d\n", a, b)
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	var got int
	if _, err := fmt.Sscan(out, &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != a+b {
		return fmt.Errorf("expected %d got %d", a+b, got)
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
	edge := [][2]int{{1, 1}, {1, 100000}, {100000, 1}, {100000, 100000}, {2, 3}, {99999, 12345}}
	idx := 0
	for ; idx < len(edge); idx++ {
		a, b := edge[idx][0], edge[idx][1]
		if err := runCase(bin, a, b); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v (a=%d b=%d)\n", idx+1, err, a, b)
			os.Exit(1)
		}
	}
	for ; idx < 100; idx++ {
		a := rng.Intn(100000) + 1
		b := rng.Intn(100000) + 1
		if err := runCase(bin, a, b); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v (a=%d b=%d)\n", idx+1, err, a, b)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
