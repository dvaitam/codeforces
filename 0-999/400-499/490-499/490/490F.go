package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n int
   r []int
   adj [][]int
   removed []bool
   sz []int
   ans int
)

// compute subtree sizes
func dfsSize(u, p int) {
   sz[u] = 1
   for _, v := range adj[u] {
       if v != p && !removed[v] {
           dfsSize(v, u)
           sz[u] += sz[v]
       }
   }
}

// find centroid
func dfsCentroid(u, p, tot int) int {
   for _, v := range adj[u] {
       if v != p && !removed[v] && sz[v] > tot/2 {
           return dfsCentroid(v, u, tot)
       }
   }
   return u
}

// lower_bound on slice a for x
func lowerBound(a []int, x int) int {
   l, r := 0, len(a)
   for l < r {
       m := (l + r) >> 1
       if a[m] < x {
           l = m + 1
       } else {
           r = m
       }
   }
   return l
}

// process centroid c
func process(c int) {
   // initial globals for this centroid
   bestL, bestR := 1, 1
   // for each subtree
   for _, v0 := range adj[c] {
       if removed[v0] {
           continue
       }
       // best in this subtree
       bestLi, bestRi := 1, 1
       // dL for left (r < r[c]), use -r for LIS >
       dL := []int{-r[c]}
       // dR for right (r > r[c])
       dR := []int{r[c]}
       // DFS left
       var dfsL func(u, p int)
       dfsL = func(u, p int) {
           v := -r[u]
           pos := lowerBound(dL, v)
           old := 0
           if pos < len(dL) {
               old = dL[pos]
               dL[pos] = v
           } else {
               old = 0
               dL = append(dL, v)
           }
           curr := pos + 1
           if curr > bestLi {
               // only count if >1 (i.e., r[u] < r[c])
               if curr > 1 {
                   bestLi = curr
               }
           }
           for _, w := range adj[u] {
               if w != p && !removed[w] {
                   dfsL(w, u)
               }
           }
           // restore
           if old != 0 || pos < len(dL) {
               dL[pos] = old
           }
           if pos == len(dL)-1 && old == 0 {
               dL = dL[:pos]
           }
       }
       // DFS right
       var dfsR func(u, p int)
       dfsR = func(u, p int) {
           v := r[u]
           pos := lowerBound(dR, v)
           old := 0
           if pos < len(dR) {
               old = dR[pos]
               dR[pos] = v
           } else {
               old = 0
               dR = append(dR, v)
           }
           curr := pos + 1
           if curr > bestRi {
               if curr > 1 {
                   bestRi = curr
               }
           }
           for _, w := range adj[u] {
               if w != p && !removed[w] {
                   dfsR(w, u)
               }
           }
           // restore
           if old != 0 || pos < len(dR) {
               dR[pos] = old
           }
           if pos == len(dR)-1 && old == 0 {
               dR = dR[:pos]
           }
       }
       dfsL(v0, c)
       dfsR(v0, c)
       // combine with globals
       if bestL + bestRi - 1 > ans {
           ans = bestL + bestRi - 1
       }
       if bestLi + bestR - 1 > ans {
           ans = bestLi + bestR - 1
       }
       // update globals
       if bestLi > bestL {
           bestL = bestLi
       }
       if bestRi > bestR {
           bestR = bestRi
       }
   }
   // also alone
   if bestL > ans {
       ans = bestL
   }
   if bestR > ans {
       ans = bestR
   }
}

// centroid decomposition
func decompose(u int) {
   dfsSize(u, -1)
   c := dfsCentroid(u, -1, sz[u])
   process(c)
   removed[c] = true
   for _, v := range adj[c] {
       if !removed[v] {
           decompose(v)
       }
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   fmt.Fscan(in, &n)
   r = make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &r[i])
   }
   adj = make([][]int, n)
   for i := 0; i < n-1; i++ {
       var a, b int
       fmt.Fscan(in, &a, &b)
       a--;
       b--;
       adj[a] = append(adj[a], b)
       adj[b] = append(adj[b], a)
   }
   removed = make([]bool, n)
   sz = make([]int, n)
   ans = 1
   decompose(0)
   // output
   fmt.Println(ans)
}
