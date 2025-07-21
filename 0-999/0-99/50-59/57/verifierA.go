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

func perimeterPos(n, x, y int) int {
	switch {
	case y == 0:
		return x
	case x == n:
		return n + y
	case y == n:
		return 2*n + (n - x)
	default:
		return 3*n + (n - y)
	}
}

func expected(n, x1, y1, x2, y2 int) int {
	t1 := perimeterPos(n, x1, y1)
	t2 := perimeterPos(n, x2, y2)
	d := t1 - t2
	if d < 0 {
		d = -d
	}
	per := 4 * n
	if d > per-d {
		d = per - d
	}
	return d
}

func genCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(1000) + 1
	choose := func() (int, int) {
		side := rng.Intn(4)
		v := rng.Intn(n + 1)
		switch side {
		case 0:
			return v, 0
		case 1:
			return n, v
		case 2:
			return n - v, n
		default:
			return 0, n - v
		}
	}
	x1, y1 := choose()
	x2, y2 := choose()
	input := fmt.Sprintf("%d %d %d %d %d\n", n, x1, y1, x2, y2)
	return input, expected(n, x1, y1, x2, y2)
}

func runCase(bin, input string, exp int) error {
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
