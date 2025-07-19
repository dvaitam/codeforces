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
   for tc := 0; tc < t; tc++ {
       var n int
       fmt.Fscan(reader, &n)
       grid := make([]string, n)
       k := 0
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &grid[i])
           for j := 0; j < n; j++ {
               if grid[i][j] != '.' {
                   k++
               }
           }
       }
       // possible assignments for classes 0,1,2
       choices := [][3]byte{
           {'*', 'O', 'X'},
           {'*', 'X', 'O'},
           {'O', '*', 'X'},
           {'X', '*', 'O'},
           {'O', 'X', '*'},
           {'X', 'O', '*'},
       }
       var use [3]byte
       limit := k / 3
       // find a valid choice
       for _, ch := range choices {
           cnt := 0
           for i := 0; i < n; i++ {
               for j := 0; j < n; j++ {
                   c := grid[i][j]
                   if c == '.' {
                       continue
                   }
                   d := (i + j) % 3
                   if ch[d] != '*' && c != ch[d] {
                       cnt++
                   }
               }
           }
           if cnt <= limit {
               use = ch
               break
           }
       }
       // output modified grid
       for i := 0; i < n; i++ {
           row := []byte(grid[i])
           for j := 0; j < n; j++ {
               c := row[j]
               if c == '.' {
                   // no change
               } else {
                   d := (i + j) % 3
                   if use[d] != '*' {
                       row[j] = use[d]
                   }
               }
           }
           writer.Write(row)
           writer.WriteByte('\n')
       }
   }
}
