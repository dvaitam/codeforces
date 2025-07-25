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
	n     int
	rects [][4]int
}

func expected(t Test) int {
	total := 0
	for _, r := range t.rects {
		width := r[2] - r[0] + 1
		height := r[3] - r[1] + 1
		total += width * height
	}
	return total
}

func genCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(100) + 1
	rects := make([][4]int, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		x1 := rng.Intn(100) + 1
		x2 := rng.Intn(100) + 1
		if x1 > x2 {
			x1, x2 = x2, x1
		}
		y1 := rng.Intn(100) + 1
		y2 := rng.Intn(100) + 1
		if y1 > y2 {
			y1, y2 = y2, y1
		}
		rects[i] = [4]int{x1, y1, x2, y2}
		fmt.Fprintf(&sb, "%d %d %d %d\n", x1, y1, x2, y2)
	}
	return sb.String(), expected(Test{n, rects})
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
