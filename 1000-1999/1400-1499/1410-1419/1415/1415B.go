package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var t int
   fmt.Fscan(reader, &t)
   for ; t > 0; t-- {
       var n, k int
       fmt.Fscan(reader, &n, &k)
       c := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &c[i])
       }
       // Track unique colors
       present := make(map[int]struct{})
       for _, v := range c {
           present[v] = struct{}{}
       }
       ans := n + 1
       // For each candidate color, compute days
       for color := range present {
           days := 0
           for i := 0; i < n; {
               if c[i] == color {
                   i++
               } else {
                   days++
                   i += k
               }
           }
           ans = min(ans, days)
       }
       fmt.Fprintln(writer, ans)
   }
}
