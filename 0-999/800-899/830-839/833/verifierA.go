package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func cbrtInt(x int64) int64 {
	r := int64(math.Round(math.Cbrt(float64(x))))
	for (r+1)*(r+1)*(r+1) <= x {
		r++
	}
	for r*r*r > x {
		r--
	}
	return r
}

func checkPair(a, b int64) string {
	prod := a * b
	c := cbrtInt(prod)
	if c*c*c != prod || a%c != 0 || b%c != 0 {
		return "No"
	}
	x := a / c
	y := b / c
	if x*y == c {
		return "Yes"
	}
	return "No"
}

func solveCase(n int, pairs [][2]int64) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(checkPair(pairs[i][0], pairs[i][1]))
	}
	return sb.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	pairs := make([][2]int64, n)
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		var a, b int64
		if rng.Intn(2) == 0 { // generate valid case
			x := int64(rng.Intn(10) + 1)
			y := int64(rng.Intn(10) + 1)
			c := x * y
			a = x * c
			b = y * c
		} else {
			a = int64(rng.Intn(1000000) + 1)
			b = int64(rng.Intn(1000000) + 1)
		}
		pairs[i][0] = a
		pairs[i][1] = b
		input.WriteString(fmt.Sprintf("%d %d\n", a, b))
	}
	expected := solveCase(n, pairs)
	if !strings.HasSuffix(expected, "\n") {
		expected += "\n"
	}
	return input.String(), expected
}

func runCase(bin, in, exp string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	if outStr != strings.TrimSpace(exp) {
		return fmt.Errorf("expected %q got %q", exp, outStr)
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
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
