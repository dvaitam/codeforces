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

func win(a, b uint64) bool {
	reversed := false
	for {
		if a == 0 || b == 0 {
			return reversed
		}
		if a > b {
			a, b = b, a
		}
		if b/a > 1 {
			return !reversed
		}
		b = b % a
		reversed = !reversed
	}
}

func solve(t int, pairs [][2]uint64) []string {
	res := make([]string, t)
	for i := 0; i < t; i++ {
		if win(pairs[i][0], pairs[i][1]) {
			res[i] = "First"
		} else {
			res[i] = "Second"
		}
	}
	return res
}

func generateCase(rng *rand.Rand) (string, []string) {
	t := rng.Intn(5) + 1
	pairs := make([][2]uint64, t)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		a := rng.Uint64()%1000 + 1
		b := rng.Uint64()%1000 + 1
		pairs[i][0] = a
		pairs[i][1] = b
		sb.WriteString(fmt.Sprintf("%d %d\n", a, b))
	}
	exp := solve(t, pairs)
	return sb.String(), exp
}

func run(bin, input string) ([]string, error) {
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
		return nil, fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	lines := strings.Fields(out.String())
	return lines, nil
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if len(out) != len(exp) {
			fmt.Fprintf(os.Stderr, "case %d wrong number of lines\ninput:\n%s", i+1, in)
			os.Exit(1)
		}
		for j := range exp {
			if out[j] != exp[j] {
				fmt.Fprintf(os.Stderr, "case %d mismatch on line %d: expected %s got %s\ninput:\n%s", i+1, j+1, exp[j], out[j], in)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
