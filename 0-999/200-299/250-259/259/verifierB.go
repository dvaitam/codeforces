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

func genCase(rng *rand.Rand) (string, string) {
	var b, c, d, f, g, h int
	var M int
	for {
		b = rng.Intn(50) + 1
		c = rng.Intn(50) + 1
		d = rng.Intn(50) + 1
		f = rng.Intn(50) + 1
		g = rng.Intn(50) + 1
		h = rng.Intn(50) + 1
		sum := b + c + d + f + g + h
		if sum%2 != 0 {
			continue
		}
		M = sum / 2
		x1 := M - b - c
		x2 := M - d - f
		x3 := M - g - h
		if x1 > 0 && x2 > 0 && x3 > 0 {
			break
		}
	}
	sum := b + c + d + f + g + h
	M = sum / 2
	x1 := M - b - c
	x2 := M - d - f
	x3 := M - g - h

	var in strings.Builder
	fmt.Fprintf(&in, "0 %d %d\n", b, c)
	fmt.Fprintf(&in, "%d 0 %d\n", d, f)
	fmt.Fprintf(&in, "%d %d 0\n", g, h)

	var out strings.Builder
	fmt.Fprintf(&out, "%d %d %d\n", x1, b, c)
	fmt.Fprintf(&out, "%d %d %d\n", d, x2, f)
	fmt.Fprintf(&out, "%d %d %d\n", g, h, x3)

	return in.String(), out.String()
}

func runCandidate(bin, input, expected string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected\n%s\ngot\n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCandidate(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
