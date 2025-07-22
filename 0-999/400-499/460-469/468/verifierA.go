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

func run(bin, input string) (string, error) {
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
	return out.String(), nil
}

func check(n int, output string) error {
	output = strings.TrimSpace(output)
	if output == "" {
		return fmt.Errorf("empty output")
	}
	lines := strings.Split(output, "\n")
	if n <= 3 {
		if len(lines) != 1 || strings.ToUpper(strings.TrimSpace(lines[0])) != "NO" {
			return fmt.Errorf("expected NO")
		}
		return nil
	}
	if strings.ToUpper(strings.TrimSpace(lines[0])) != "YES" {
		return fmt.Errorf("missing YES")
	}
	if len(lines) != n {
		return fmt.Errorf("expected %d operation lines, got %d", n-1, len(lines)-1)
	}
	counts := make(map[int64]int)
	for i := 1; i <= n; i++ {
		counts[int64(i)]++
	}
	for idx := 1; idx < len(lines); idx++ {
		line := strings.TrimSpace(lines[idx])
		var a, b, c int64
		var op string
		if _, err := fmt.Sscanf(line, "%d %s %d = %d", &a, &op, &b, &c); err != nil {
			return fmt.Errorf("line %d: cannot parse", idx)
		}
		if counts[a] == 0 {
			return fmt.Errorf("line %d: %d not available", idx, a)
		}
		counts[a]--
		if counts[b] == 0 {
			return fmt.Errorf("line %d: %d not available", idx, b)
		}
		counts[b]--
		var res int64
		switch op {
		case "+":
			res = a + b
		case "-":
			res = a - b
		case "*":
			res = a * b
		default:
			return fmt.Errorf("line %d: invalid operator", idx)
		}
		if res != c {
			return fmt.Errorf("line %d: incorrect result", idx)
		}
		if math.Abs(float64(res)) > 1e18 {
			return fmt.Errorf("line %d: result out of range", idx)
		}
		counts[c]++
	}
	if len(counts) != 1 || counts[24] != 1 {
		return fmt.Errorf("final value is not 24")
	}
	return nil
}

func generateCase(rng *rand.Rand) int {
	if rng.Intn(4) == 0 {
		return rng.Intn(4) + 1
	}
	return rng.Intn(100) + 1
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := generateCase(rng)
		input := fmt.Sprintf("%d\n", n)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := check(n, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: %d\noutput:\n%s\n", i+1, err, n, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
