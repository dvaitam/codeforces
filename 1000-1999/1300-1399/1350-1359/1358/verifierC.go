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

func expected(x1, y1, x2, y2 int64) string {
	dx := x2 - x1
	dy := y2 - y1
	val := dx*dy + 1
	return fmt.Sprintf("%d", val)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		x1 := rng.Int63n(1_000_000_000) + 1
		y1 := rng.Int63n(1_000_000_000) + 1
		x2 := x1 + rng.Int63n(1_000_000)
		y2 := y1 + rng.Int63n(1_000_000)
		input := fmt.Sprintf("1\n%d %d %d %d\n", x1, y1, x2, y2)
		want := expected(x1, y1, x2, y2)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
