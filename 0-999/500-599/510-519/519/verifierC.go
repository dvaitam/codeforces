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

func solveCase(input string) string {
	reader := strings.NewReader(input)
	var n, m int64
	fmt.Fscan(reader, &n, &m)
	total := (n + m) / 3
	ans := total
	if n < ans {
		ans = n
	}
	if m < ans {
		ans = m
	}
	return fmt.Sprint(ans)
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(100) + 1
	m := rng.Intn(100) + 1
	return fmt.Sprintf("%d %d\n", n, m)
}

func runCase(bin, input, expected string) error {
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
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %q got %q", expected, got)
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
		inp := generateCase(rng)
		exp := solveCase(inp)
		if err := runCase(bin, inp, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, inp)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
