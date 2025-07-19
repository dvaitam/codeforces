package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

type Edge struct {
   to int
   w  float64
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
      return
   }
   adj := make([][]Edge, n+1)
   for i := 0; i < m; i++ {
      var a, b int
      var c float64
      fmt.Fscan(reader, &a, &b, &c)
      adj[a] = append(adj[a], Edge{to: b, w: c})
      adj[b] = append(adj[b], Edge{to: a, w: c})
   }
   const InfCoeff = int(1e9)
   coeff := make([]int, n+1)
   cons := make([]float64, n+1)
   ans := make([]float64, n+1)
   for i := 1; i <= n; i++ {
      coeff[i] = InfCoeff
   }
   // process each component
   for i := 1; i <= n; i++ {
      if coeff[i] != InfCoeff {
         continue
      }
      // BFS
      queue := []int{i}
      coeff[i] = 1
      cons[i] = 0
      vis := make([]int, 0, 1)
      temp := make([]float64, 0, 1)
      x := 0.0
      xSet := false
      for qi := 0; qi < len(queue); qi++ {
         v := queue[qi]
         // record for median
         temp = append(temp, -cons[v]*float64(coeff[v]))
         vis = append(vis, v)
         for _, e := range adj[v] {
            u := e.to
            c := e.w
            if coeff[u] == InfCoeff {
               coeff[u] = -coeff[v]
               cons[u] = c - cons[v]
               queue = append(queue, u)
            } else {
               sumC := float64(coeff[u]+coeff[v])
               sumCons := cons[u] + cons[v]
               if int(sumC) == 0 {
                  if math.Abs(sumCons-c) > 1e-8 {
                     fmt.Fprintln(writer, "NO")
                     return
                  }
                  continue
               }
               val := (c - sumCons) / sumC
               if xSet {
                  if math.Abs(val-x) > 1e-8 {
                     fmt.Fprintln(writer, "NO")
                     return
                  }
               } else {
                  x = val
                  xSet = true
               }
            }
         }
      }
      if !xSet {
         sort.Float64s(temp)
         x = temp[len(temp)/2]
      }
      for _, u := range vis {
         ans[u] = x*float64(coeff[u]) + cons[u]
      }
   }
   // output
   fmt.Fprintln(writer, "YES")
   for i := 1; i <= n; i++ {
      if i > 1 {
         writer.WriteByte(' ')
      }
      fmt.Fprintf(writer, "%.5f", ans[i])
   }
   writer.WriteByte('\n')
}
