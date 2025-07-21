package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n, m int
   board [][]byte
   visited [][]bool
   dirs = [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
)

// dfs returns true if a cycle is found starting from (x,y)
func dfs(x, y, px, py, depth int, color byte) bool {
   visited[x][y] = true
   for _, d := range dirs {
       nx, ny := x+d[0], y+d[1]
       if nx < 0 || nx >= n || ny < 0 || ny >= m {
           continue
       }
       if board[nx][ny] != color {
           continue
       }
       if nx == px && ny == py {
           // don't go back to parent
           continue
       }
       if visited[nx][ny] {
           // found a cycle if length >= 4
           if depth+1 >= 4 {
               return true
           }
           continue
       }
       if dfs(nx, ny, x, y, depth+1, color) {
           return true
       }
   }
   return false
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   fmt.Fscan(reader, &n, &m)
   board = make([][]byte, n)
   for i := 0; i < n; i++ {
       var line string
       fmt.Fscan(reader, &line)
       board[i] = []byte(line)
   }
   visited = make([][]bool, n)
   for i := range visited {
       visited[i] = make([]bool, m)
   }
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           if !visited[i][j] {
               if dfs(i, j, -1, -1, 1, board[i][j]) {
                   fmt.Println("Yes")
                   return
               }
           }
       }
   }
   fmt.Println("No")
}
