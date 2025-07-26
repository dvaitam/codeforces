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

func expected(n, l, r int, a, b []int) string {
	for i := 0; i < l-1; i++ {
		if a[i] != b[i] {
			return "LIE"
		}
	}
	for i := r; i < n; i++ {
		if a[i] != b[i] {
			return "LIE"
		}
	}
	freq := make([]int, n+1)
	for i := l - 1; i <= r-1; i++ {
		freq[a[i]]++
		freq[b[i]]--
	}
	for i := 1; i <= n; i++ {
		if freq[i] != 0 {
			return "LIE"
		}
	}
	return "TRUTH"
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

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 1
	if rng.Float64() < 0.1 {
		n = rng.Intn(20) + 1
	}
	l := rng.Intn(n) + 1
	r := rng.Intn(n-l+1) + l
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(n) + 1
		b[i] = rng.Intn(n) + 1
	}
	var bld strings.Builder
	fmt.Fprintf(&bld, "%d %d %d\n", n, l, r)
	for i := 0; i < n; i++ {
		if i > 0 {
			bld.WriteByte(' ')
		}
		fmt.Fprintf(&bld, "%d", a[i])
	}
	bld.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			bld.WriteByte(' ')
		}
		fmt.Fprintf(&bld, "%d", b[i])
	}
	bld.WriteByte('\n')
	exp := expected(n, l, r, a, b)
	return bld.String(), exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, strings.TrimSpace(out), input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
