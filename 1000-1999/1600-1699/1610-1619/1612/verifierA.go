package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
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
	return strings.TrimSpace(out.String()), nil
}

func check(x, y int, output string) error {
	parts := strings.Fields(output)
	if len(parts) != 2 {
		return fmt.Errorf("expected two integers, got %q", output)
	}
	a, err1 := strconv.Atoi(parts[0])
	b, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		return fmt.Errorf("failed to parse integers from %q", output)
	}
	if (x+y)%2 == 1 {
		if a != -1 || b != -1 {
			return fmt.Errorf("expected -1 -1 for x=%d y=%d, got %d %d", x, y, a, b)
		}
		return nil
	}
	if a == -1 && b == -1 {
		return fmt.Errorf("solution exists for x=%d y=%d but got -1 -1", x, y)
	}
	if a < 0 || b < 0 || a > x || b > y {
		return fmt.Errorf("invalid coordinates %d %d for x=%d y=%d", a, b, x, y)
	}
	if a+b != (x+y)/2 {
		return fmt.Errorf("a+b mismatch for x=%d y=%d: got %d+%d", x, y, a, b)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for x := 0; x <= 50; x++ {
		for y := 0; y <= 50; y++ {
			input := fmt.Sprintf("1\n%d %d\n", x, y)
			out, err := run(bin, input)
			if err != nil {
				fmt.Fprintf(os.Stderr, "case x=%d y=%d failed: %v\n", x, y, err)
				os.Exit(1)
			}
			if err := check(x, y, out); err != nil {
				fmt.Fprintf(os.Stderr, "case x=%d y=%d failed: %v\n", x, y, err)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
