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
   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for t > 0 {
       t--
       var n, m int
       fmt.Fscan(reader, &n, &m)
       grid := make([]string, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &grid[i])
       }
       hor := make([][]int, m+2)
       ver := make([][]int, n+2)
       for i := 1; i <= n; i++ {
           row := grid[i-1]
           for j := 1; j <= m; j++ {
               c := row[j-1]
               if c == 'U' {
                   ver[i] = append(ver[i], j)
               } else if c == 'L' {
                   hor[j] = append(hor[j], i)
               }
           }
       }
       board := make([][]byte, n+2)
       for i := range board {
           board[i] = make([]byte, m+2)
       }
       fail := false
       // fill horizontal pairs
       for j := 1; j <= m; j++ {
           if len(hor[j])%2 != 0 {
               fail = true
               break
           }
           for idx, irow := range hor[j] {
               if idx%2 == 0 {
                   board[irow][j] = 'B'
                   board[irow][j+1] = 'W'
               } else {
                   board[irow][j] = 'W'
                   board[irow][j+1] = 'B'
               }
           }
       }
       // fill vertical pairs
       if !fail {
           for i := 1; i <= n; i++ {
               if len(ver[i])%2 != 0 {
                   fail = true
                   break
               }
               for idx, jcol := range ver[i] {
                   if idx%2 == 0 {
                       board[i][jcol] = 'B'
                       board[i+1][jcol] = 'W'
                   } else {
                       board[i][jcol] = 'W'
                       board[i+1][jcol] = 'B'
                   }
               }
           }
       }
       // output result
       if fail {
           writer.WriteString("-1\n")
       } else {
           for i := 1; i <= n; i++ {
               for j := 1; j <= m; j++ {
                   writer.WriteByte(board[i][j])
               }
               writer.WriteByte('\n')
           }
       }
   }
}
