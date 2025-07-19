package main

import (
   "bufio"
   "fmt"
   "os"
)

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

var (
   n       int
   adj     [][]int
   parent1 []int
   siz     []int
   reader  = bufio.NewReader(os.Stdin)
)

func readInt() int {
   var x int
   fmt.Fscan(reader, &x)
   return x
}

func dfsSize(u, p int) {
   parent1[u] = p
   siz[u] = 1
   for _, v := range adj[u] {
       if v == p {
           continue
       }
       dfsSize(v, u)
       siz[u] += siz[v]
   }
}

func main() {
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   n = readInt()
   adj = make([][]int, n+1)
   parent1 = make([]int, n+1)
   siz = make([]int, n+1)
   for i := 0; i < n-1; i++ {
       u := readInt()
       v := readInt()
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   dfsSize(1, 0)
   // find centroid
   rt, best := 1, n+1
   for i := 1; i <= n; i++ {
       mx := n - siz[i]
       for _, v := range adj[i] {
           if v != parent1[i] {
               mx = max(mx, siz[v])
           }
       }
       if mx < best {
           best = mx
           rt = i
       }
   }
   // collect components
   var child []int
   comp := make([]int, 0, len(adj[rt]))
   for _, v := range adj[rt] {
       child = append(child, v)
       if v == parent1[rt] {
           comp = append(comp, n-siz[rt])
       } else {
           comp = append(comp, siz[v])
       }
   }
   tn := len(child)
   // dp
   f := make([][]bool, tn+1)
   las := make([][]bool, tn+1)
   for i := 0; i <= tn; i++ {
       f[i] = make([]bool, n+1)
       las[i] = make([]bool, n+1)
   }
   f[0][0] = true
   for i := 1; i <= tn; i++ {
       w := comp[i-1]
       for j := 0; j <= n; j++ {
           if f[i-1][j] {
               f[i][j] = true
               las[i][j] = false
               if j+w <= n {
                   f[i][j+w] = true
                   las[i][j+w] = true
               }
           }
       }
   }
   // choose partition
   nw := 0
   for j := 0; j <= n; j++ {
       if f[tn][j] && (j+1)*(n-j) > 2*n*n/9 {
           nw = j
           break
       }
   }
   // assign values
   val := make([]int, n+1)
   for i := 1; i <= nw; i++ {
       val[i] = i
   }
   for i := nw + 1; i < n; i++ {
       val[i] = (i-nw)*(nw+1)
   }
   col := make(map[int]bool, tn)
   // backtrack
   rem := nw
   for i := tn; i >= 1; i-- {
       v := child[i-1]
       if las[i][rem] {
           col[v] = true
           rem -= comp[i-1]
       } else {
           col[v] = false
       }
   }
   // assign order values
   parent := make([]int, n+1)
   vval := make([]int, n+1)
   cur := 1
   // dfs assign
   var assign func(u int)
   assign = func(u int) {
       vval[u] = val[cur]
       cur++
       for _, w := range adj[u] {
           if w == parent[u] {
               continue
           }
           parent[w] = u
           assign(w)
       }
   }
   // start with centroid
   parent[rt] = 0
   vval[rt] = 0
   // group 1 then 0
   for _, u := range child {
       if col[u] {
           parent[u] = rt
           assign(u)
       }
   }
   for _, u := range child {
       if !col[u] {
           parent[u] = rt
           assign(u)
       }
   }
   // output
   for i := 1; i <= n; i++ {
       if i == rt {
           continue
       }
       fmt.Fprintf(out, "%d %d %d\n", i, parent[i], vval[i]-vval[parent[i]])
   }
}
