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
   var T int
   if _, err := fmt.Fscan(reader, &T); err != nil {
       return
   }
   for tc := 0; tc < T; tc++ {
       solve(reader, writer)
   }
}

func solve(reader *bufio.Reader, writer *bufio.Writer) {
   var n, m int
   fmt.Fscan(reader, &n, &m)
   grid := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &grid[i])
   }
   flg := [2]bool{true, true}
   for i := 0; i < n; i++ {
       row := grid[i]
       for j := 0; j < m; j++ {
           bit := (i ^ j) & 1
           switch row[j] {
           case 'R':
               flg[bit] = false
           case 'W':
               flg[bit^1] = false
           }
       }
   }
   if flg[0] {
       fmt.Fprintln(writer, "YES")
       for i := 0; i < n; i++ {
           for j := 0; j < m; j++ {
               if (i^j)&1 == 1 {
                   writer.WriteByte('R')
               } else {
                   writer.WriteByte('W')
               }
           }
           writer.WriteByte('\n')
       }
   } else if flg[1] {
       fmt.Fprintln(writer, "YES")
       for i := 0; i < n; i++ {
           for j := 0; j < m; j++ {
               if (i^j)&1 == 1 {
                   writer.WriteByte('W')
               } else {
                   writer.WriteByte('R')
               }
           }
           writer.WriteByte('\n')
       }
   } else {
       fmt.Fprintln(writer, "NO")
   }
}
