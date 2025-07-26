package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // compute sizes of active roots
   size := make([]int, n)
   activeIdx := make([]int, 0, n)
   for i := 0; i < n; i++ {
       if a[i] == -1 {
           // active startup
           size[i] = 1
           activeIdx = append(activeIdx, i)
       }
   }
   for i := 0; i < n; i++ {
       if a[i] != -1 {
           // acquired under a[i]
           size[a[i]-1]++
       }
   }
   // compute sum over pairs of sizes
   ans := 0
   m := len(activeIdx)
   for i := 0; i < m; i++ {
       si := size[activeIdx[i]]
       for j := i + 1; j < m; j++ {
           sj := size[activeIdx[j]]
           ans = (ans + si*sj) % MOD
       }
   }
   fmt.Println(ans)
}
