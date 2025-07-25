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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	k := rng.Intn(4) + 1
	per := n * m * k
	a := rng.Intn(per) + 1
	b := rng.Intn(per-1) + 1
	if b >= a {
		b++
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	sb.WriteString(fmt.Sprintf("%d %d\n", a, b))

	perEntrance := m * k
	entranceA := (a-1)/perEntrance + 1
	entranceB := (b-1)/perEntrance + 1
	floorA := ((a-1)%perEntrance)/k + 1
	floorB := ((b-1)%perEntrance)/k + 1
	downA := min((floorA-1)*5, 10+(floorA-1))
	upB := min((floorB-1)*5, 10+(floorB-1))
	diff := abs(entranceA - entranceB)
	walk := min(diff, n-diff) * 15
	exp := fmt.Sprintf("%d\n", downA+walk+upB)
	return sb.String(), exp
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
