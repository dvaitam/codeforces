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

func expected(k, r int) int {
	for i := 1; i <= 10; i++ {
		total := i * k
		if total%10 == 0 || total%10 == r {
			return i
		}
	}
	return -1
}

func runCase(bin string, k, r int) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	input := fmt.Sprintf("%d %d\n", k, r)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var ans int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &ans); err != nil {
		return fmt.Errorf("bad output: %s", out.String())
	}
	if ans != expected(k, r) {
		return fmt.Errorf("expected %d got %d", expected(k, r), ans)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := []struct{ k, r int }{
		{1, 1},
		{117, 3},
		{10, 1},
		{999, 9},
		{1000, 5},
	}
	for len(cases) < 105 {
		cases = append(cases, struct{ k, r int }{
			rng.Intn(1000) + 1,
			rng.Intn(9) + 1,
		})
	}

	for i, tc := range cases {
		if err := runCase(bin, tc.k, tc.r); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: %d %d\n", i+1, err, tc.k, tc.r)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
