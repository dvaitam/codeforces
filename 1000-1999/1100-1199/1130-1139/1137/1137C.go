package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, m, d int
   fmt.Fscan(reader, &n, &m, &d)
   g := make([][]int, n+1)
   o := make([][]int, n+1)
   for i := 0; i < m; i++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
       g[x] = append(g[x], y)
       o[y] = append(o[y], x)
   }
   s := make([][]byte, n+1)
   for i := 1; i <= n; i++ {
       var str string
       fmt.Fscan(reader, &str)
       s[i] = []byte(str)
   }
   // first dfs to get order
   u := make([]bool, n+1)
   order := make([]int, 0, n)
   type frame struct{ v, idx int }
   for i := 1; i <= n; i++ {
       if !u[i] {
           stack := []frame{{i, 0}}
           u[i] = true
           for len(stack) > 0 {
               fr := &stack[len(stack)-1]
               if fr.idx < len(g[fr.v]) {
                   to := g[fr.v][fr.idx]
                   fr.idx++
                   if !u[to] {
                       u[to] = true
                       stack = append(stack, frame{to, 0})
                   }
               } else {
                   order = append(order, fr.v)
                   stack = stack[:len(stack)-1]
               }
           }
       }
   }
   // components
   c := make([]int, n+1)
   dst := make([]int, n+1)
   var cc int
   vcomp := make([][]int, 1) // 1-indexed
   dd := []int{0}
   w := make([][]bool, n+1)
   for i := 1; i <= n; i++ {
       w[i] = make([]bool, d)
   }
   for idx := len(order) - 1; idx >= 0; idx-- {
       u0 := order[idx]
       if c[u0] != 0 {
           continue
       }
       cc++
       // new component
       vcomp = append(vcomp, nil)
       dd = append(dd, 0)
       // dfs on reverse graph to assign component and dst
       stack := []struct{ v, dist int }{{u0, 0}}
       c[u0] = cc
       dst[u0] = 0
       vcomp[cc] = append(vcomp[cc], u0)
       for len(stack) > 0 {
           el := stack[len(stack)-1]
           stack = stack[:len(stack)-1]
           x := el.v
           d0 := el.dist
           for _, y := range o[x] {
               if c[y] == 0 {
                   c[y] = cc
                   dst[y] = d0 + 1
                   vcomp[cc] = append(vcomp[cc], y)
                   stack = append(stack, struct{ v, dist int }{y, d0 + 1})
               }
           }
       }
       // bfs on time states to find cycle length
       qx := []int{u0}
       qy := []int{0}
       w[u0][0] = true
       for i := 0; i < len(qx); i++ {
           x := qx[i]
           z := (qy[i] + 1) % d
           for _, y := range g[x] {
               if c[y] == cc && !w[y][z] {
                   w[y][z] = true
                   qx = append(qx, y)
                   qy = append(qy, z)
               }
           }
       }
       // find minimal cycle length
       for k := 1; k < d; k++ {
           if w[u0][k] {
               dd[cc] = k
               break
           }
       }
       if dd[cc] == 0 {
           dd[cc] = d
       }
   }
   // build component graph edges and reverse graph
   compG := make([][]int, cc+1)
   outdeg := make([]int, cc+1)
   revG := make([][]int, cc+1)
   for x := 1; x <= n; x++ {
       cx := c[x]
       for _, y := range g[x] {
           cy := c[y]
           if cy != cx {
               compG[cx] = append(compG[cx], cy)
           }
       }
   }
   // dedupe and build reverse
   for i := 1; i <= cc; i++ {
       lst := compG[i]
       if len(lst) > 1 {
           // unique
           mapi := make(map[int]struct{}, len(lst))
           uniq := lst[:0]
           for _, y := range lst {
               if _, ok := mapi[y]; !ok {
                   mapi[y] = struct{}{}
                   uniq = append(uniq, y)
               }
           }
           compG[i] = uniq
       }
       outdeg[i] = len(compG[i])
       for _, j := range compG[i] {
           revG[j] = append(revG[j], i)
       }
   }
   // Kahn on reverse graph to get DP order
   queue := make([]int, 0, cc)
   for i := 1; i <= cc; i++ {
       if outdeg[i] == 0 {
           queue = append(queue, i)
       }
   }
   dpOrder := make([]int, 0, cc)
   for i := 0; i < len(queue); i++ {
       u := queue[i]
       dpOrder = append(dpOrder, u)
       for _, v := range revG[u] {
           outdeg[v]--
           if outdeg[v] == 0 {
               queue = append(queue, v)
           }
       }
   }
   // DP f for each component
   f := make([][]int, cc+1)
   for _, comp := range dpOrder {
       // initialize
       f[comp] = make([]int, dd[comp])
       // merge from successors
       for _, x := range vcomp[comp] {
           for _, y := range g[x] {
               c2 := c[y]
               if c2 == comp {
                   continue
               }
               for j := 0; j < dd[comp]; j++ {
                   ja := (j + dst[x]) % dd[comp]
                   for k := j; k < d; k += dd[comp] {
                       idx2 := (k + 1 + dst[y]) % dd[c2]
                       if f[c2][idx2] > f[comp][ja] {
                           f[comp][ja] = f[c2][idx2]
                       }
                   }
               }
           }
       }
       // add own ones
       for _, x := range vcomp[comp] {
           for j := 0; j < dd[comp]; j++ {
               ja := (j + dst[x]) % dd[comp]
               for k := j; k < d; k += dd[comp] {
                   if s[x][k] == '1' {
                       f[comp][ja]++
                       break
                   }
               }
           }
       }
   }
   // result
   c1 := c[1]
   res := f[c1][dst[1]%dd[c1]]
   fmt.Fprintln(writer, res)
}
