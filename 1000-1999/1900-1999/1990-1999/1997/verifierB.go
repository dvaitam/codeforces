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

func expectedAnswerB(n int, a, b string) int {
	ans := 0
	for i := 0; i < n-2; i++ {
		if a[i] != b[i] && a[i] == a[i+2] && b[i] == b[i+2] && a[i+1] == '.' && b[i+1] == '.' {
			ans++
		}
	}
	return ans
}

func generateCaseB(rng *rand.Rand) (int, string, string) {
	n := rng.Intn(8) + 3 // at least 3
	row1 := make([]byte, n)
	row2 := make([]byte, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			row1[i] = '.'
		} else {
			row1[i] = 'x'
		}
		if rng.Intn(2) == 0 {
			row2[i] = '.'
		} else {
			row2[i] = 'x'
		}
	}
	return n, string(row1), string(row2)
}

func runCaseB(bin string, n int, a, b string) error {
	input := fmt.Sprintf("1\n%d\n%s\n%s\n", n, a, b)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expected := fmt.Sprint(expectedAnswerB(n, a, b))
	if got != expected {
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
		n, r1, r2 := generateCaseB(rng)
		if err := runCaseB(bin, n, r1, r2); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%d\n%s\n%s\n", i+1, err, n, r1, r2)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
