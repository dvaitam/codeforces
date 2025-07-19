package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type pair struct{ f, s int }

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   v := make([][]int, n)
   for i := 0; i < m; i++ {
       var f, s int
       fmt.Fscan(reader, &f, &s)
       f--
       s--
       v[f] = append(v[f], s)
       v[s] = append(v[s], f)
   }
   mod := 1 << 15
   // BFS from 0
   d1 := make([]int, n)
   u1 := make([]bool, n)
   for i := range d1 {
       d1[i] = mod
   }
   d1[0], u1[0] = 0, true
   queue1 := []int{0}
   for len(queue1) > 0 {
       p := queue1[0]
       queue1 = queue1[1:]
       for _, nei := range v[p] {
           if !u1[nei] {
               u1[nei] = true
               d1[nei] = d1[p] + 1
               queue1 = append(queue1, nei)
           }
       }
   }
   // BFS from n-1
   d2 := make([]int, n)
   u2 := make([]bool, n)
   for i := range d2 {
       d2[i] = mod
   }
   d2[n-1], u2[n-1] = 0, true
   queue2 := []int{n - 1}
   for len(queue2) > 0 {
       p := queue2[0]
       queue2 = queue2[1:]
       for _, nei := range v[p] {
           if !u2[nei] {
               u2[nei] = true
               d2[nei] = d2[p] + 1
               queue2 = append(queue2, nei)
           }
       }
   }
   // BFS on pairs
   d := make([][]int, n)
   used := make([][]bool, n)
   pf := make([][]int, n)
   ps := make([][]int, n)
   for i := 0; i < n; i++ {
       d[i] = make([]int, n)
       used[i] = make([]bool, n)
       pf[i] = make([]int, n)
       ps[i] = make([]int, n)
       for j := 0; j < n; j++ {
           d[i][j] = mod
       }
   }
   queue := make([]pair, 0, n*n)
   queue = append(queue, pair{0, n - 1})
   used[0][n-1] = true
   d[0][n-1] = 0
   head := 0
   found := false
   for head < len(queue) && !found {
       t := queue[head]
       head++
       t1, t2 := t.f, t.s
       // sort neighbors by reverse distances
       neigh1 := append([]int(nil), v[t1]...)
       neigh2 := append([]int(nil), v[t2]...)
       sort.Slice(neigh1, func(i, j int) bool { return d2[neigh1[i]] < d2[neigh1[j]] })
       sort.Slice(neigh2, func(i, j int) bool { return d1[neigh2[i]] < d1[neigh2[j]] })
       cnt1, cnt2 := 0, 0
       i := 0
       for cnt1 < 4 && i < len(neigh1) && !found {
           j := 0
           for cnt2 < 4 && j < len(neigh2) {
               from, to := neigh1[i], neigh2[j]
               if !used[from][to] && from != to {
                   used[from][to] = true
                   pf[from][to] = t1
                   ps[from][to] = t2
                   d[from][to] = d[t1][t2] + 1
                   if from == n-1 && to == 0 {
                       found = true
                       break
                   }
                   queue = append(queue, pair{from, to})
               }
               j++
               cnt2++
           }
           i++
           cnt1++
       }
   }
   if d[n-1][0] == mod {
       fmt.Println(-1)
       return
   }
   fmt.Println(d[n-1][0])
   // reconstruct paths
   ansf := []int{}
   anss := []int{}
   curf, curs := n-1, 0
   for curf != 0 || curs != 0 {
       ansf = append(ansf, curf+1)
       anss = append(anss, curs+1)
       prevf, prevs := pf[curf][curs], ps[curf][curs]
       curf, curs = prevf, prevs
   }
   for i := len(ansf) - 1; i >= 0; i-- {
       fmt.Printf("%d ", ansf[i])
   }
   fmt.Println()
   for i := len(anss) - 1; i >= 0; i-- {
       fmt.Printf("%d ", anss[i])
   }
   fmt.Println()
}
