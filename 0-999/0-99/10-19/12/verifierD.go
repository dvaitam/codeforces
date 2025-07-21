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

func solveNaive(B, I, R []int) int {
	n := len(B)
	cnt := 0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if B[j] > B[i] && I[j] > I[i] && R[j] > R[i] {
				cnt++
				break
			}
		}
	}
	return cnt
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		n := rng.Intn(20) + 1
		B := make([]int, n)
		I := make([]int, n)
		Rr := make([]int, n)
		for i := 0; i < n; i++ {
			B[i] = rng.Intn(100)
			I[i] = rng.Intn(100)
			Rr[i] = rng.Intn(100)
		}
		var input strings.Builder
		fmt.Fprintf(&input, "%d\n", n)
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(B[i]))
		}
		input.WriteByte('\n')
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(I[i]))
		}
		input.WriteByte('\n')
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(Rr[i]))
		}
		input.WriteByte('\n')
		expected := strconv.Itoa(solveNaive(B, I, Rr))
		out, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", t+1, err, input.String())
			os.Exit(1)
		}
		if out != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", t+1, expected, out, input.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
