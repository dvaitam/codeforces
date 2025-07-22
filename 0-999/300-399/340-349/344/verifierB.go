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

func expectedAnswerB(a, b, c int64) string {
	sumABc := a + b - c
	sumBCa := b + c - a
	sumACb := a + c - b
	if sumABc < 0 || sumBCa < 0 || sumACb < 0 || sumABc%2 != 0 || sumBCa%2 != 0 || sumACb%2 != 0 {
		return "Impossible"
	}
	x := sumABc / 2
	y := sumBCa / 2
	z := sumACb / 2
	if x < 0 || y < 0 || z < 0 {
		return "Impossible"
	}
	return fmt.Sprintf("%d %d %d", x, y, z)
}

func generateCaseB(rng *rand.Rand) (int64, int64, int64) {
	a := rng.Int63n(1000000) + 1
	b := rng.Int63n(1000000) + 1
	c := rng.Int63n(1000000) + 1
	return a, b, c
}

func runCase(bin string, a, b, c int64) error {
	input := fmt.Sprintf("%d %d %d\n", a, b, c)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expected := expectedAnswerB(a, b, c)
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
		a, b, c := generateCaseB(rng)
		if err := runCase(bin, a, b, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%d %d %d\n", i+1, err, a, b, c)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
