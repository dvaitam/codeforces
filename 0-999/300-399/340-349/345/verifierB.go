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

func run(bin string, input string) (string, error) {
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

func countBadBases(n int) int {
	if strings.Contains(strconv.Itoa(n), "13") {
		return -1
	}
	cnt := 0
	for b := 2; b <= n; b++ {
		x := n
		var digits []int
		for x > 0 {
			digits = append(digits, x%b)
			x /= b
		}
		var sb strings.Builder
		for i := len(digits) - 1; i >= 0; i-- {
			sb.WriteString(strconv.Itoa(digits[i]))
		}
		if strings.Contains(sb.String(), "13") {
			cnt++
		}
	}
	return cnt
}

func runCase(bin string, n int) error {
	input := fmt.Sprintf("%d\n", n)
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	got, err := strconv.Atoi(strings.TrimSpace(out))
	if err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	exp := countBadBases(n)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(1000) + 1
		if err := runCase(bin, n); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v (n=%d)\n", i+1, err, n)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
