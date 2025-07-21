package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   parents := make([]int, n+1)
   for i := 2; i <= n; i++ {
       fmt.Fscan(reader, &parents[i])
   }
   children := make([][]int, n+1)
   for i := 2; i <= n; i++ {
       p := parents[i]
       children[p] = append(children[p], i)
   }
   dp := make([]int64, n+1)
   // process nodes in reverse order (since parents have smaller indices)
   for i := n; i >= 1; i-- {
       if len(children[i]) == 0 {
           // leaf: two painting options (white, black) + recursion (node stays red)
           dp[i] = 3
       } else {
           prod := int64(1)
           for _, c := range children[i] {
               prod = prod * dp[c] % MOD
           }
           dp[i] = (prod + 2) % MOD
       }
   }
   fmt.Fprint(writer, dp[1])
}
