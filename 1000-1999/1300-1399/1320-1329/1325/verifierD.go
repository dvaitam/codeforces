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

func solveD(u, v int64) []int64 {
	diff := v - u
	if diff < 0 || diff&1 == 1 {
		return []int64{-1}
	}
	if diff == 0 {
		if u == 0 {
			return []int64{0}
		}
		return []int64{1, u}
	}
	diff >>= 1
	if diff&u == 0 {
		return []int64{2, diff, diff ^ u}
	}
	return []int64{3, diff, diff, u}
}

func verifyCase(bin string, u, v int64) error {
	input := fmt.Sprintf("%d %d\n", u, v)
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	fields := strings.Fields(out)
	expected := solveD(u, v)
	if int64(len(fields)) != int64(len(expected)) {
		return fmt.Errorf("expected %d numbers got %d", len(expected), len(fields))
	}
	for i, f := range fields {
		got, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid number %q", f)
		}
		if got != expected[i] {
			return fmt.Errorf("expected %v got %v", expected, fields)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		u := rng.Int63n(1_000_000_000)
		diff := rng.Int63n(1_000_000_000)
		v := u + diff
		if err := verifyCase(bin, u, v); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
