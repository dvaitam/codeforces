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

func expected(n int64) int {
	switch {
	case n == 1:
		return 0
	case n == 2:
		return 1
	case n == 3:
		return 2
	case n%2 == 0:
		return 2
	default:
		return 3
	}
}

func runCase(bin string, n int64) error {
	input := fmt.Sprintf("1\n%d\n", n)
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
	var got int
	if _, err := fmt.Fscan(&out, &got); err != nil {
		return fmt.Errorf("failed to parse output: %v", err)
	}
	want := expected(n)
	if got != want {
		return fmt.Errorf("n=%d expected %d got %d", n, want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1000000000}
	for len(cases) < 100 {
		cases = append(cases, rng.Int63n(1000000000)+1)
	}
	for i, n := range cases {
		if err := runCase(bin, n); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
