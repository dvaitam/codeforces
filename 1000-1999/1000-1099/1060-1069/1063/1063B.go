package main

import (
   "bufio"
   "fmt"
   "os"
)

type cell struct {
   r, c, lsteps, rsteps int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   var sr, sc int
   var x, y int
   fmt.Fscan(reader, &n, &m)
   fmt.Fscan(reader, &sr, &sc)
   fmt.Fscan(reader, &x, &y)
   maze := make([][]byte, n)
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       maze[i] = []byte(s)
   }
   visited := make([][]bool, n)
   for i := 0; i < n; i++ {
       visited[i] = make([]bool, m)
   }
   q := make([]cell, 0, n*m)
   head := 0
   // 0-based start
   sr--
   sc--
   q = append(q, cell{sr, sc, 0, 0})
   visited[sr][sc] = true
   ans := 1
   for head < len(q) {
       cur := q[head]
       head++
       r, c := cur.r, cur.c
       ls, rs := cur.lsteps, cur.rsteps
       // move left
       if c > 0 && ls < x && !visited[r][c-1] && maze[r][c-1] == '.' {
           visited[r][c-1] = true
           ans++
           q = append(q, cell{r, c - 1, ls + 1, rs})
       }
       // move right
       if c+1 < m && rs < y && !visited[r][c+1] && maze[r][c+1] == '.' {
           visited[r][c+1] = true
           ans++
           q = append(q, cell{r, c + 1, ls, rs + 1})
       }
       // move up
       for i := r - 1; i >= 0; i-- {
           if visited[i][c] || maze[i][c] != '.' {
               break
           }
           visited[i][c] = true
           ans++
           // left from new cell
           if c > 0 && ls < x && !visited[i][c-1] && maze[i][c-1] == '.' {
               visited[i][c-1] = true
               ans++
               q = append(q, cell{i, c - 1, ls + 1, rs})
           }
           // right from new cell
           if c+1 < m && rs < y && !visited[i][c+1] && maze[i][c+1] == '.' {
               visited[i][c+1] = true
               ans++
               q = append(q, cell{i, c + 1, ls, rs + 1})
           }
       }
       // move down
       for i := r + 1; i < n; i++ {
           if visited[i][c] || maze[i][c] != '.' {
               break
           }
           visited[i][c] = true
           ans++
           // left from new cell
           if c > 0 && ls < x && !visited[i][c-1] && maze[i][c-1] == '.' {
               visited[i][c-1] = true
               ans++
               q = append(q, cell{i, c - 1, ls + 1, rs})
           }
           // right from new cell
           if c+1 < m && rs < y && !visited[i][c+1] && maze[i][c+1] == '.' {
               visited[i][c+1] = true
               ans++
               q = append(q, cell{i, c + 1, ls, rs + 1})
           }
       }
   }
   fmt.Println(ans)
}
