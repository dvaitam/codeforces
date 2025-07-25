package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
	"unicode"
)

func solveBoard(input string) string {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	weights := map[rune]int{'Q': 9, 'R': 5, 'B': 3, 'N': 3, 'P': 1}
	white, black := 0, 0
	for _, line := range lines {
		for _, r := range line {
			if r == '.' {
				continue
			}
			w := weights[unicode.ToUpper(r)]
			if unicode.IsUpper(r) {
				white += w
			} else {
				black += w
			}
		}
	}
	if white > black {
		return "White"
	} else if black > white {
		return "Black"
	}
	return "Draw"
}

func generateCase(rng *rand.Rand) string {
	pieces := []byte{'.', 'Q', 'q', 'R', 'r', 'B', 'b', 'N', 'n', 'P', 'p', 'K', 'k'}
	var sb strings.Builder
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			sb.WriteByte(pieces[rng.Intn(len(pieces))])
		}
		if i < 7 {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func runCase(bin, input, expected string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input + "\n")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		inp := generateCase(rng)
		exp := solveBoard(inp)
		if err := runCase(bin, inp, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\n", i+1, err, inp)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
