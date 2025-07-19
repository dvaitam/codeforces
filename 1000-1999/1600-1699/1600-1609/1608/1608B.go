package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   in  = bufio.NewReader(os.Stdin)
   out = bufio.NewWriter(os.Stdout)
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}

func solve() {
   var n, a, b int
   fmt.Fscan(in, &n, &a, &b)
   if abs(a-b) > 1 || a+b > n-2 {
       fmt.Fprintln(out, -1)
       return
   }
   perm := make([]int, n)
   for i := 0; i < n; i++ {
       perm[i] = i + 1
   }
   temp := min(a, b)
   if a > b {
       i := 1
       for temp > 0 {
           perm[i], perm[i+1] = perm[i+1], perm[i]
           i += 2
           temp--
       }
       for k := i; k < n-1; k++ {
           perm[k], perm[k+1] = perm[k+1], perm[k]
       }
   } else if a < b {
       i := n - 2
       for temp > 0 {
           perm[i], perm[i-1] = perm[i-1], perm[i]
           i -= 2
           temp--
       }
       for k := i; k > 0; k-- {
           perm[k], perm[k-1] = perm[k-1], perm[k]
       }
   } else {
       i := 1
       for temp > 0 {
           perm[i], perm[i+1] = perm[i+1], perm[i]
           i += 2
           temp--
       }
   }
   for i := 0; i < n; i++ {
       if i > 0 {
           fmt.Fprint(out, " ")
       }
       fmt.Fprint(out, perm[i])
   }
   fmt.Fprintln(out)
}

func main() {
   defer out.Flush()
   var tc int
   fmt.Fscan(in, &tc)
   for tc > 0 {
       solve()
       tc--
   }
}
