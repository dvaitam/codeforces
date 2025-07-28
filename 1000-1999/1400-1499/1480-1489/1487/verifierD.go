package main

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func solveD(n int64) int64 {
	m := int64(math.Sqrt(float64(2*n - 1)))
	ans := (m - 1) / 2
	if ans < 0 {
		ans = 0
	}
	return ans
}

func run(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	for t := 0; t < 100; t++ {
		n := rng.Int63n(1e9) + 1
		input := fmt.Sprintf("1\n%d\n", n)
		expected := strconv.FormatInt(solveD(n), 10)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n%s", t+1, err, got)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("case %d failed:\ninput:%sexpected %s got %s\n", t+1, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
