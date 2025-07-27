package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 998244353

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for t > 0 {
       t--
       var n, k int
       fmt.Fscan(reader, &n, &k)
       a := make([]int, n+1)
       pos := make([]int, n+1)
       for i := 1; i <= n; i++ {
           fmt.Fscan(reader, &a[i])
           pos[a[i]] = i
       }
       b := make([]int, k+1)
       need := make([]int, n+1)
       for i := 1; i <= k; i++ {
           fmt.Fscan(reader, &b[i])
           need[b[i]] = i
       }
       ans := 1
       for i := 1; i <= k; i++ {
           x := b[i]
           idx := pos[x]
           cnt := 0
           // check left neighbor
           if idx > 1 {
               y := a[idx-1]
               if need[y] < i {
                   cnt++
               }
           }
           // check right neighbor
           if idx < n {
               y := a[idx+1]
               if need[y] < i {
                   cnt++
               }
           }
           ans = ans * cnt % mod
       }
       fmt.Fprintln(writer, ans)
   }
}
