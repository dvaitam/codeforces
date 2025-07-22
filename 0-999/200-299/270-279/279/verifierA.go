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

func solveA(x, y int) string {
	if x == 0 && y == 0 {
		return "0"
	}
	dirs := [4][2]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}
	cx, cy := 0, 0
	d := 0
	turns := 0
	first := true
	for k := 1; ; k++ {
		for rep := 0; rep < 2; rep++ {
			if first {
				first = false
			} else {
				d = (d + 1) % 4
				turns++
			}
			dx, dy := dirs[d][0], dirs[d][1]
			for step := 0; step < k; step++ {
				cx += dx
				cy += dy
				if cx == x && cy == y {
					return fmt.Sprintf("%d", turns)
				}
			}
		}
	}
}

func generateCaseA(rng *rand.Rand) (string, string) {
	x := rng.Intn(201) - 100
	y := rng.Intn(201) - 100
	input := fmt.Sprintf("%d %d\n", x, y)
	expected := solveA(x, y)
	return input, expected
}

func runCase(bin, input, expected string) error {
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
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", expected, got)
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
	for i := 0; i < 100; i++ {
		in, exp := generateCaseA(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
