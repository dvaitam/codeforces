package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func run(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout")
		}
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func checkCase(a, b int64, output string) error {
	tokens := strings.Fields(output)
	if b == 1 {
		if len(tokens) != 1 || strings.ToUpper(tokens[0]) != "NO" {
			return fmt.Errorf("expected NO")
		}
		return nil
	}
	if len(tokens) < 4 || strings.ToUpper(tokens[0]) != "YES" {
		return fmt.Errorf("expected YES x y z")
	}
	x, err := strconv.ParseInt(tokens[1], 10, 64)
	if err != nil {
		return fmt.Errorf("bad x")
	}
	y, err := strconv.ParseInt(tokens[2], 10, 64)
	if err != nil {
		return fmt.Errorf("bad y")
	}
	z, err := strconv.ParseInt(tokens[3], 10, 64)
	if err != nil {
		return fmt.Errorf("bad z")
	}
	if x <= 0 || y <= 0 || z <= 0 {
		return fmt.Errorf("numbers must be positive")
	}
	if x == y || x == z || y == z {
		return fmt.Errorf("numbers must be distinct")
	}
	good := func(v int64) bool { return v%(a*b) == 0 }
	nearly := func(v int64) bool { return v%a == 0 && !good(v) }
	cntGood := 0
	if good(x) {
		cntGood++
	}
	if good(y) {
		cntGood++
	}
	if good(z) {
		cntGood++
	}
	if cntGood != 1 {
		return fmt.Errorf("exactly one number must be good")
	}
	if !((good(x) && nearly(y) && nearly(z)) ||
		(good(y) && nearly(x) && nearly(z)) ||
		(good(z) && nearly(x) && nearly(y))) {
		return fmt.Errorf("nearly good condition failed")
	}
	if x+y != z && x+z != y && y+z != x {
		return fmt.Errorf("sum condition failed")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(42)
	for i := 1; i <= 100; i++ {
		a := rand.Int63n(1_000_000) + 1
		b := rand.Int63n(1_000_000) + 1
		if i%10 == 0 {
			b = 1
		}
		input := fmt.Sprintf("1\n%d %d\n", a, b)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\n", i, err)
			os.Exit(1)
		}
		if err := checkCase(a, b, out); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\ninput:%soutput:%s\n", i, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
