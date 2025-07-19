package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var N, M int
   if _, err := fmt.Fscan(reader, &N, &M); err != nil {
       return
   }
   seg := make([]string, N)
   for i := 0; i < N; i++ {
       fmt.Fscan(reader, &seg[i])
   }
   // 8 directions: diagonals and orthogonals
   dirs := [][2]int{{1, 1}, {1, -1}, {1, 0}, {0, 1}, {0, -1}, {-1, -1}, {-1, 0}, {-1, 1}}

   // ch returns true if all 8 neighbors around (x,y) are '#'
   ch := func(x, y int) bool {
       if x > 0 && x < N-1 && y > 0 && y < M-1 {
           for _, d := range dirs {
               nx, ny := x+d[0], y+d[1]
               if seg[nx][ny] != '#' {
                   return false
               }
           }
           return true
       }
       return false
   }

   // check if cell (i,j) is part of any full 3x3 block
   check := func(i, j int) bool {
       for _, d := range dirs {
           xx, yy := i+d[0], j+d[1]
           if ch(xx, yy) {
               return true
           }
       }
       return false
   }

   for i := 0; i < N; i++ {
       for j := 0; j < M; j++ {
           if seg[i][j] == '#' {
               if !check(i, j) {
                   fmt.Fprintln(writer, "NO")
                   return
               }
           }
       }
   }
   fmt.Fprintln(writer, "YES")
}
