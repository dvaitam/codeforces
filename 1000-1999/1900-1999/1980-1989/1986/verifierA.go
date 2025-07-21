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

type Test struct {
	x1, x2, x3 int
}

func expected(t Test) int {
	minv := t.x1
	if t.x2 < minv {
		minv = t.x2
	}
	if t.x3 < minv {
		minv = t.x3
	}
	maxv := t.x1
	if t.x2 > maxv {
		maxv = t.x2
	}
	if t.x3 > maxv {
		maxv = t.x3
	}
	return maxv - minv
}

func genCase(rng *rand.Rand) (string, int) {
	c := Test{rng.Intn(10) + 1, rng.Intn(10) + 1, rng.Intn(10) + 1}
	input := fmt.Sprintf("1\n%d %d %d\n", c.x1, c.x2, c.x3)
	return input, expected(c)
}

func runCase(bin, input string, exp int) error {
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
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
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
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
