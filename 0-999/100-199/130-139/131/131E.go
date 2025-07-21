package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   fmt.Fscan(in, &n, &m)
   queens := make([]struct{ r, c, att int }, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &queens[i].r, &queens[i].c)
   }
   idx := make([]int, m)
   for i := 0; i < m; i++ {
       idx[i] = i
   }
   // Horizontal (rows)
   sort.Slice(idx, func(i, j int) bool {
       if queens[idx[i]].r != queens[idx[j]].r {
           return queens[idx[i]].r < queens[idx[j]].r
       }
       return queens[idx[i]].c < queens[idx[j]].c
   })
   for i := 0; i+1 < m; i++ {
       a, b := idx[i], idx[i+1]
       if queens[a].r == queens[b].r {
           queens[a].att++
           queens[b].att++
       }
   }
   // Vertical (columns)
   for i := 0; i < m; i++ {
       idx[i] = i
   }
   sort.Slice(idx, func(i, j int) bool {
       if queens[idx[i]].c != queens[idx[j]].c {
           return queens[idx[i]].c < queens[idx[j]].c
       }
       return queens[idx[i]].r < queens[idx[j]].r
   })
   for i := 0; i+1 < m; i++ {
       a, b := idx[i], idx[i+1]
       if queens[a].c == queens[b].c {
           queens[a].att++
           queens[b].att++
       }
   }
   // Main diagonal (r-c)
   for i := 0; i < m; i++ {
       idx[i] = i
   }
   sort.Slice(idx, func(i, j int) bool {
       ai, aj := idx[i], idx[j]
       di := queens[ai].r - queens[ai].c
       dj := queens[aj].r - queens[aj].c
       if di != dj {
           return di < dj
       }
       return queens[ai].r < queens[aj].r
   })
   for i := 0; i+1 < m; i++ {
       a, b := idx[i], idx[i+1]
       if queens[a].r-queens[a].c == queens[b].r-queens[b].c {
           queens[a].att++
           queens[b].att++
       }
   }
   // Anti-diagonal (r+c)
   for i := 0; i < m; i++ {
       idx[i] = i
   }
   sort.Slice(idx, func(i, j int) bool {
       ai, aj := idx[i], idx[j]
       si := queens[ai].r + queens[ai].c
       sj := queens[aj].r + queens[aj].c
       if si != sj {
           return si < sj
       }
       return queens[ai].r < queens[aj].r
   })
   for i := 0; i+1 < m; i++ {
       a, b := idx[i], idx[i+1]
       if queens[a].r+queens[a].c == queens[b].r+queens[b].c {
           queens[a].att++
           queens[b].att++
       }
   }
   // Count results
   t := make([]int, 9)
   for i := 0; i < m; i++ {
       if queens[i].att >= 0 && queens[i].att <= 8 {
           t[queens[i].att]++
       }
   }
   out := bufio.NewWriter(os.Stdout)
   for i := 0; i < 9; i++ {
       if i > 0 {
           fmt.Fprint(out, " ")
       }
       fmt.Fprint(out, t[i])
   }
   fmt.Fprintln(out)
   out.Flush()
}
