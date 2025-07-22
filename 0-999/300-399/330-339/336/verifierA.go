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

func expectedAnswerA(x, y int64) string {
	L := abs64(x) + abs64(y)
	x1 := L
	if x <= 0 {
		x1 = -L
	}
	y1 := int64(0)
	x2 := int64(0)
	y2 := L
	if y <= 0 {
		y2 = -L
	}
	if x1 >= x2 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
	}
	return fmt.Sprintf("%d %d %d %d", x1, y1, x2, y2)
}

func abs64(a int64) int64 {
	if a < 0 {
		return -a
	}
	return a
}

func generateCaseA(rng *rand.Rand) (int64, int64) {
	x := rng.Int63n(2001) - 1000
	for x == 0 {
		x = rng.Int63n(2001) - 1000
	}
	y := rng.Int63n(2001) - 1000
	for y == 0 {
		y = rng.Int63n(2001) - 1000
	}
	return x, y
}

func runCaseA(bin string, x, y int64) error {
	input := fmt.Sprintf("%d %d\n", x, y)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expected := expectedAnswerA(x, y)
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
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
		x, y := generateCaseA(rng)
		if err := runCaseA(bin, x, y); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%d %d\n", i+1, err, x, y)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
