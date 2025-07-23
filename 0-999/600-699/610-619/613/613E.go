package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 1000000007

var (
	board [2][]byte
	word  []byte
	n, k  int
	ans   int
	used  [][]bool
)

func dfs(r, c, pos int) {
	if pos == k {
		ans++
		if ans >= MOD {
			ans -= MOD
		}
		return
	}
	used[r][c] = true
	// four directions
	dirs := [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	for _, d := range dirs {
		nr := r + d[0]
		nc := c + d[1]
		if nr < 0 || nr >= 2 || nc < 0 || nc >= n {
			continue
		}
		if used[nr][nc] {
			continue
		}
		if board[nr][nc] != word[pos] {
			continue
		}
		dfs(nr, nc, pos+1)
	}
	used[r][c] = false
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	board[0] = make([]byte, 0)
	board[1] = make([]byte, 0)
	var line string
	fmt.Fscan(reader, &line)
	board[0] = []byte(line)
	fmt.Fscan(reader, &line)
	board[1] = []byte(line)
	n = len(board[0])
	fmt.Fscan(reader, &line) // read empty line or maybe not
	if len(line) != 0 && len(board[1]) != n {
		// second line may have been blank line due to scanning
		// but problem guarantees an empty line
	}
	fmt.Fscan(reader, &line)
	word = []byte(line)
	k = len(word)

	used = make([][]bool, 2)
	for i := 0; i < 2; i++ {
		used[i] = make([]bool, n)
	}

	ans = 0
	if k == 0 {
		fmt.Println(0)
		return
	}
	for r := 0; r < 2; r++ {
		for c := 0; c < n; c++ {
			if board[r][c] == word[0] {
				dfs(r, c, 1)
			}
		}
	}
	fmt.Println(ans % MOD)
}
