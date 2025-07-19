package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   rdr := bufio.NewReader(os.Stdin)
   wrt := bufio.NewWriter(os.Stdout)
   defer wrt.Flush()
   var n, m int
   fmt.Fscan(rdr, &n, &m)
   x := make([]int, n+3)
   y := make([]int, n+3)
   for i := 1; i <= n; i++ {
       fmt.Fscan(rdr, &x[i], &y[i])
   }
   x[0], y[0] = x[n], y[n]
   x[n+1], y[n+1] = x[1], y[1]
   x[n+2], y[n+2] = x[2], y[2]
   // read lengths and their indices
   orig := make(map[int][]int)
   for i := 1; i <= m; i++ {
       var t int
       fmt.Fscan(rdr, &t)
       orig[t] = append(orig[t], i)
   }
   // try both parities
   for start := 1; start <= 2; start++ {
       mp := make(map[int][]int, len(orig))
       for k, v := range orig {
           cp := make([]int, len(v))
           copy(cp, v)
           mp[k] = cp
       }
       ans := make([]int, n+1)
       for i := 1; i <= n; i++ {
           ans[i] = -1
       }
       bad := false
       for j := start; j <= n; j += 2 {
           t := abs(x[j]-x[j-1]) + abs(y[j]-y[j-1]) + abs(x[j]-x[j+1]) + abs(y[j]-y[j+1])
           list, ok := mp[t]
           if !ok || len(list) == 0 {
               bad = true
               break
           }
           ans[j] = list[len(list)-1]
           mp[t] = list[:len(list)-1]
       }
       if !bad {
           fmt.Fprintln(wrt, "YES")
           for i := 1; i <= n; i++ {
               if i > 1 {
                   wrt.WriteByte(' ')
               }
               fmt.Fprint(wrt, ans[i])
           }
           fmt.Fprintln(wrt)
           return
       }
   }
   fmt.Fprintln(wrt, "NO")
}

func abs(a int) int {
   if a < 0 {
       return -a
   }
   return a
}
