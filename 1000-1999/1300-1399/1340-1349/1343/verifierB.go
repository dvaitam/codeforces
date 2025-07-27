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

func runProg(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
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

func genCase(rng *rand.Rand) int {
	return 2 * (rng.Intn(100) + 1) // even between 2 and 200
}

func runCase(bin string, n int) error {
	input := fmt.Sprintf("1\n%d\n", n)
	out, err := runProg(bin, input)
	if err != nil {
		return err
	}
	tokens := strings.Fields(out)
	if len(tokens) == 0 {
		return fmt.Errorf("no output")
	}
	needNO := (n/2)%2 == 1
	ans := strings.ToUpper(tokens[0])
	if needNO {
		if ans != "NO" {
			return fmt.Errorf("expected NO for n=%d", n)
		}
		return nil
	}
	if ans != "YES" {
		return fmt.Errorf("expected YES for n=%d", n)
	}
	if len(tokens)-1 != n {
		return fmt.Errorf("expected %d numbers got %d", n, len(tokens)-1)
	}
	seen := make(map[int]bool)
	sumEven := 0
	sumOdd := 0
	for i := 0; i < n; i++ {
		v, err := strconv.Atoi(tokens[i+1])
		if err != nil {
			return fmt.Errorf("invalid number: %v", err)
		}
		if v <= 0 {
			return fmt.Errorf("non positive value %d", v)
		}
		if seen[v] {
			return fmt.Errorf("duplicate value %d", v)
		}
		seen[v] = true
		if i < n/2 {
			if v%2 != 0 {
				return fmt.Errorf("first half not even")
			}
			sumEven += v
		} else {
			if v%2 == 0 {
				return fmt.Errorf("second half not odd")
			}
			sumOdd += v
		}
	}
	if sumEven != sumOdd {
		return fmt.Errorf("sum mismatch %d vs %d", sumEven, sumOdd)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := genCase(rng)
		if err := runCase(bin, n); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
