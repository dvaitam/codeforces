package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Board [4]string

func run(bin string, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(k int, board Board) string {
	counts := make([]int, 10)
	for i := 0; i < 4; i++ {
		for _, c := range board[i] {
			if c >= '1' && c <= '9' {
				counts[c-'0']++
			}
		}
	}
	limit := 2 * k
	for d := 1; d <= 9; d++ {
		if counts[d] > limit {
			return "NO"
		}
	}
	return "YES"
}

func runCase(bin string, k int, board Board) error {
	input := fmt.Sprintf("%d\n%s\n%s\n%s\n%s\n", k, board[0], board[1], board[2], board[3])
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	exp := expected(k, board)
	if out != exp {
		return fmt.Errorf("expected %s got %s", exp, out)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(42)
	digits := []rune(".123456789")
	total := 100
	for i := 0; i < total; i++ {
		k := rand.Intn(5) + 1
		var b Board
		for r := 0; r < 4; r++ {
			row := make([]rune, 4)
			for c := 0; c < 4; c++ {
				row[c] = digits[rand.Intn(len(digits))]
			}
			b[r] = string(row)
		}
		if err := runCase(bin, k, b); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			fmt.Fprintf(os.Stderr, "k=%d board=\n%s\n%s\n%s\n%s\n", k, b[0], b[1], b[2], b[3])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", total)
}
