package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const bitLen = 22

type Query struct {
   l, r, id int
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   fmt.Fscan(in, &n)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   var m int
   fmt.Fscan(in, &m)
   qs := make([]Query, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &qs[i].l, &qs[i].r)
       qs[i].l--
       qs[i].r--
       qs[i].id = i
   }
   sort.Slice(qs, func(i, j int) bool {
       return qs[i].r < qs[j].r
   })

   lb := make([]int, bitLen)
   at := make([]int, bitLen)
   res := make([]int, m)

   u := 0
   for r := 0; r < n; r++ {
       x := a[r]
       pos := r
       for i := bitLen - 1; i >= 0; i-- {
           if (x>>i)&1 == 0 {
               continue
           }
           if lb[i] == 0 {
               lb[i] = x
               at[i] = pos
               break
           }
           if at[i] < pos {
               // swap lb[i], x
               lb[i], x = x, lb[i]
               at[i], pos = pos, at[i]
           }
           x ^= lb[i]
       }
       for u < m && qs[u].r == r {
           ul := qs[u].l
           id := qs[u].id
           ans := 0
           for i := bitLen - 1; i >= 0; i-- {
               if at[i] >= ul && (ans^lb[i]) > ans {
                   ans ^= lb[i]
               }
           }
           res[id] = ans
           u++
       }
   }
   for i := 0; i < m; i++ {
       fmt.Fprintln(out, res[i])
   }
}
