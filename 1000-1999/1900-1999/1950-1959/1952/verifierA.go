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

func runCandidate(bin, input string) (string, error) {
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

func generateCase(rng *rand.Rand) (string, int) {
	a := rng.Intn(2001) - 1000
	b := rng.Intn(2001) - 1000
	input := fmt.Sprintf("%d %d\n", a, b)
	return input, a + b
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		val, err := strconv.Atoi(out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse output %q\n", i+1, out)
			os.Exit(1)
		}
		if val != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", i+1, expect, val)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
