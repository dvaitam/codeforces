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

func expected(d, L, v1, v2 int) float64 {
	return float64(L-d) / float64(v1+v2)
}

func generateCase(rng *rand.Rand) (string, float64) {
	d := rng.Intn(9999) + 1
	L := d + rng.Intn(10000-d) + 1
	v1 := rng.Intn(10000) + 1
	v2 := rng.Intn(10000) + 1
	input := fmt.Sprintf("%d %d %d %d\n", d, L, v1, v2)
	return input, expected(d, L, v1, v2)
}

func runCase(bin, input string, exp float64) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got float64
	if _, err := fmt.Sscan(out.String(), &got); err != nil {
		return fmt.Errorf("failed to parse output: %v", err)
	}
	diff := math.Abs(got - exp)
	if diff > 1e-6*math.Max(1, math.Abs(exp)) {
		return fmt.Errorf("expected %.10f got %.10f", exp, got)
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
