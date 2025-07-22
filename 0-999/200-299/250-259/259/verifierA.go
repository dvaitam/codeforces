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

func solve(board []string) string {
	for i := 0; i < 8; i++ {
		var target strings.Builder
		for j := 0; j < 8; j++ {
			if (i+j)%2 == 0 {
				target.WriteByte('W')
			} else {
				target.WriteByte('B')
			}
		}
		t := target.String()
		doubled := board[i] + board[i]
		if !strings.Contains(doubled, t) {
			return "NO"
		}
	}
	return "YES"
}

func genCase(rng *rand.Rand) (string, string) {
	board := make([]string, 8)
	var sb strings.Builder
	for i := 0; i < 8; i++ {
		row := make([]byte, 8)
		for j := 0; j < 8; j++ {
			if rng.Intn(2) == 0 {
				row[j] = 'W'
			} else {
				row[j] = 'B'
			}
		}
		board[i] = string(row)
		sb.WriteString(board[i])
		sb.WriteByte('\n')
	}
	exp := solve(board)
	sbInput := sb.String()
	return sbInput, exp
}

func runCandidate(bin, input, expected string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
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
		in, exp := genCase(rng)
		if err := runCandidate(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
