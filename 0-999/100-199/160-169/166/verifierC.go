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

func expected(n, x int, arr []int) int {
	L, E := 0, 0
	for _, v := range arr {
		if v < x {
			L++
		} else if v == x {
			E++
		}
	}
	k := 0
	for {
		N := n + k
		pos := (N + 1) / 2
		if pos > L && pos <= L+E+k {
			return k
		}
		k++
	}
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(10) + 1
	x := rng.Intn(100) + 1
	arr := make([]int, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, x)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(100) + 1
		if i > 0 {
			fmt.Fprint(&sb, " ")
		}
		fmt.Fprint(&sb, arr[i])
	}
	fmt.Fprintln(&sb)
	return sb.String(), expected(n, x, arr)
}

func runCase(exe, input string, expectedAns int) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expectedAns {
		return fmt.Errorf("expected %d got %d", expectedAns, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
