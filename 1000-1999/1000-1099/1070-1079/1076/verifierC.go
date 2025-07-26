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

func solveCase(d float64) string {
	D := d*d - 4*d
	if D < 0 {
		return "N"
	}
	sqrtD := math.Sqrt(D)
	a := (d + sqrtD) / 2
	b := (d - sqrtD) / 2
	return fmt.Sprintf("Y %.15f %.15f", a, b)
}

func generateCase(rng *rand.Rand) (string, string) {
	t := rng.Intn(5) + 1
	var in strings.Builder
	in.WriteString(fmt.Sprintf("%d\n", t))
	var out strings.Builder
	for i := 0; i < t; i++ {
		d := float64(rng.Intn(1001))
		in.WriteString(fmt.Sprintf("%d\n", int(d)))
		if i > 0 {
			out.WriteByte('\n')
		}
		out.WriteString(solveCase(d))
	}
	return in.String(), out.String()
}

func runCase(bin, in, expect string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	expect = strings.TrimSpace(expect)
	if got != expect {
		return fmt.Errorf("expected:\n%s\n\ngot:\n%s", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := generateCase(rng)
		if err := runCase(bin, in, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
