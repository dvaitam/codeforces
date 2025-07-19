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
   const maxX = 100
   faces := []int{2, 3, 4, 5, 6, 7}
   dp := make([]int, maxX+1)
   // initialize dp
   dp[0] = 0
   for i := 1; i <= maxX; i++ {
       dp[i] = maxX + 5
       for _, f := range faces {
           if i >= f && dp[i-f]+1 < dp[i] {
               dp[i] = dp[i-f] + 1
           }
       }
   }
   for i := 0; i < t; i++ {
       var x int
       fmt.Fscan(reader, &x)
       if x >= 0 && x <= maxX {
           fmt.Fprintln(writer, dp[x])
       } else {
           // out of range, but constraints guarantee 2<=x<=100
           fmt.Fprintln(writer, 0)
       }
   }
