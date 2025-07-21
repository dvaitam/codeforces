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

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   a := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // ans[i] = number of distinct elements in a[i..n]
   ans := make([]int, n+2)
   // seen values up to max a_i (<= 100000)
   const maxA = 100000
   seen := make([]bool, maxA+1)
   distinct := 0
   for i := n; i >= 1; i-- {
       v := a[i]
       if !seen[v] {
           seen[v] = true
           distinct++
       }
       ans[i] = distinct
   }
   // process queries
   for i := 0; i < m; i++ {
       var l int
       fmt.Fscan(reader, &l)
       // l is between 1 and n
       fmt.Fprintln(writer, ans[l])
   }
}
