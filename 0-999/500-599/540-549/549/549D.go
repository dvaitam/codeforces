package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   c := make([][]int, n+2)
   for i := range c {
       c[i] = make([]int, m+2)
   }
   for i := 1; i <= n; i++ {
       var line string
       fmt.Fscan(reader, &line)
       for j := 1; j <= m; j++ {
           if line[j-1] == 'W' {
               c[i][j] = 1
           } else {
               c[i][j] = -1
           }
       }
   }
   count := 0
   for i := 1; i <= n; i++ {
       for j := 1; j <= m; j++ {
           k := c[i][j] - c[i+1][j] - c[i][j+1] + c[i+1][j+1]
           if k != 0 {
               count++
           }
       }
   }
   fmt.Println(count)
}
