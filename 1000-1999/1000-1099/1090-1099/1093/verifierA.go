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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runCase(bin string, x int) error {
	input := fmt.Sprintf("1\n%d\n", x)
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	tokens := strings.Fields(out)
	if len(tokens) != 1 {
		return fmt.Errorf("expected one number, got %d", len(tokens))
	}
	n, err := strconv.Atoi(tokens[0])
	if err != nil {
		return fmt.Errorf("invalid integer %q", tokens[0])
	}
	if n <= 0 || 2*n > x || 7*n < x {
		return fmt.Errorf("invalid number of rolls %d for x=%d", n, x)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		x := rng.Intn(99) + 2 // 2..100
		if err := runCase(bin, x); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
