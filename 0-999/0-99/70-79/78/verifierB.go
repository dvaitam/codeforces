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

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(94) + 7 // 7..100
	input := fmt.Sprintf("%d\n", n)
	return input, n
}

func validate(output string, n int) error {
	out := strings.TrimSpace(output)
	if len(out) != n {
		return fmt.Errorf("length mismatch: expected %d got %d", n, len(out))
	}
	colors := "ROYGBIV"
	freq := make(map[rune]int)
	for _, ch := range out {
		if !strings.ContainsRune(colors, ch) {
			return fmt.Errorf("invalid character %q", ch)
		}
		freq[ch]++
	}
	for _, ch := range colors {
		if freq[ch] == 0 {
			return fmt.Errorf("color %c missing", ch)
		}
	}
	for i := 0; i < n; i++ {
		set := make(map[byte]struct{})
		for j := 0; j < 4; j++ {
			set[out[(i+j)%n]] = struct{}{}
		}
		if len(set) != 4 {
			return fmt.Errorf("repeated colors in segment starting %d", i)
		}
	}
	return nil
}

func runCase(bin string, input string, n int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if err := validate(out.String(), n); err != nil {
		return err
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, n := generateCase(rng)
		if err := runCase(bin, in, n); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
