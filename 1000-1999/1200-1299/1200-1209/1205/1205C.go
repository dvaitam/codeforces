package main

import (
    "bufio"
    "fmt"
    "os"
)

var (
    in  = bufio.NewReader(os.Stdin)
    out = bufio.NewWriter(os.Stdout)
)

func ask(x1, y1, x2, y2 int) int {
    fmt.Fprintf(out, "? %d %d %d %d\n", x1, y1, x2, y2)
    out.Flush()
    var res int
    fmt.Fscan(in, &res)
    return res
}

func bfs(board [][]int, startX, startY int, forward bool) {
    n := len(board)
    type pair struct{ x, y int }
    q := []pair{{startX, startY}}
    for len(q) > 0 {
        cur := q[0]
        q = q[1:]
        x, y := cur.x, cur.y
        if forward {
            dirs := [][2]int{{2,0},{0,2},{1,1}}
            for _, d := range dirs {
                nx, ny := x+d[0], y+d[1]
                if nx >= 0 && nx < n && ny >= 0 && ny < n && board[nx][ny] == -1 {
                    res := ask(x+1, y+1, nx+1, ny+1)
                    if res == 1 {
                        board[nx][ny] = board[x][y]
                    } else {
                        board[nx][ny] = 1 - board[x][y]
                    }
                    q = append(q, pair{nx, ny})
                }
            }
        } else {
            dirs := [][2]int{{-2,0},{0,-2},{-1,-1}}
            for _, d := range dirs {
                nx, ny := x+d[0], y+d[1]
                if nx >= 0 && nx < n && ny >= 0 && ny < n && board[nx][ny] == -1 {
                    res := ask(nx+1, ny+1, x+1, y+1)
                    if res == 1 {
                        board[nx][ny] = board[x][y]
                    } else {
                        board[nx][ny] = 1 - board[x][y]
                    }
                    q = append(q, pair{nx, ny})
                }
            }
        }
    }
}

func computeDP(board [][]int) [][][][]bool {
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

func main() {
    defer out.Flush()
    var n int
    fmt.Fscan(in, &n)

    board := make([][]int, n)
    for i := range board {
        board[i] = make([]int, n)
        for j := range board[i] {
            board[i][j] = -1
        }
    }
    board[0][0] = 1
    board[n-1][n-1] = 0

    bfs(board, 0, 0, true)
    bfs(board, n-1, n-1, false)

    board2 := make([][]int, n)
    for i := range board2 {
        board2[i] = make([]int, n)
        for j := range board2[i] {
            if (i+j)%2 == 0 {
                board2[i][j] = board[i][j]
            } else {
                board2[i][j] = 1 - board[i][j]
            }
        }
    }

    dp1 := computeDP(board)
    dp2 := computeDP(board2)

    var x1, y1, x2, y2 int
    found := false
    for i1 := 0; i1 < n && !found; i1++ {
        for j1 := 0; j1 < n && !found; j1++ {
            for i2 := 0; i2 < n && !found; i2++ {
                for j2 := 0; j2 < n && !found; j2++ {
                    if i2 < i1 || j2 < j1 || i1+j1+2 > i2+j2 {
                        continue
                    }
                    if dp1[i1][j1][i2][j2] != dp2[i1][j1][i2][j2] {
                        x1, y1, x2, y2 = i1, j1, i2, j2
                        found = true
                    }
                }
            }
        }
    }

    res := ask(x1+1, y1+1, x2+1, y2+1)
    finalBoard := board
    if res != 0 {
        if !dp1[x1][y1][x2][y2] {
            finalBoard = board2
        }
    } else {
        if dp1[x1][y1][x2][y2] {
            finalBoard = board2
        }
    }

    fmt.Fprintln(out, "!")
    for i := 0; i < n; i++ {
        for j := 0; j < n; j++ {
            fmt.Fprint(out, finalBoard[i][j])
            if j+1 == n {
                fmt.Fprintln(out)
            } else {
                fmt.Fprint(out, " ")
            }
        }
    }
}

