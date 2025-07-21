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

func solve(input string) string {
	fields := strings.Fields(input)
	if len(fields) == 0 {
		return ""
	}
	n, _ := strconv.Atoi(fields[0])
	idx := 1
	var sx, sy, sx2, sy2 int64
	for i := 0; i < n; i++ {
		x, _ := strconv.ParseInt(fields[idx], 10, 64)
		y, _ := strconv.ParseInt(fields[idx+1], 10, 64)
		idx += 2
		sx += x
		sy += y
		sx2 += x * x
		sy2 += y * y
	}
	N := int64(n)
	ansX := N*sx2 - sx*sx
	ansY := N*sy2 - sy*sy
	ans := ansX + ansY
	return fmt.Sprintf("%d\n", ans)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(9) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		x := rng.Intn(21) - 10
		y := rng.Intn(21) - 10
		sb.WriteString(fmt.Sprintf("%d %d\n", x, y))
	}
	input := sb.String()
	expected := solve(strings.TrimSpace(input))
	return input, expected
}

func runCase(exe string, in, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(expected), strings.TrimSpace(out.String()))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
