package main

import (
	"bufio"
	"fmt"
	"os"
)

func c(h, w int) int {
	tot := h - 1
	for i := 0; i < h; i++ {
		if i%2 == 0 {
			tot += w - 1
		} else {
			if i+1 < h {
				tot += 2 * (w - 1)
			} else {
				tot += w - 1
			}
		}
	}
	return tot
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var h, w, k int
	if _, err := fmt.Fscan(reader, &h, &w, &k); err != nil {
		return
	}
	// special case k == 1
	if k == 1 {
		totE := 0
		for i := 0; i < 2*h-1; i++ {
			var s string
			fmt.Fscan(reader, &s)
			for _, ch := range s {
				if ch == 'E' {
					totE++
				}
			}
		}
		need := 3 * (w*(h-1) + h*(w-1))
		if totE*4 >= need {
			fmt.Println("YES")
			for i := 0; i < h; i++ {
				for j := 0; j < w; j++ {
					if j > 0 {
						fmt.Print(" ")
					}
					fmt.Print(1)
				}
				fmt.Println()
			}
		} else {
			fmt.Println("NO")
		}
		return
	}
	// read row and col constraints
	row := make([]string, h)
	col := make([]string, h-1)
	for i := 0; i < 2*h-1; i++ {
		if i%2 == 0 {
			fmt.Fscan(reader, &row[i/2])
		} else {
			fmt.Fscan(reader, &col[i/2])
		}
	}
	// board initialization
	board := make([][]int, h)
	for i := range board {
		board[i] = make([]int, w)
	}
	board[0][0] = 1

	// fill functions
	fill1 := func(i, j int) {
		if j > 0 && board[i][j-1] != 0 {
			if row[i][j-1] == 'E' {
				board[i][j] = board[i][j-1]
			} else {
				board[i][j] = 3 - board[i][j-1]
			}
		} else if i > 0 && board[i-1][j] != 0 {
			if col[i-1][j] == 'E' {
				board[i][j] = board[i-1][j]
			} else {
				board[i][j] = 3 - board[i-1][j]
			}
		}
	}
	fill2 := func(r, c0 int) {
		var poll [3]int
		// up
		if r > 0 && board[r-1][c0] != 0 {
			if col[r-1][c0] == 'E' {
				poll[board[r-1][c0]]++
			} else {
				poll[3-board[r-1][c0]]++
			}
		}
		// left
		if c0 > 0 && board[r][c0-1] != 0 {
			if row[r][c0-1] == 'E' {
				poll[board[r][c0-1]]++
			} else {
				poll[3-board[r][c0-1]]++
			}
		}
		// down
		if r+1 < h && board[r+1][c0] != 0 {
			if col[r][c0] == 'E' {
				poll[board[r+1][c0]]++
			} else {
				poll[3-board[r+1][c0]]++
			}
		}
		// right
		if c0+1 < w && board[r][c0+1] != 0 {
			if row[r][c0] == 'E' {
				poll[board[r][c0+1]]++
			} else {
				poll[3-board[r][c0+1]]++
			}
		}
		if poll[1] >= poll[2] {
			board[r][c0] = 1
		} else {
			board[r][c0] = 2
		}
	}

	// decide orientation
	if c(h, w) >= c(w, h) {
		for i := 1; i < h; i++ {
			fill1(i, 0)
		}
		for i := 0; i < h; i++ {
			if i%2 == 0 {
				for j := 1; j < w; j++ {
					fill1(i, j)
				}
			}
		}
		for i := 0; i < h; i++ {
			if i%2 != 0 {
				for j := 1; j < w; j++ {
					fill2(i, j)
				}
			}
		}
	} else {
		for j := 1; j < w; j++ {
			fill1(0, j)
		}
		for j := 0; j < w; j++ {
			if j%2 == 0 {
				for i := 1; i < h; i++ {
					fill1(i, j)
				}
			}
		}
		for j := 0; j < w; j++ {
			if j%2 != 0 {
				for i := 1; i < h; i++ {
					fill2(i, j)
				}
			}
		}
	}
	// output
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	fmt.Fprintln(writer, "YES")
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			if j > 0 {
				writer.WriteString(" ")
			}
			writer.WriteString(fmt.Sprint(board[i][j]))
		}
		writer.WriteByte('\n')
	}
}
