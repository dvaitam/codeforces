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

func solve() {
   var n, m int
   fmt.Fscan(in, &n, &m)
   p := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &p[i])
   }
   l, r := 1, n+1
   for r-l > 1 {
       mid := (l + r) >> 1
       if p[mid] <= m {
           l = mid
       } else {
           r = mid
       }
   }
   if p[l] == m {
       fmt.Fprintln(out, 0)
   } else {
       id := 0
       for i := 1; i <= n; i++ {
           if p[i] == m {
               id = i
           }
       }
       if p[l] <= m {
           fmt.Fprintln(out, 1)
           fmt.Fprintln(out, id, l)
       }
       if p[l] > m {
           fmt.Fprintln(out, 2)
           fmt.Fprintln(out, id, l)
           p[id], p[l] = p[l], p[id]
           id = l
           l, r = 1, n+1
           for r-l > 1 {
               mid := (l + r) >> 1
               if p[mid] <= m {
                   l = mid
               } else {
                   r = mid
               }
           }
           fmt.Fprintln(out, id, l)
       }
   }
}

func main() {
   defer out.Flush()
   var t int
   fmt.Fscan(in, &t)
   for i := 0; i < t; i++ {
       solve()
   }
}
