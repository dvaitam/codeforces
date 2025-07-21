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

func solve(n, b int, a []int) int {
	ans := b
	for i := 0; i < n; i++ {
		if a[i] > b {
			continue
		}
		k := b / a[i]
		for j := i + 1; j < n; j++ {
			if a[j] <= a[i] {
				continue
			}
			val := b + k*(a[j]-a[i])
			if val > ans {
				ans = val
			}
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(8) + 1
	b := rng.Intn(1000) + 1
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(1000) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, b))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d ", a[i]))
	}
	sb.WriteString("\n")
	return sb.String(), solve(n, b, a)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		got, err2 := strconv.Atoi(strings.TrimSpace(out))
		if err2 != nil || got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
