package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   a := make([][]byte, n+2)
   s := make([][]byte, n+2)
   for i := range a {
       a[i] = make([]byte, m+2)
       s[i] = make([]byte, m+2)
   }
   for i := 1; i <= n; i++ {
       var row string
       fmt.Fscan(in, &row)
       for j := 1; j <= m; j++ {
           a[i][j] = row[j-1]
           s[i][j] = '.'
       }
   }
   ok := true
   for i := 1; i <= n; i++ {
       for j := 1; j <= m; j++ {
           if a[i][j] != 'b' {
               continue
           }
           cI := 'a' + 2*((i/2+j/2)%2) + (i % 2)
           // vertical placement
           if a[i-1][j] == 'w' && a[i+1][j] == 'w' {
               c := byte(cI)
               s[i-1][j], s[i][j], s[i+1][j] = c, c, c
               a[i-1][j], a[i][j], a[i+1][j] = '.', '.', '.'
           } else if a[i][j-1] == 'w' && a[i][j+1] == 'w' {
               c := byte(cI)
               s[i][j-1], s[i][j], s[i][j+1] = c, c, c
               a[i][j-1], a[i][j], a[i][j+1] = '.', '.', '.'
           }
       }
   }
   for i := 1; i <= n && ok; i++ {
       for j := 1; j <= m; j++ {
           if a[i][j] != '.' {
               ok = false
               break
           }
       }
   }
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   if ok {
       fmt.Fprintln(out, "YES")
       for i := 1; i <= n; i++ {
           for j := 1; j <= m; j++ {
               out.WriteByte(s[i][j])
           }
           out.WriteByte('\n')
       }
   } else {
       fmt.Fprintln(out, "NO")
   }
}
