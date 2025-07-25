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

func solve(a []int) int {
	n := len(a)
	pos1, posn := 0, 0
	for i, v := range a {
		if v == 1 {
			pos1 = i
		}
		if v == n {
			posn = i
		}
	}
	dist := func(x, y int) int {
		if x > y {
			return x - y
		}
		return y - x
	}
	ans := dist(pos1, 0)
	if v := dist(pos1, n-1); v > ans {
		ans = v
	}
	if v := dist(posn, 0); v > ans {
		ans = v
	}
	if v := dist(posn, n-1); v > ans {
		ans = v
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(99) + 2 // [2,100]
	perm := rng.Perm(n)
	for i := range perm {
		perm[i]++
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i, v := range perm {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	expected := fmt.Sprintf("%d\n", solve(perm))
	return sb.String(), expected
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
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
		input, expect := generateCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, strings.TrimSpace(expect), strings.TrimSpace(out), input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
