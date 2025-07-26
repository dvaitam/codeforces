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

func expected(arr []int) int {
	last := make(map[int]int)
	best := len(arr) + 10
	for i, v := range arr {
		if p, ok := last[v]; ok {
			if i-p+1 < best {
				best = i - p + 1
			}
		}
		last[v] = i
	}
	if best > len(arr) {
		return -1
	}
	return best
}

func generateCase(rng *rand.Rand) (string, []string) {
	t := rng.Intn(10) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	exp := make([]string, t)
	for i := 0; i < t; i++ {
		n := rng.Intn(10) + 1
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(5)
		}
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", arr[j]))
		}
		sb.WriteByte('\n')
		exp[i] = fmt.Sprintf("%d", expected(arr))
	}
	return sb.String(), exp
}

func runCase(bin, input string, exp []string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != len(exp) {
		return fmt.Errorf("expected %d lines got %d", len(exp), len(lines))
	}
	for i, line := range lines {
		if strings.TrimSpace(line) != exp[i] {
			return fmt.Errorf("line %d expected %s got %s", i+1, exp[i], strings.TrimSpace(line))
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
