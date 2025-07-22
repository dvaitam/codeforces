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

func isLucky(x int) bool {
	if x == 0 {
		return true
	}
	for x > 0 {
		d := x % 10
		if d != 4 && d != 7 {
			return false
		}
		x /= 10
	}
	return true
}

func solveA(n int) string {
	count4 := 0
	for n%7 != 0 && n >= 4 {
		n -= 4
		count4++
	}
	if n%7 != 0 {
		return "-1"
	}
	return strings.Repeat("4", count4) + strings.Repeat("7", n/7)
}

func runCase(bin string, n int) error {
	input := fmt.Sprintf("%d\n", n)
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
	got := strings.TrimSpace(out.String())
	exp := solveA(n)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []int{1, 2, 3, 4, 5, 6, 7, 8, 11, 44, 47, 51, 1000000}
	for i := 0; i < 100; i++ {
		cases = append(cases, rng.Intn(1000000)+1)
	}
	for i, n := range cases {
		if err := runCase(bin, n); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: %d\n", i+1, err, n)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
