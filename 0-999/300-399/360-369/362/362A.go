package main

import (
   "bufio"
   "fmt"
   "os"
)

func bfs(startR, startC int) [8][8]int {
   const N = 8
   var dist [8][8]int
   for i := 0; i < N; i++ {
       for j := 0; j < N; j++ {
           dist[i][j] = -1
       }
   }
   // moves: semiknight moves
   moves := [4][2]int{{2, 2}, {2, -2}, {-2, 2}, {-2, -2}}
   q := make([][2]int, 0, N*N)
   dist[startR][startC] = 0
   q = append(q, [2]int{startR, startC})
   for head := 0; head < len(q); head++ {
       r, c := q[head][0], q[head][1]
       for _, m := range moves {
           nr, nc := r+m[0], c+m[1]
           if nr >= 0 && nr < N && nc >= 0 && nc < N {
               if dist[nr][nc] == -1 {
                   dist[nr][nc] = dist[r][c] + 1
                   q = append(q, [2]int{nr, nc})
               }
           }
       }
   }
   return dist
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var t int
   if _, err := fmt.Fscan(in, &t); err != nil {
       return
   }
   // directions for reading lines
   for tc := 0; tc < t; tc++ {
       // read board
       board := make([]string, 8)
       for i := 0; i < 8; i++ {
           fmt.Fscan(in, &board[i])
       }
       // locate Ks
       var r1, c1, r2, c2 int
       found := 0
       for i := 0; i < 8; i++ {
           for j := 0; j < 8; j++ {
               if board[i][j] == 'K' {
                   if found == 0 {
                       r1, c1 = i, j
                   } else {
                       r2, c2 = i, j
                   }
                   found++
               }
           }
       }
       d1 := bfs(r1, c1)
       d2 := bfs(r2, c2)
       ok := false
       for i := 0; i < 8 && !ok; i++ {
           for j := 0; j < 8; j++ {
               // only good squares count for meeting
               if board[i][j] == '#' {
                   continue
               }
               di1, di2 := d1[i][j], d2[i][j]
               if di1 >= 0 && di2 >= 0 && (di1&1) == (di2&1) {
                   ok = true
                   break
               }
           }
       }
       if ok {
           fmt.Println("YES")
       } else {
           fmt.Println("NO")
       }
   }
}
