package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m, k int
   if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
       return
   }
   // occupied seats
   occ := make([][]bool, n+1)
   for i := 0; i <= n; i++ {
       occ[i] = make([]bool, m+1)
   }
   for t := 0; t < k; t++ {
       var x1, y1 int
       fmt.Fscan(in, &x1, &y1)
       if !occ[x1][y1] {
           occ[x1][y1] = true
           fmt.Fprintln(out, x1, y1)
           continue
       }
       // search nearest free seat
       found := false
       var rx, ry int
       // increasing Manhattan distance d
       for d := 1; !found; d++ {
           // dx from 0 to d
           for dx := 0; dx <= d && !found; dx++ {
               dy := d - dx
               // candidate rows in lex order
               rows := []int{}
               if dx == 0 {
                   if x1 >= 1 && x1 <= n {
                       rows = append(rows, x1)
                   }
               } else {
                   if x1-dx >= 1 {
                       rows = append(rows, x1-dx)
                   }
                   if x1+dx <= n {
                       rows = append(rows, x1+dx)
                   }
               }
               for _, x2 := range rows {
                   // candidate columns in lex order
                   cols := []int{}
                   if dy == 0 {
                       if y1 >= 1 && y1 <= m {
                           cols = append(cols, y1)
                       }
                   } else {
                       if y1-dy >= 1 {
                           cols = append(cols, y1-dy)
                       }
                       if y1+dy <= m {
                           cols = append(cols, y1+dy)
                       }
                   }
                   for _, y2 := range cols {
                       if !occ[x2][y2] {
                           found = true
                           rx, ry = x2, y2
                           break
                       }
                   }
                   if found {
                       break
                   }
               }
           }
       }
       occ[rx][ry] = true
       fmt.Fprintln(out, rx, ry)
   }
}
