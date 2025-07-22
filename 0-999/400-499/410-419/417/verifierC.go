package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func allowedK(n int) int {
	if n%2 == 0 {
		return n/2 - 1
	}
	return n / 2
}

func runCase(bin string, n, k int) error {
	input := fmt.Sprintf("%d %d\n", n, k)
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

	rdr := bufio.NewReader(strings.NewReader(out.String()))
	var first int
	if _, err := fmt.Fscan(rdr, &first); err != nil {
		return fmt.Errorf("bad output")
	}
	if first == -1 {
		if k <= allowedK(n) {
			return fmt.Errorf("unexpected -1 output")
		}
		return nil
	}
	total := first
	if total != n*k {
		return fmt.Errorf("expected %d pairs, got %d", n*k, total)
	}
	wins := make([]int, n+1)
	seen := make(map[[2]int]bool)
	for i := 0; i < total; i++ {
		var a, b int
		if _, err := fmt.Fscan(rdr, &a, &b); err != nil {
			return fmt.Errorf("not enough pairs")
		}
		if a < 1 || a > n || b < 1 || b > n || a == b {
			return fmt.Errorf("invalid pair")
		}
		key := [2]int{a, b}
		if seen[key] {
			return fmt.Errorf("duplicate pair")
		}
		seen[key] = true
		wins[a]++
	}
	for i := 1; i <= n; i++ {
		if wins[i] != k {
			return fmt.Errorf("team %d wins %d times", i, wins[i])
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
		n := rng.Intn(10) + 2 // at least 2
		k := rng.Intn(n)
		if err := runCase(bin, n, k); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: n=%d k=%d\n", i+1, err, n, k)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
