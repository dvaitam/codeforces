package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var T int
   fmt.Fscan(in, &T)
   for T > 0 {
       T--
       var n int
       fmt.Fscan(in, &n)
       a := make([]int, n+1)
       b := make([]int, n+1)
       c := make([]int, n+1)
       d := make([]int, n+1)
       num := make([]int, n+1)
       res := make([]int, n+1)
       vis := make([]bool, n+1)

       for i := 1; i <= n; i++ {
           fmt.Fscan(in, &a[i])
           if c[a[i]] == 0 {
               c[a[i]] = i
           } else {
               d[a[i]] = i
           }
       }
       for i := 1; i <= n; i++ {
           fmt.Fscan(in, &b[i])
           if c[b[i]] == 0 {
               c[b[i]] = i
           } else {
               d[b[i]] = i
           }
       }
       for i := 1; i <= n; i++ {
           num[a[i]]++
           num[b[i]]++
       }
       bad := false
       for i := 1; i <= n; i++ {
           if num[i] != 2 {
               bad = true
               break
           }
       }
       if bad {
           fmt.Fprintln(out, -1)
           continue
       }
       // process cycles
       for i := 1; i <= n; i++ {
           if vis[i] {
               continue
           }
           cur := i
           x := c[i]
           buf := []int{x}
           cnt := 0
           if a[x] != i {
               res[x] = 1
               cnt++
           }
           // traverse cycle
           for {
               vis[cur] = true
               nxt := a[x] + b[x] - cur
               x = c[nxt] + d[nxt] - x
               cur = nxt
               if vis[cur] {
                   break
               }
               buf = append(buf, x)
               if b[x] == nxt {
                   res[x] = 1
                   cnt++
               }
           }
           if cnt > len(buf)-cnt {
               for _, v := range buf {
                   if res[v] == 1 {
                       res[v] = 0
                   } else {
                       res[v] = 1
                   }
               }
           }
       }
       // collect result
       var outIdx []int
       for i := 1; i <= n; i++ {
           if res[i] == 1 {
               outIdx = append(outIdx, i)
           }
       }
       fmt.Fprintln(out, len(outIdx))
       for i, v := range outIdx {
           if i > 0 {
               out.WriteByte(' ')
           }
           fmt.Fprint(out, v)
       }
       out.WriteByte('\n')
   }
}
