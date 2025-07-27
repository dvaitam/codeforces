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

func runProg(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
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

func expectedX(n int64) int64 {
	for k := int64(2); ; k++ {
		denom := (int64(1) << k) - 1
		if n%denom == 0 {
			return n / denom
		}
	}
}

func genCase(rng *rand.Rand) int64 {
	k := rng.Intn(28) + 2 // 2..29
	denom := int64(1<<k) - 1
	maxX := int64(1e9) / denom
	if maxX < 1 {
		maxX = 1
	}
	x := rng.Int63n(maxX) + 1
	return x * denom
}

func runCase(bin string, n int64) error {
	input := fmt.Sprintf("1\n%d\n", n)
	got, err := runProg(bin, input)
	if err != nil {
		return err
	}
	x, err := strconv.ParseInt(strings.TrimSpace(got), 10, 64)
	if err != nil {
		return fmt.Errorf("cannot parse output: %v", err)
	}
	exp := expectedX(n)
	if x != exp {
		return fmt.Errorf("expected %d got %d", exp, x)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
