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

func expectedAnswer(magnets []string) string {
	if len(magnets) == 0 {
		return "0"
	}
	count := 1
	prev := magnets[0]
	for _, m := range magnets[1:] {
		if m != prev {
			count++
		}
		prev = m
	}
	return fmt.Sprint(count)
}

func generateCase(rng *rand.Rand) []string {
	n := rng.Intn(100) + 1 // 1..100 magnets
	mags := make([]string, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			mags[i] = "01"
		} else {
			mags[i] = "10"
		}
	}
	return mags
}

func runCase(bin string, mags []string) error {
	input := fmt.Sprintf("%d\n%s\n", len(mags), strings.Join(mags, "\n"))
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expected := expectedAnswer(mags)
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
		mags := generateCase(rng)
		if err := runCase(bin, mags); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%d\n%s\n", i+1, err, len(mags), strings.Join(mags, "\n"))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
