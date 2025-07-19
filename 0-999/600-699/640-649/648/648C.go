package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   grid := make([]string, n)
   a, b := -1, -1
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       grid[i] = s
       if a < 0 {
           for j := 0; j < m; j++ {
               if s[j] == 'S' {
                   a, b = i, j
               }
           }
       }
   }
   vis := make([][]bool, n)
   for i := range vis {
       vis[i] = make([]bool, m)
   }
   i, j := a, b
   var sb strings.Builder
   for {
       vis[i][j] = true
       // try moves: D, U, R, L
       if i+1 < n && grid[i+1][j] == '*' && !vis[i+1][j] {
           i++
           sb.WriteByte('D')
           continue
       }
       if i-1 >= 0 && grid[i-1][j] == '*' && !vis[i-1][j] {
           i--
           sb.WriteByte('U')
           continue
       }
       if j+1 < m && grid[i][j+1] == '*' && !vis[i][j+1] {
           j++
           sb.WriteByte('R')
           continue
       }
       if j-1 >= 0 && grid[i][j-1] == '*' && !vis[i][j-1] {
           j--
           sb.WriteByte('L')
           continue
       }
       // close the cycle by moving back to start
       if abs(i-a)+abs(j-b) == 1 {
           switch {
           case a-i == 1:
               sb.WriteByte('D')
           case a-i == -1:
               sb.WriteByte('U')
           case b-j == 1:
               sb.WriteByte('R')
           default:
               sb.WriteByte('L')
           }
           break
       }
   }
   fmt.Println(sb.String())
}
