package main

import (
   "bufio"
   "fmt"
   "os"
)

type Cross struct {
   x, y, s int
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   // grid with padding
   grid := make([][]byte, n+2)
   for i := range grid {
       grid[i] = make([]byte, m+2)
       for j := range grid[i] {
           grid[i][j] = '.'
       }
   }
   // read grid rows
   for i := 1; i <= n; i++ {
       var row string
       fmt.Fscan(in, &row)
       for j := 1; j <= m; j++ {
           grid[i][j] = row[j-1]
       }
   }
   was := make([][]bool, n+2)
   for i := range was {
       was[i] = make([]bool, m+2)
   }
   var ans []Cross
   // find crosses
   for i := 1; i <= n; i++ {
       for j := 1; j <= m; j++ {
           if grid[i][j] != '*' {
               continue
           }
           // expand
           s := 0
           for d := 1; ; d++ {
               if i-d < 1 || i+d > n || j-d < 1 || j+d > m {
                   break
               }
               if grid[i-d][j] == '*' && grid[i+d][j] == '*' && grid[i][j-d] == '*' && grid[i][j+d] == '*' {
                   s = d
               } else {
                   break
               }
           }
           if s > 0 {
               ans = append(ans, Cross{i, j, s})
               // mark covered cells
               was[i][j] = true
               for d := 1; d <= s; d++ {
                   was[i-d][j] = true
                   was[i+d][j] = true
                   was[i][j-d] = true
                   was[i][j+d] = true
               }
           }
       }
   }
   // verify all stars covered
   for i := 1; i <= n; i++ {
       for j := 1; j <= m; j++ {
           if grid[i][j] == '*' && !was[i][j] {
               fmt.Fprintln(out, -1)
               return
           }
       }
   }
   // output
   fmt.Fprintln(out, len(ans))
   for _, c := range ans {
       fmt.Fprintf(out, "%d %d %d\n", c.x, c.y, c.s)
   }
}
