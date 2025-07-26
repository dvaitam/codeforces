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

func runCandidate(bin, input string) (string, error) {
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

func generateCase(rng *rand.Rand) (string, int) {
	sign := 1
	if rng.Intn(2) == 0 {
		sign = -1
	}
	input := fmt.Sprintf("%d\n", sign)
	return input, sign
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, sign := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		var a, b float64
		if _, err := fmt.Sscan(out, &a, &b); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: bad output: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		inv := 1 / math.Sqrt2
		if math.Abs(a-inv) > 1e-6 || math.Abs(b-inv*float64(sign)) > 1e-6 {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %.10f %.10f got %s\ninput:\n%s", i+1, inv, inv*float64(sign), out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
