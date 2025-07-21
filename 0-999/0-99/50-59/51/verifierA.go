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

func canonical(a, b, c, d int) string {
	r1 := fmt.Sprintf("%d%d%d%d", a, b, c, d)
	r2 := fmt.Sprintf("%d%d%d%d", c, a, d, b)
	r3 := fmt.Sprintf("%d%d%d%d", d, c, b, a)
	r4 := fmt.Sprintf("%d%d%d%d", b, d, a, c)
	min := r1
	if r2 < min {
		min = r2
	}
	if r3 < min {
		min = r3
	}
	if r4 < min {
		min = r4
	}
	return min
}

func solve(amulets [][4]int) string {
	seen := make(map[string]struct{})
	for _, q := range amulets {
		s := canonical(q[0], q[1], q[2], q[3])
		seen[s] = struct{}{}
	}
	return fmt.Sprintf("%d", len(seen))
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	am := make([][4]int, n)
	for i := 0; i < n; i++ {
		a := rng.Intn(6) + 1
		b := rng.Intn(6) + 1
		c := rng.Intn(6) + 1
		d := rng.Intn(6) + 1
		am[i] = [4]int{a, b, c, d}
		sb.WriteString(fmt.Sprintf("%d %d\n%d %d\n", a, b, c, d))
		if i != n-1 {
			sb.WriteString("**\n")
		}
	}
	input := sb.String()
	expected := solve(am)
	return input, expected
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

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := generateCase(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
