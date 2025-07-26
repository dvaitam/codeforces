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

func run(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func expectedF(a []int) (bool, []int) {
	n := len(a)
	total := 1 << n
	for mask := 1; mask < total-1; mask++ {
		g1, g2 := 0, 0
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				g1 = gcd(g1, a[i])
			} else {
				g2 = gcd(g2, a[i])
			}
		}
		if g1 == 1 && g2 == 1 {
			res := make([]int, n)
			for i := 0; i < n; i++ {
				if mask&(1<<i) != 0 {
					res[i] = 1
				} else {
					res[i] = 2
				}
			}
			return true, res
		}
	}
	return false, nil
}

func generateCase(rng *rand.Rand) ([]byte, bool, []int) {
	n := rng.Intn(6) + 2
	a := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(50) + 1
	}
	var b bytes.Buffer
	fmt.Fprintf(&b, "%d\n", n)
	for i, v := range a {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprint(&b, v)
	}
	b.WriteByte('\n')
	ok, assignment := expectedF(a)
	return b.Bytes(), ok, assignment
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input, ok, assign := generateCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, string(input))
			os.Exit(1)
		}
		lines := strings.Split(strings.TrimSpace(out), "\n")
		if !ok {
			if strings.TrimSpace(lines[0]) != "NO" {
				fmt.Fprintf(os.Stderr, "case %d failed: expected NO got %s\ninput:\n%s", i, out, string(input))
				os.Exit(1)
			}
			continue
		}
		if strings.TrimSpace(lines[0]) != "YES" {
			fmt.Fprintf(os.Stderr, "case %d failed: expected YES got %s\ninput:\n%s", i, out, string(input))
			os.Exit(1)
		}
		if len(lines) < 2 {
			fmt.Fprintf(os.Stderr, "case %d failed: incomplete output\ninput:\n%s", i, string(input))
			os.Exit(1)
		}
		fields := strings.Fields(lines[1])
		if len(fields) != len(assign) {
			fmt.Fprintf(os.Stderr, "case %d failed: wrong output length\ninput:\n%s", i, string(input))
			os.Exit(1)
		}
		for j, f := range fields {
			if f != fmt.Sprint(assign[j]) {
				fmt.Fprintf(os.Stderr, "case %d failed: expected %v got %s\ninput:\n%s", i, assign, out, string(input))
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
