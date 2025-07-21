package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "strconv"
)

// toposort returns a topological ordering of graph with n nodes (1..n) given adjacency lists and indegrees.
// If a cycle exists, returns (nil, false).
func toposort(n int, adj [][]int, indeg []int) ([]int, bool) {
   order := make([]int, 0, n)
   q := make([]int, 0, n)
   // init queue
   for i := 1; i <= n; i++ {
       if indeg[i] == 0 {
           q = append(q, i)
       }
   }
   for qi := 0; qi < len(q); qi++ {
       u := q[qi]
       order = append(order, u)
       for _, v := range adj[u] {
           indeg[v]--
           if indeg[v] == 0 {
               q = append(q, v)
           }
       }
   }
   if len(order) != n {
       return nil, false
   }
   return order, true
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   // Edge groups
   var LR [][2]int
   var TB [][2]int
   var LT, LB, RT, BR [][2]int
   // Read strings
   for i := 0; i < n+m; i++ {
       var a, b string
       var p, q int
       fmt.Fscan(in, &a, &p, &b, &q)
       // normalize (a,p) as first side
       side1, side2 := a, b
       x, y := p, q
       // ensure consistent ordering: we'll handle both orders
       if side1 == "L" {
           if side2 == "R" {
               LR = append(LR, [2]int{x, y})
               continue
           } else if side2 == "T" {
               LT = append(LT, [2]int{x, y})
               continue
           } else if side2 == "B" {
               LB = append(LB, [2]int{x, y})
               continue
           }
       }
       if side1 == "R" {
           if side2 == "L" {
               // flip
               LR = append(LR, [2]int{y, x})
               continue
           } else if side2 == "T" {
               RT = append(RT, [2]int{x, y})
               continue
           } else if side2 == "B" {
               BR = append(BR, [2]int{y, x}) // store B->R as (b, uR)
               continue
           }
       }
       if side1 == "T" {
           if side2 == "B" {
               TB = append(TB, [2]int{q, p}) // store B->T: use p for T, q for B, but we need B->T: (b, t)
               // Actually q is q(B), p is p(T)? side1=T p, side2=B q => B=q, T=p
               continue
           } else if side2 == "L" {
               LT = append(LT, [2]int{q, p})
               continue
           } else if side2 == "R" {
               RT = append(RT, [2]int{y, x})
               continue
           }
       }
       if side1 == "B" {
           if side2 == "T" {
               TB = append(TB, [2]int{x, y}) // side1=B x, side2=T y
               continue
           } else if side2 == "L" {
               LB = append(LB, [2]int{y, x})
               continue
           } else if side2 == "R" {
               BR = append(BR, [2]int{x, y})
               continue
           }
       }
       // should not reach
   }
   // Prepare row constraints
   rowAdj := make([][]int, n+1)
   rowIndeg := make([]int, n+1)
   // LR: sort by L increasing
   sortPair := func(a [][2]int, asc bool) {
       if asc {
           // increasing on first
           sort.Slice(a, func(i, j int) bool { return a[i][0] < a[j][0] })
       } else {
           sort.Slice(a, func(i, j int) bool { return a[i][0] > a[j][0] })
       }
   }
   // use built-in sort
   // LR
   sort.Slice(LR, func(i, j int) bool { return LR[i][0] < LR[j][0] })
   for i := 0; i+1 < len(LR); i++ {
       u := LR[i][1]
       v := LR[i+1][1]
       rowAdj[u] = append(rowAdj[u], v)
       rowIndeg[v]++
   }
   // BR: sort by B inc
   sort.Slice(BR, func(i, j int) bool { return BR[i][0] < BR[j][0] })
   for i := 0; i+1 < len(BR); i++ {
       // v_j+1 -> v_j
       u1 := BR[i][1]
       u2 := BR[i+1][1]
       // need row_pos[u1] > row_pos[u2] => u2->u1
       rowAdj[u2] = append(rowAdj[u2], u1)
       rowIndeg[u1]++
   }
   // Toposort rows
   rowOrder, ok := toposort(n, rowAdj, rowIndeg)
   if !ok {
       fmt.Println("No solution")
       return
   }
   // Prepare col constraints
   colAdj := make([][]int, m+1)
   colIndeg := make([]int, m+1)
   // LB: sort by L inc
   sort.Slice(LB, func(i, j int) bool { return LB[i][0] < LB[j][0] })
   for i := 0; i+1 < len(LB); i++ {
       u := LB[i][1]
       v := LB[i+1][1]
       colAdj[u] = append(colAdj[u], v)
       colIndeg[v]++
   }
   // LT: sort by L inc, dec on T
   sort.Slice(LT, func(i, j int) bool { return LT[i][0] < LT[j][0] })
   for i := 0; i+1 < len(LT); i++ {
       u := LT[i+1][1]
       v := LT[i][1]
       colAdj[u] = append(colAdj[u], v)
       colIndeg[v]++
   }
   // RT: sort by R dec, dec on T
   sort.Slice(RT, func(i, j int) bool { return RT[i][0] > RT[j][0] })
   for i := 0; i+1 < len(RT); i++ {
       u := RT[i+1][1]
       v := RT[i][1]
       colAdj[u] = append(colAdj[u], v)
       colIndeg[v]++
   }
   // TB: sort by B inc, dec on T
   sort.Slice(TB, func(i, j int) bool { return TB[i][0] < TB[j][0] })
   for i := 0; i+1 < len(TB); i++ {
       u := TB[i+1][1]
       v := TB[i][1]
       colAdj[u] = append(colAdj[u], v)
       colIndeg[v]++
   }
   // Toposort cols
   colOrder, ok := toposort(m, colAdj, colIndeg)
   if !ok {
       fmt.Println("No solution")
       return
   }
   // Output row order and col order
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for i, v := range rowOrder {
       if i > 0 {
           out.WriteByte(' ')
       }
       out.WriteString(strconv.Itoa(v))
   }
   out.WriteByte('\n')
   for i, v := range colOrder {
       if i > 0 {
           out.WriteByte(' ')
       }
       out.WriteString(strconv.Itoa(v))
   }
   out.WriteByte('\n')
}
