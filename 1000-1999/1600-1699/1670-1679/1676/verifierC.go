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

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 2
	m := rng.Intn(8) + 1
	words := make([]string, n)
	for i := 0; i < n; i++ {
		var sb strings.Builder
		for j := 0; j < m; j++ {
			sb.WriteByte(byte('a' + rng.Intn(26)))
		}
		words[i] = sb.String()
	}
	minDiff := 1<<31 - 1
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			diff := 0
			for k := 0; k < m; k++ {
				a := words[i][k]
				b := words[j][k]
				if a > b {
					diff += int(a - b)
				} else {
					diff += int(b - a)
				}
			}
			if diff < minDiff {
				minDiff = diff
			}
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n", n, m)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%s\n", words[i])
	}
	return sb.String(), fmt.Sprintf("%d", minDiff)
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
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
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
