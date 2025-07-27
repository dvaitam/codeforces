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

func runCandidate(bin, input string) (string, error) {
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

func solve(input string) string {
	in := strings.NewReader(input)
	var t int
	fmt.Fscan(in, &t)
	var out strings.Builder
	for i := 0; i < t; i++ {
		var n, d int
		fmt.Fscan(in, &n, &d)
		a := make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &a[j])
		}
		res := a[0]
		days := d
		for j := 1; j < n && days > 0; j++ {
			cost := j
			move := days / cost
			if move > a[j] {
				move = a[j]
			}
			res += move
			days -= move * cost
		}
		fmt.Fprintf(&out, "%d", res)
		if i+1 < t {
			out.WriteByte('\n')
		}
	}
	return out.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	t := rng.Intn(3) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rng.Intn(5) + 1
		d := rng.Intn(10)
		fmt.Fprintf(&sb, "%d %d\n", n, d)
		for j := 0; j < n; j++ {
			val := rng.Intn(10)
			if j+1 == n {
				fmt.Fprintf(&sb, "%d\n", val)
			} else {
				fmt.Fprintf(&sb, "%d ", val)
			}
		}
	}
	input := sb.String()
	return input, solve(input)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		got, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
