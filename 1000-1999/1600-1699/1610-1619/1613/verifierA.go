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

func digits(x int64) int {
	d := 0
	for x > 0 {
		d++
		x /= 10
	}
	return d
}

func pow10(k int) int64 {
	res := int64(1)
	for i := 0; i < k; i++ {
		res *= 10
	}
	return res
}

func compare(x1, p1, x2, p2 int64) string {
	d1 := digits(x1)
	d2 := digits(x2)
	len1 := d1 + int(p1)
	len2 := d2 + int(p2)
	if len1 > len2 {
		return ">"
	}
	if len1 < len2 {
		return "<"
	}
	diff := int64(d1 - d2)
	if diff > 0 {
		x2 *= pow10(int(diff))
		p2 -= diff
	} else if diff < 0 {
		x1 *= pow10(int(-diff))
		p1 += diff
	}
	if x1 > x2 {
		return ">"
	}
	if x1 < x2 {
		return "<"
	}
	return "="
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) (string, string) {
	x1 := rng.Int63n(1_000_000) + 1
	p1 := rng.Int63n(1_000_000)
	x2 := rng.Int63n(1_000_000) + 1
	p2 := rng.Int63n(1_000_000)
	input := fmt.Sprintf("1\n%d %d\n%d %d\n", x1, p1, x2, p2)
	expect := compare(x1, p1, x2, p2)
	return input, expect
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
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
