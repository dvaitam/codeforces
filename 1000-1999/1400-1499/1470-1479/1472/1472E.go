package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   fmt.Fscan(reader, &t)
   for t > 0 {
       t--
       var n int
       fmt.Fscan(reader, &n)
       h := make([]int, n)
       w := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &h[i], &w[i])
           if h[i] > w[i] {
               h[i], w[i] = w[i], h[i]
           }
       }
       p := make([]int, n)
       for i := 0; i < n; i++ {
           p[i] = i
       }
       sort.Slice(p, func(i, j int) bool {
           return h[p[i]] < h[p[j]]
       })
       tmp := -1
       ans := make([]int, n)
       for i := 0; i < n; i++ {
           ans[i] = -1
       }
       i := 0
       for i < n {
           j := i
           hh := h[p[i]]
           for j < n && h[p[j]] == hh {
               j++
           }
           for k := i; k < j; k++ {
               idx := p[k]
               if tmp != -1 && w[tmp] < w[idx] {
                   ans[idx] = tmp
               }
           }
           for k := i; k < j; k++ {
               idx := p[k]
               if tmp == -1 || w[tmp] > w[idx] {
                   tmp = idx
               }
           }
           i = j
       }
       for i := 0; i < n; i++ {
           if ans[i] >= 0 {
               fmt.Fprint(writer, ans[i]+1, " ")
           } else {
               fmt.Fprint(writer, -1, " ")
           }
       }
       fmt.Fprintln(writer)
   }
}
