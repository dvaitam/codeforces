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

func isKnight(x1, y1, x2, y2 int) bool {
	dx := x1 - x2
	if dx < 0 {
		dx = -dx
	}
	dy := y1 - y2
	if dy < 0 {
		dy = -dy
	}
	return (dx == 1 && dy == 2) || (dx == 2 && dy == 1)
}

func posToStr(x, y int) string {
	return fmt.Sprintf("%c%d", 'a'+x-1, y)
}

func solveCase(rX, rY, kX, kY int) int {
	count := 0
	for x := 1; x <= 8; x++ {
		for y := 1; y <= 8; y++ {
			if x == rX && y == rY {
				continue
			}
			if x == kX && y == kY {
				continue
			}
			if x == rX || y == rY {
				continue
			}
			if isKnight(rX, rY, x, y) {
				continue
			}
			if isKnight(kX, kY, x, y) {
				continue
			}
			count++
		}
	}
	return count
}

func generateCase(rng *rand.Rand) (string, string) {
	for {
		rX := rng.Intn(8) + 1
		rY := rng.Intn(8) + 1
		kX := rng.Intn(8) + 1
		kY := rng.Intn(8) + 1
		if rX == kX && rY == kY {
			continue
		}
		if rX == kX || rY == kY {
			continue
		}
		if isKnight(rX, rY, kX, kY) {
			continue
		}
		input := posToStr(rX, rY) + "\n" + posToStr(kX, kY) + "\n"
		exp := fmt.Sprintf("%d", solveCase(rX, rY, kX, kY))
		return input, exp
	}
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
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
