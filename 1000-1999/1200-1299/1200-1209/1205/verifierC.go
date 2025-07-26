package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Board [][]int

func genBoard(rng *rand.Rand) Board {
	n := rng.Intn(3)*2 + 3 // 3,5,7
	b := make(Board, n)
	for i := range b {
		b[i] = make([]int, n)
		for j := range b[i] {
			if i == 0 && j == 0 {
				b[i][j] = 1
			} else if i == n-1 && j == n-1 {
				b[i][j] = 0
			} else {
				b[i][j] = rng.Intn(2)
			}
		}
	}
	return b
}

func computeDP(board Board) [][][][]bool {
	n := len(board)
	dp := make([][][][]bool, n)
	for i := range dp {
		dp[i] = make([][][]bool, n)
		for j := range dp[i] {
			dp[i][j] = make([][]bool, n)
			for a := range dp[i][j] {
				dp[i][j][a] = make([]bool, n)
			}
		}
	}
	for d := 0; d < 2*n-1; d++ {
		for x1 := 0; x1 < n; x1++ {
			for y1 := 0; y1 < n; y1++ {
				s := x1 + y1 + d
				if s >= 2*n-1 {
					continue
				}
				for x2 := 0; x2 < n; x2++ {
					y2 := s - x2
					if y2 < 0 || y2 >= n {
						continue
					}
					if x2 < x1 || y2 < y1 {
						continue
					}
					if board[x1][y1] != board[x2][y2] {
						dp[x1][y1][x2][y2] = false
						continue
					}
					if d <= 1 {
						dp[x1][y1][x2][y2] = true
						continue
					}
					v := false
					if x1+1 < n && x2-1 >= 0 && dp[x1+1][y1][x2-1][y2] {
						v = true
					}
					if x1+1 < n && y2-1 >= 0 && dp[x1+1][y1][x2][y2-1] {
						v = true
					}
					if y1+1 < n && x2-1 >= 0 && dp[x1][y1+1][x2-1][y2] {
						v = true
					}
					if y1+1 < n && y2-1 >= 0 && dp[x1][y1+1][x2][y2-1] {
						v = true
					}
					dp[x1][y1][x2][y2] = v
				}
			}
		}
	}
	return dp
}

func runCase(bin string, board Board) error {
	n := len(board)
	dp := computeDP(board)
	cmd := exec.Command(bin)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}

	writer := bufio.NewWriter(stdin)
	reader := bufio.NewReader(stdout)

	fmt.Fprintf(writer, "%d\n", n)
	writer.Flush()

	queries := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			cmd.Process.Kill()
			return fmt.Errorf("failed to read from program: %v", err)
		}
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "?") {
			parts := strings.Fields(line)
			if len(parts) != 5 {
				cmd.Process.Kill()
				return fmt.Errorf("bad query format: %s", line)
			}
			var x1, y1, x2, y2 int
			fmt.Sscan(parts[1], &x1)
			fmt.Sscan(parts[2], &y1)
			fmt.Sscan(parts[3], &x2)
			fmt.Sscan(parts[4], &y2)
			if x1 < 1 || x1 > n || y1 < 1 || y1 > n || x2 < 1 || x2 > n || y2 < 1 || y2 > n || x1 > x2 || y1 > y2 || x1+y1+2 > x2+y2 {
				cmd.Process.Kill()
				return fmt.Errorf("invalid query coordinates: %s", line)
			}
			ans := 0
			if dp[x1-1][y1-1][x2-1][y2-1] {
				ans = 1
			}
			fmt.Fprintf(writer, "%d\n", ans)
			writer.Flush()
			queries++
			if queries > n*n {
				cmd.Process.Kill()
				return fmt.Errorf("too many queries")
			}
		} else if strings.HasPrefix(line, "!") {
			final := make(Board, n)
			for i := 0; i < n; i++ {
				row, err := reader.ReadString('\n')
				if err != nil {
					cmd.Process.Kill()
					return fmt.Errorf("failed to read row: %v", err)
				}
				row = strings.TrimSpace(row)
				parts := strings.Fields(row)
				if len(parts) != n {
					cmd.Process.Kill()
					return fmt.Errorf("bad row output: %s", row)
				}
				final[i] = make([]int, n)
				for j, p := range parts {
					var v int
					fmt.Sscan(p, &v)
					final[i][j] = v
				}
			}
			cmd.Wait()
			for i := 0; i < n; i++ {
				for j := 0; j < n; j++ {
					if final[i][j] != board[i][j] {
						return fmt.Errorf("wrong board")
					}
				}
			}
			return nil
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 1; i <= 100; i++ {
		board := genBoard(rng)
		if err := runCase(bin, board); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i, err)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
