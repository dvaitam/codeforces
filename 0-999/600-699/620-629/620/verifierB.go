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

var seg = []int{6, 2, 5, 5, 4, 5, 6, 3, 7, 6}

func countSegments(x int) int {
	if x == 0 {
		return seg[0]
	}
	total := 0
	for x > 0 {
		total += seg[x%10]
		x /= 10
	}
	return total
}

func expected(a, b int) string {
	total := 0
	for i := a; i <= b; i++ {
		total += countSegments(i)
	}
	return fmt.Sprintf("%d", total)
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
	if got != strings.TrimSpace(want) {
		return fmt.Errorf("expected %s got %s", want, got)
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

	tests := []struct{ a, b int }{
		{1, 3},
		{10, 10},
	}
	for i := 0; i < 100; i++ {
		a := rng.Intn(1_000_000) + 1
		b := a + rng.Intn(20)
		if b > 1_000_000 {
			b = 1_000_000
		}
		tests = append(tests, struct{ a, b int }{a, b})
	}

	for i, tc := range tests {
		input := fmt.Sprintf("%d %d\n", tc.a, tc.b)
		want := expected(tc.a, tc.b)
		if err := runCase(bin, input, want); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
