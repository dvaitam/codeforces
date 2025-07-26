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

func solve(h1, m1, h2, m2 int) string {
	t1 := h1*60 + m1
	t2 := h2*60 + m2
	mid := (t1 + t2) / 2
	return fmt.Sprintf("%02d:%02d", mid/60, mid%60)
}

func generateCase(rng *rand.Rand) (string, string) {
	start := rng.Intn(24*60 - 2)
	length := rng.Intn(60)*2 + 2
	if start+length >= 24*60 {
		length = (24*60 - start - 2) / 2 * 2
		if length < 2 {
			length = 2
		}
	}
	end := start + length
	h1, m1 := start/60, start%60
	h2, m2 := end/60, end%60
	input := fmt.Sprintf("%02d:%02d\n%02d:%02d\n", h1, m1, h2, m2)
	return input, solve(h1, m1, h2, m2)
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

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d mismatch: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
