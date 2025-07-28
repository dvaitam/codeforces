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

func solveB(n int) []int {
	m := n - 1
	d := 1
	for d<<1 <= m {
		d <<= 1
	}
	res := make([]int, 0, n)
	for x := m; x >= d; x-- {
		res = append(res, x)
	}
	d--
	res = append(res, 0)
	for d > 0 {
		res = append(res, d)
		d--
	}
	return res
}

func generateCase(rng *rand.Rand) (string, []int) {
	n := rng.Intn(100) + 2
	input := fmt.Sprintf("1\n%d\n", n)
	return input, solveB(n)
}

func runCase(bin, input string, exp []int) error {
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
	parts := strings.Fields(strings.TrimSpace(out.String()))
	if len(parts) != len(exp) {
		return fmt.Errorf("expected %d numbers got %d", len(exp), len(parts))
	}
	for i, p := range parts {
		var v int
		if _, err := fmt.Sscan(p, &v); err != nil {
			return fmt.Errorf("bad int on position %d: %v", i+1, err)
		}
		if v != exp[i] {
			return fmt.Errorf("pos %d expected %d got %d", i+1, exp[i], v)
		}
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
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
