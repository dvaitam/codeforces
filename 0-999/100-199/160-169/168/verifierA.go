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

func expected(n, x, y int64) int64 {
	required := (n*y + 99) / 100
	if required < x {
		return 0
	}
	return required - x
}

func runCase(bin string, n, x, y int64) error {
	input := fmt.Sprintf("%d %d %d\n", n, x, y)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Fscan(&out, &got); err != nil {
		return fmt.Errorf("failed to parse output: %v", err)
	}
	want := expected(n, x, y)
	if got != want {
		return fmt.Errorf("expected %d got %d", want, got)
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
	type tc struct{ n, x, y int64 }
	cases := []tc{
		{10, 1, 14},
		{10, 10, 100},
		{1, 1, 10000},
		{5, 1, 120},
		{9999, 5000, 50},
		{10000, 1, 1},
		{10000, 10000, 10000},
		{2, 1, 51},
		{7, 3, 0},
		{3, 2, 200},
	}
	for i := 0; i < 100; i++ {
		n := rng.Int63n(10000) + 1
		x := rng.Int63n(n) + 1
		y := rng.Int63n(10000) + 1
		cases = append(cases, tc{n, x, y})
	}
	for i, c := range cases {
		if err := runCase(bin, c.n, c.x, c.y); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
