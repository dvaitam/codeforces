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

type Test struct {
	n   int
	b   int
	d   int
	arr []int
}

func expected(t Test) int {
	waste := 0
	empties := 0
	for _, a := range t.arr {
		if a > t.b {
			continue
		}
		waste += a
		if waste > t.d {
			empties++
			waste = 0
		}
	}
	return empties
}

func genCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(20) + 1
	b := rng.Intn(1000) + 1
	d := b + rng.Intn(1000)
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(1000) + 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, b, d)
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	t := Test{n, b, d, arr}
	return sb.String(), expected(t)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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
