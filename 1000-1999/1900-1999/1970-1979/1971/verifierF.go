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

func countPoints(r int) int {
	limit := r + 1
	count := 0
	r2 := r * r
	r2n := (r + 1) * (r + 1)
	for x := -limit; x <= limit; x++ {
		for y := -limit; y <= limit; y++ {
			d2 := x*x + y*y
			if d2 >= r2 && d2 < r2n {
				count++
			}
		}
	}
	return count
}

type caseF struct{ r int }

func genCase(rng *rand.Rand) caseF {
	return caseF{rng.Intn(50) + 1}
}

func runCase(bin string, tc caseF) error {
	input := fmt.Sprintf("1\n%d\n", tc.r)
	cmd := exec.Command(bin)
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
	exp := countPoints(tc.r)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
