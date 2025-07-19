package main

import (
   "bufio"
   "fmt"
   "os"
)

type edge struct {
   to, id int
}

func main() {
   r := bufio.NewReader(os.Stdin)
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()

   var n int
   fmt.Fscan(r, &n)
   adj := make([][]edge, n+1)
   we := make([]int64, n-1)
   for i := 0; i < n-1; i++ {
       var x, y int
       var wgt int64
       fmt.Fscan(r, &x, &y, &wgt)
       we[i] = wgt
       adj[x] = append(adj[x], edge{y, i})
       adj[y] = append(adj[y], edge{x, i})
   }

   hag := make([]float64, n-1)
   size := make([]int, n+1)
   poop := float64(n) * float64(n-1) * float64(n-2) / 6.0

   type item struct{ u, p, eid, stage int }
   stack := make([]item, 0, 2*n)
   stack = append(stack, item{1, 0, -1, 0})
   for len(stack) > 0 {
       it := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       u, p, eid, stage := it.u, it.p, it.eid, it.stage
       if stage == 0 {
           stack = append(stack, item{u, p, eid, 1})
           for _, e := range adj[u] {
               if e.to == p {
                   continue
               }
               stack = append(stack, item{e.to, u, e.id, 0})
           }
       } else {
           size[u] = 1
           for _, e := range adj[u] {
               if e.to == p {
                   continue
               }
               size[u] += size[e.to]
           }
           if eid != -1 {
               ne := float64(size[u])
               op := float64(n) - ne
               hag[eid] = ((op*(op-1)/2.0)*ne/poop + (ne*(ne-1)/2.0)*op/poop) * 2.0
           }
       }
   }

   var ans float64
   for i := 0; i < n-1; i++ {
       ans += hag[i] * float64(we[i])
   }

   var q int
   fmt.Fscan(r, &q)
   for i := 0; i < q; i++ {
       var x int
       var nw int64
       fmt.Fscan(r, &x, &nw)
       eid := x - 1
       ans += hag[eid] * float64(nw - we[eid])
       we[eid] = nw
       fmt.Fprintf(w, "%.15f\n", ans)
   }
}
