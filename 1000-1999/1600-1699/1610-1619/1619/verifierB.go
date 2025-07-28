package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("time limit")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

func intSqrt(n int64) int64 {
	x := int64(math.Sqrt(float64(n)))
	for (x+1)*(x+1) <= n {
		x++
	}
	for x*x > n {
		x--
	}
	return x
}

func intCbrt(n int64) int64 {
	x := int64(math.Cbrt(float64(n)))
	for (x+1)*(x+1)*(x+1) <= n {
		x++
	}
	for x*x*x > n {
		x--
	}
	return x
}

func intSixthRoot(n int64) int64 {
	x := int64(math.Pow(float64(n), 1.0/6.0))
	pow6 := func(v int64) int64 {
		res := int64(1)
		for i := 0; i < 6; i++ {
			res *= v
		}
		return res
	}
	for pow6(x+1) <= n {
		x++
	}
	for pow6(x) > n {
		x--
	}
	return x
}

func generateCase(rng *rand.Rand) (string, string) {
	t := rng.Intn(10) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	var out strings.Builder
	for i := 0; i < t; i++ {
		n := rng.Int63n(1_000_000_000) + 1
		fmt.Fprintf(&sb, "%d\n", n)
		sq := intSqrt(n)
		cb := intCbrt(n)
		sixth := intSixthRoot(n)
		fmt.Fprintf(&out, "%d\n", sq+cb-sixth)
	}
	return sb.String(), strings.TrimSpace(out.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runBinary(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
