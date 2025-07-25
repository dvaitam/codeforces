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

func expected(n, c1, c2 int64, s string) int64 {
	adults := int64(strings.Count(s, "1"))
	if adults == 0 {
		return 0
	}
	if adults > n {
		adults = n
	}
	best := int64(1<<63 - 1)
	for k := int64(1); k <= adults; k++ {
		q := n / k
		r := n % k
		cost := r*(c1+c2*q*q) + (k-r)*(c1+c2*(q-1)*(q-1))
		if cost < best {
			best = cost
		}
	}
	return best
}

func run(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) (string, int64) {
	n := int64(rng.Intn(20) + 1)
	if rng.Float64() < 0.1 {
		n = int64(rng.Intn(200) + 1)
	}
	c1 := int64(rng.Intn(10) + 1)
	c2 := int64(rng.Intn(10) + 1)
	sb := make([]byte, n)
	hasAdult := false
	for i := int64(0); i < n; i++ {
		if rng.Intn(2) == 0 {
			sb[i] = '0'
		} else {
			sb[i] = '1'
			hasAdult = true
		}
	}
	if !hasAdult {
		sb[rng.Intn(int(n))] = '1'
	}
	s := string(sb)
	input := fmt.Sprintf("%d %d %d\n%s\n", n, c1, c2, s)
	return input, expected(n, c1, c2, s)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := genCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: cannot parse output: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
