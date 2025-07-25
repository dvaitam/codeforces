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

func expected(board [4]string) string {
	check := func(a, b, c byte) bool {
		x, dot := 0, 0
		if a == 'x' {
			x++
		} else if a == '.' {
			dot++
		}
		if b == 'x' {
			x++
		} else if b == '.' {
			dot++
		}
		if c == 'x' {
			x++
		} else if c == '.' {
			dot++
		}
		return x == 2 && dot == 1
	}
	yes := false
	for i := 0; i < 4 && !yes; i++ {
		for j := 0; j <= 1; j++ {
			if check(board[i][j], board[i][j+1], board[i][j+2]) {
				yes = true
				break
			}
		}
	}
	for j := 0; j < 4 && !yes; j++ {
		for i := 0; i <= 1; i++ {
			if check(board[i][j], board[i+1][j], board[i+2][j]) {
				yes = true
				break
			}
		}
	}
	for i := 0; i <= 1 && !yes; i++ {
		for j := 0; j <= 1; j++ {
			if check(board[i][j], board[i+1][j+1], board[i+2][j+2]) {
				yes = true
				break
			}
		}
	}
	for i := 0; i <= 1 && !yes; i++ {
		for j := 2; j < 4; j++ {
			if check(board[i][j], board[i+1][j-1], board[i+2][j-2]) {
				yes = true
				break
			}
		}
	}
	if yes {
		return "YES"
	}
	return "NO"
}

func runCase(bin string, board [4]string) error {
	var input strings.Builder
	for i := 0; i < 4; i++ {
		input.WriteString(board[i])
		if i+1 < 4 {
			input.WriteByte('\n')
		}
	}
	input.WriteByte('\n')
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	want := expected(board)
	if got != want {
		return fmt.Errorf("expected %s got %s", want, got)
	}
	return nil
}

func randBoard(rng *rand.Rand) [4]string {
	var b [4]string
	for i := 0; i < 4; i++ {
		row := make([]byte, 4)
		for j := 0; j < 4; j++ {
			switch rng.Intn(3) {
			case 0:
				row[j] = '.'
			case 1:
				row[j] = 'x'
			default:
				row[j] = 'o'
			}
		}
		b[i] = string(row)
	}
	return b
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := [][4]string{
		{"....", "....", "....", "...."},
		{"xxx.", "....", "....", "...."},
		{"xx..", "x...", "....", "...."},
		{"..x.", "..x.", "..x.", "...."},
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, randBoard(rng))
	}
	for i, c := range cases {
		if err := runCase(bin, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\nboard: %v\n", i+1, err, c)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
