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

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func expected(x1, y1, x2, y2 int64) string {
	dx := abs64(x1 - x2)
	dy := abs64(y1 - y2)
	if dx > dy {
		return fmt.Sprintf("%d", dx)
	}
	return fmt.Sprintf("%d", dy)
}

func runCase(bin, input, want string) error {
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
	want = strings.TrimSpace(want)
	if got != want {
		return fmt.Errorf("expected %s got %s", want, got)
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

	tests := []struct{ x1, y1, x2, y2 int64 }{
		{0, 0, 0, 0},
		{0, 0, 1, 1},
	}
	for i := 0; i < 100; i++ {
		x1 := rng.Int63n(2_000_000_001) - 1_000_000_000
		y1 := rng.Int63n(2_000_000_001) - 1_000_000_000
		x2 := rng.Int63n(2_000_000_001) - 1_000_000_000
		y2 := rng.Int63n(2_000_000_001) - 1_000_000_000
		tests = append(tests, struct{ x1, y1, x2, y2 int64 }{x1, y1, x2, y2})
	}

	for i, tc := range tests {
		input := fmt.Sprintf("%d %d\n%d %d\n", tc.x1, tc.y1, tc.x2, tc.y2)
		want := expected(tc.x1, tc.y1, tc.x2, tc.y2)
		if err := runCase(bin, input, want); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
