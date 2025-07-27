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

const maxN = 100
const maxM = 100

var grundy [maxN + 2][maxM + 2]int

func init() {
	seen := make([]int, maxN*maxM+5)
	timestamp := 1
	for i := maxN; i >= 1; i-- {
		for j := maxM; j >= 1; j-- {
			timestamp++
			seen[0] = timestamp
			for x := i; x <= maxN; x++ {
				for y := j; y <= maxM; y++ {
					if x == i && y == j {
						continue
					}
					v := grundy[x][y]
					if v < len(seen) {
						seen[v] = timestamp
					}
				}
			}
			mex := 0
			for seen[mex] == timestamp {
				mex++
			}
			grundy[i][j] = mex
		}
	}
}

func expected(board [][]int) string {
	xor := 0
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			if board[i][j]&1 == 1 {
				xor ^= grundy[i+1][j+1]
			}
		}
	}
	if xor != 0 {
		return "Ashish"
	}
	return "Jeel"
}

func runCase(bin string, board [][]int) error {
	n := len(board)
	m := len(board[0])
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("1\n%d %d\n", n, m))
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			sb.WriteString(fmt.Sprintf("%d", board[i][j]))
			if j+1 < m {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
	}
	input := sb.String()

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
	want := expected(board)
	if strings.ToLower(got) != strings.ToLower(want) {
		return fmt.Errorf("expected %s got %s", want, got)
	}
	return nil
}

func randomBoard(rng *rand.Rand) [][]int {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	b := make([][]int, n)
	for i := range b {
		b[i] = make([]int, m)
		for j := range b[i] {
			b[i][j] = rng.Intn(3)
		}
	}
	return b
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		board := randomBoard(rng)
		if err := runCase(bin, board); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
