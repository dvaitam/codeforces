package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func run(bin string, input string) (string, error) {
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

func expectedValue(s string, p float64) float64 {
	count1 := 0
	countQ := 0
	for _, ch := range s {
		if ch == '1' {
			count1++
		} else if ch == '?' {
			countQ++
		}
	}
	n := float64(len(s))
	return (float64(count1) + float64(countQ)*p) / n
}

func runCase(bin string, s string, p float64) error {
	input := fmt.Sprintf("%s\n%f\n", s, p)
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	val, err := strconv.ParseFloat(out, 64)
	if err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	exp := expectedValue(s, p)
	if math.Abs(val-exp) > 1e-5 {
		return fmt.Errorf("expected %.5f got %.5f", exp, val)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(50) + 1
		var sb strings.Builder
		for j := 0; j < n; j++ {
			r := rng.Intn(3)
			switch r {
			case 0:
				sb.WriteByte('0')
			case 1:
				sb.WriteByte('1')
			default:
				sb.WriteByte('?')
			}
		}
		p := rng.Float64()
		if err := runCase(bin, sb.String(), p); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v (s=%s p=%f)\n", i+1, err, sb.String(), p)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
