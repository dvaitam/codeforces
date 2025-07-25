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

const MOD int64 = 1000000007

func runCandidate(bin, input string) (string, error) {
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

func expected(n int, x int64, d []int64) string {
	counts := make([]int64, x+1)
	counts[0] = 1
	for dist := int64(0); dist <= x; dist++ {
		c := counts[dist]
		if c == 0 {
			continue
		}
		for _, step := range d {
			nd := dist + step
			if nd <= x {
				counts[nd] = (counts[nd] + c) % MOD
			}
		}
	}
	sum := int64(0)
	for i := int64(0); i <= x; i++ {
		sum = (sum + counts[i]) % MOD
	}
	return fmt.Sprintf("%d", sum%MOD)
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 1
	x := int64(rng.Intn(10) + 1)
	d := make([]int64, n)
	for i := 0; i < n; i++ {
		d[i] = int64(rng.Intn(3) + 1)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, x))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", d[i]))
	}
	sb.WriteByte('\n')
	expect := expected(n, x, d)
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// simple edge case
	in := "1 1\n1\n"
	exp := "2"
	out, err := runCandidate(bin, in)
	if err != nil {
		fmt.Fprintf(os.Stderr, "edge case failed: %v\n", err)
		os.Exit(1)
	}
	if out != exp {
		fmt.Fprintf(os.Stderr, "edge case failed: expected %s got %s\n", exp, out)
		os.Exit(1)
	}

	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
