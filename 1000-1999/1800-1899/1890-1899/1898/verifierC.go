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

func expectedC(N, M, K int) string {
	rows := 2*N - 1
	cols := 2*M - 1
	C := make([][]byte, rows)
	for i := 0; i < rows; i++ {
		C[i] = make([]byte, cols)
		for j := 0; j < cols; j++ {
			C[i][j] = '#'
		}
	}
	minMoves := N + M - 2
	if K < minMoves || (K%2) != (minMoves%2) {
		return "NO"
	}
	maxRow := 2*N - 2
	maxCol := 2*M - 2
	now := 0
	if K%4 == minMoves%4 {
		now = 0
		for j := 1; j < maxCol; j += 2 {
			if now == 0 {
				C[0][j] = 'R'
			} else {
				C[0][j] = 'B'
			}
			now ^= 1
		}
		for i := 1; i < maxRow; i += 2 {
			if now == 0 {
				C[i][maxCol] = 'R'
			} else {
				C[i][maxCol] = 'B'
			}
			now ^= 1
		}
		if now == 0 {
			C[maxRow-1][maxCol] = 'B'
			C[maxRow-1][maxCol-2] = 'B'
			C[maxRow-2][maxCol-1] = 'R'
			C[maxRow][maxCol-1] = 'R'
		} else {
			C[maxRow-1][maxCol] = 'R'
			C[maxRow-1][maxCol-2] = 'R'
			C[maxRow-2][maxCol-1] = 'B'
			C[maxRow][maxCol-1] = 'B'
		}
	} else {
		now = 0
		for j := 1; j < maxCol; j += 2 {
			if now == 0 {
				C[0][j] = 'R'
			} else {
				C[0][j] = 'B'
			}
			now ^= 1
		}
		for i := 1; i < maxRow; i += 2 {
			if now == 0 {
				C[i][maxCol] = 'R'
			} else {
				C[i][maxCol] = 'B'
			}
			now ^= 1
		}
		if now == 1 {
			C[maxRow-1][maxCol] = 'B'
			C[maxRow-1][maxCol-2] = 'B'
			C[maxRow-2][maxCol-1] = 'R'
			C[maxRow][maxCol-1] = 'R'
		} else {
			C[maxRow-1][maxCol] = 'R'
			C[maxRow-1][maxCol-2] = 'R'
			C[maxRow-2][maxCol-1] = 'B'
			C[maxRow][maxCol-1] = 'B'
		}
	}
	var sb strings.Builder
	sb.WriteString("YES\n")
	for i := 0; i <= maxRow; i += 2 {
		for j := 1; j <= maxCol; j += 2 {
			if C[i][j] == 'B' {
				sb.WriteString("B ")
			} else {
				sb.WriteString("R ")
			}
		}
		sb.WriteByte('\n')
	}
	for i := 1; i <= maxRow; i += 2 {
		for j := 0; j <= maxCol; j += 2 {
			if C[i][j] == 'B' {
				sb.WriteString("B ")
			} else {
				sb.WriteString("R ")
			}
		}
		sb.WriteByte('\n')
	}
	return strings.TrimRight(sb.String(), "\n")
}

func generateCase(rng *rand.Rand) (string, string) {
	N := rng.Intn(14) + 3 // 3..16
	M := rng.Intn(14) + 3
	K := rng.Intn(50)
	input := fmt.Sprintf("1\n%d %d %d\n", N, M, K)
	expect := expectedC(N, M, K)
	return input, expect
}

func runCase(bin, input, exp string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != exp {
		return fmt.Errorf("expected:\n%s\n---got:\n%s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
