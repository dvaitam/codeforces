package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type pair struct { F, S int }

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // total nodes N = n+1
   N := n + 1
   // read first node
   var firstF int
   fmt.Fscan(reader, &firstF)
   if firstF == 0 {
       fmt.Fprint(writer, -1)
       return
   }
   // queue for BFS
   q := make([]pair, 0, N)
   q = append(q, pair{firstF, 1})
   // read remaining nodes
   a := make([]pair, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i].F)
       a[i].S = i + 2
   }
   // sort by descending F
   sort.Slice(a, func(i, j int) bool {
       return a[i].F > a[j].F
   })
   // visited marker
   u := make([]bool, N+2)
   ans := make([]pair, 0)
   id := 0
   // BFS
   for head := 0; head < len(q); head++ {
       v := q[head]
       for v.F > 0 && id < len(a) {
           if !u[a[id].S] {
               v.F--
               ans = append(ans, pair{v.S, a[id].S})
               q = append(q, a[id])
               u[v.S] = true
               u[a[id].S] = true
               id++
           } else {
               break
           }
       }
   }
   // check all nodes visited
   for i := 1; i <= N; i++ {
       if !u[i] {
           fmt.Fprint(writer, -1)
           return
       }
   }
   // output
   fmt.Fprintln(writer, len(ans))
   for _, p := range ans {
       fmt.Fprintln(writer, p.F, p.S)
   }
}
