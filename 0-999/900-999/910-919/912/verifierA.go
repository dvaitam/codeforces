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

func runCandidate(bin, input string) (string, error) {
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

func solveCase(A, B, x, y, z int64) string {
	needYellow := 2*x + y
	needBlue := y + 3*z
	add := int64(0)
	if needYellow > A {
		add += needYellow - A
	}
	if needBlue > B {
		add += needBlue - B
	}
	return fmt.Sprint(add)
}

func generateCase(rng *rand.Rand) (int64, int64, int64, int64, int64) {
	A := rng.Int63n(1000)
	B := rng.Int63n(1000)
	x := rng.Int63n(500)
	y := rng.Int63n(500)
	z := rng.Int63n(500)
	return A, B, x, y, z
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		A, B, x, y, z := generateCase(rng)
		input := fmt.Sprintf("%d %d\n%d %d %d\n", A, B, x, y, z)
		expect := solveCase(A, B, x, y, z)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if out != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
