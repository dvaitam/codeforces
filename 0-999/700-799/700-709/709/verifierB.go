package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func expected(xs []int, a int) int {
	n := len(xs)
	if n <= 1 {
		return 0
	}
	sort.Ints(xs)
	mn := xs[0]
	mx := xs[n-1]
	secondMn := xs[1]
	secondMx := xs[n-2]
	opt1 := (mx - secondMn) + min(abs(a-secondMn), abs(a-mx))
	opt2 := (secondMx - mn) + min(abs(a-mn), abs(a-secondMx))
	if opt2 < opt1 {
		return opt2
	}
	return opt1
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func genCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(20) + 1
	a := rng.Intn(2001) - 1000
	xs := make([]int, n)
	for i := range xs {
		xs[i] = rng.Intn(2001) - 1000
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, a)
	for i, v := range xs {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String(), expected(xs, a)
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
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		var got int
		if _, err := fmt.Fscan(strings.NewReader(out), &got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: bad output: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
