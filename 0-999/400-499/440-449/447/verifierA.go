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

func expectedAnswerA(p int, xs []int) int {
	used := make([]bool, p)
	for i, x := range xs {
		h := x % p
		if used[h] {
			return i + 1
		}
		used[h] = true
	}
	return -1
}

func generateCaseA(rng *rand.Rand) (int, []int) {
	p := rng.Intn(299) + 2 // 2..300
	n := rng.Intn(299) + 2 // 2..300
	xs := make([]int, n)
	for i := range xs {
		xs[i] = rng.Intn(1_000_000_000)
	}
	return p, xs
}

func runCaseA(bin string, p int, xs []int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", p, len(xs)))
	for _, x := range xs {
		sb.WriteString(fmt.Sprintf("%d\n", x))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expected := fmt.Sprint(expectedAnswerA(p, xs))
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
		p, xs := generateCaseA(rng)
		if err := runCaseA(bin, p, xs); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%d %d\n%v\n", i+1, err, p, len(xs), xs)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
