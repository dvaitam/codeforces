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

   var N, D, U, Q int
   if _, err := fmt.Fscan(reader, &N, &D, &U, &Q); err != nil {
       return
   }
   H := make([]int, N)
   for i := 0; i < N; i++ {
       fmt.Fscan(reader, &H[i])
   }
   // read updates
   upU := make([][2]int, U)
   for i := 0; i < U; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       upU[i][0], upU[i][1] = a, b
   }
   // read queries
   type Query struct{ x, y, v, idx int }
   qs := make([]Query, Q)
   for i := 0; i < Q; i++ {
       fmt.Fscan(reader, &qs[i].x, &qs[i].y, &qs[i].v)
       qs[i].idx = i
   }
   // sort queries by day v
   ord := make([]int, Q)
   for i := range ord {
       ord[i] = i
   }
   sort.Slice(ord, func(i, j int) bool {
       return qs[ord[i]].v < qs[ord[j]].v
   })
   // adjacency sets
   adj := make([]map[int]struct{}, N)
   for i := 0; i < N; i++ {
       adj[i] = make(map[int]struct{})
   }
   ans := make([]int, Q)
   qi := 0
   // day 0 queries (before any update)
   for qi < Q && qs[ord[qi]].v == 0 {
       q := qs[ord[qi]]
       ans[q.idx] = queryAns(q.x, q.y, H, adj)
       qi++
   }
   // days 1..U
   for day := 1; day <= U; day++ {
       a := upU[day-1][0]
       b := upU[day-1][1]
       // toggle trust
       if _, ok := adj[a][b]; ok {
           delete(adj[a], b)
           delete(adj[b], a)
       } else {
           adj[a][b] = struct{}{}
           adj[b][a] = struct{}{}
       }
       // process queries for this day
       for qi < Q && qs[ord[qi]].v == day {
           q := qs[ord[qi]]
           ans[q.idx] = queryAns(q.x, q.y, H, adj)
           qi++
       }
   }
   // output answers
   for i := 0; i < Q; i++ {
       fmt.Fprintln(writer, ans[i])
   }
}

// compute answer for one query given current adj
func queryAns(x, y int, H []int, adj []map[int]struct{}) int {
   Nx := adj[x]
   Ny := adj[y]
   if len(Nx) == 0 || len(Ny) == 0 {
       return -1
   }
   sx := make([]int, 0, len(Nx))
   for u := range Nx {
       sx = append(sx, H[u])
   }
   sy := make([]int, 0, len(Ny))
   for v := range Ny {
       sy = append(sy, H[v])
   }
   sort.Ints(sx)
   sort.Ints(sy)
   i, j, best := 0, 0, int(1e18)
   for i < len(sx) && j < len(sy) {
       a, b := sx[i], sy[j]
       d := a - b
       if d < 0 {
           d = -d
       }
       if d < best {
           best = d
       }
       if sx[i] < sy[j] {
           i++
       } else {
           j++
       }
   }
   return best
}
