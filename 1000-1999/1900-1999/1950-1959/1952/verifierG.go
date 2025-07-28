package main

import (
	"bytes"
	"fmt"
	"math"
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) (string, float64) {
	x := rng.Float64()*1000 + 0.1
	return fmt.Sprintf("%.6f\n", x), math.Log2(x)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
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
		got, err := strconv.ParseFloat(out, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse output %q\n", i+1, out)
			os.Exit(1)
		}
		if math.Abs(got-exp) > 1e-6 {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %.6f got %.6f\n", i+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
