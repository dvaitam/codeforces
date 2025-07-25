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

func expectedAnswerD(board [][]byte) string {
	n := len(board)
	size := 2*n - 1
	center := n - 1
	valid := make([][]bool, size)
	for i := range valid {
		valid[i] = make([]bool, size)
		for j := range valid[i] {
			valid[i][j] = true
		}
	}
	valid[center][center] = false
	for dx := -center; dx <= center; dx++ {
		for dy := -center; dy <= center; dy++ {
			kx := dx + center
			ky := dy + center
			if kx == center && ky == center {
				continue
			}
			ok := true
			for i := 0; i < n && ok; i++ {
				for j := 0; j < n; j++ {
					if board[i][j] != 'o' {
						continue
					}
					ii := i + dx
					jj := j + dy
					if ii >= 0 && ii < n && jj >= 0 && jj < n {
						if board[ii][jj] == '.' {
							ok = false
							break
						}
					}
				}
			}
			valid[kx][ky] = ok
		}
	}
	got := make([][]byte, n)
	for i := range got {
		got[i] = make([]byte, n)
		for j := range got[i] {
			got[i][j] = '.'
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if board[i][j] == 'o' {
				got[i][j] = 'o'
				for dx := -center; dx <= center; dx++ {
					for dy := -center; dy <= center; dy++ {
						kx := dx + center
						ky := dy + center
						if !valid[kx][ky] {
							continue
						}
						ii := i + dx
						jj := j + dy
						if ii >= 0 && ii < n && jj >= 0 && jj < n {
							if got[ii][jj] == '.' {
								got[ii][jj] = 'x'
							}
						}
					}
				}
			}
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if board[i][j] == 'x' && got[i][j] != 'x' {
				return "NO"
			}
			if board[i][j] == '.' && got[i][j] == 'x' {
				return "NO"
			}
		}
	}
	var sb strings.Builder
	sb.WriteString("YES\n")
	for dx := -center; dx <= center; dx++ {
		for dy := -center; dy <= center; dy++ {
			kx := dx + center
			ky := dy + center
			if kx == center && ky == center {
				sb.WriteByte('o')
			} else if valid[kx][ky] {
				sb.WriteByte('x')
			} else {
				sb.WriteByte('.')
			}
		}
		if dx != center {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func generateCaseD(rng *rand.Rand) [][]byte {
	n := rng.Intn(4) + 1
	board := make([][]byte, n)
	hasO := false
	for i := range board {
		board[i] = make([]byte, n)
		for j := range board[i] {
			r := rng.Intn(3)
			switch r {
			case 0:
				board[i][j] = '.'
			case 1:
				board[i][j] = 'x'
			case 2:
				board[i][j] = 'o'
				hasO = true
			}
		}
	}
	if !hasO {
		board[0][0] = 'o'
	}
	return board
}

func runCaseD(bin string, board [][]byte) error {
	var sb strings.Builder
	n := len(board)
	sb.WriteString(fmt.Sprint(n, "\n"))
	for _, row := range board {
		sb.WriteString(string(row))
		sb.WriteByte('\n')
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
	expected := strings.TrimSpace(expectedAnswerD(board))
	if got != expected {
		return fmt.Errorf("expected\n%s\ngot\n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		board := generateCaseD(rng)
		if err := runCaseD(bin, board); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
