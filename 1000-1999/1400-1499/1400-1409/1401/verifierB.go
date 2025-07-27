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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(x1, y1, z1, x2, y2, z2 int64) int64 {
	pos := min(z1, y2)
	safe := x2 + y2
	neg := y1 - safe
	if neg < 0 {
		neg = 0
	}
	return 2*pos - 2*neg
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func runCase(bin string, x1, y1, z1, x2, y2, z2 int64) error {
	input := fmt.Sprintf("1\n%d %d %d\n%d %d %d\n", x1, y1, z1, x2, y2, z2)
	gotStr, err := run(bin, input)
	if err != nil {
		return err
	}
	var got int64
	if _, err := fmt.Sscan(gotStr, &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	exp := solve(x1, y1, z1, x2, y2, z2)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
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
		x1 := rng.Int63n(1000)
		y1 := rng.Int63n(1000)
		z1 := rng.Int63n(1000)
		x2 := rng.Int63n(1000)
		y2 := rng.Int63n(1000)
		z2 := rng.Int63n(1000)
		if err := runCase(bin, x1, y1, z1, x2, y2, z2); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
