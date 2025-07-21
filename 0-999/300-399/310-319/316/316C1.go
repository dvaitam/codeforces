package main

import (
   "bufio"
   "fmt"
   "os"
)

var n, m, K int
var pos [][2]int
var edges [][2]int
var costMat [][]int
var cand [][]int
var order []int
var lbSum []int
var best int

func dfs(idx int, used uint64, cur int) {
   if cur+lbSum[idx] >= best {
       return
   }
   if idx == K {
       if cur < best {
           best = cur
       }
       return
   }
   k := order[idx]
   for _, e := range cand[k] {
       c := costMat[k][e]
       if cur+c >= best {
           continue
       }
       u, v := edges[e][0], edges[e][1]
       bit := (uint64(1) << u) | (uint64(1) << v)
       if used&bit != 0 {
           continue
       }
       dfs(idx+1, used|bit, cur+c)
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fscan(in, &n, &m)
   N := n * m
   K = N / 2
   pos = make([][2]int, K)
   cnt := make([]int, K)
   grid := make([]int, N)
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           var x int
           fmt.Fscan(in, &x)
           x--
           idx := i*m + j
           grid[idx] = x
           if cnt[x] < 2 {
               pos[x][cnt[x]] = idx
           }
           cnt[x]++
       }
   }
   // build edges
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           u := i*m + j
           if j+1 < m {
               edges = append(edges, [2]int{u, u + 1})
           }
           if i+1 < n {
               edges = append(edges, [2]int{u, u + m})
           }
       }
   }
   E := len(edges)
   costMat = make([][]int, K)
   cand = make([][]int, K)
   for k := 0; k < K; k++ {
       costMat[k] = make([]int, E)
       for e := 0; e < E; e++ {
           u, v := edges[e][0], edges[e][1]
           p0, p1 := pos[k][0], pos[k][1]
           c1 := 0
           if p0 != u {
               c1++
           }
           if p1 != v {
               c1++
           }
           c2 := 0
           if p0 != v {
               c2++
           }
           if p1 != u {
               c2++
           }
           costMat[k][e] = c1
           if c2 < c1 {
               costMat[k][e] = c2
           }
       }
       // create candidate edge list sorted by cost
       idxs := make([]int, E)
       for e := 0; e < E; e++ {
           idxs[e] = e
       }
       // simple insertion sort by cost
       for i := 1; i < E; i++ {
           v := idxs[i]
           j := i
           for j > 0 && costMat[k][idxs[j-1]] > costMat[k][v] {
               idxs[j] = idxs[j-1]
               j--
           }
           idxs[j] = v
       }
       cand[k] = idxs
   }
   // order IDs
   order = make([]int, K)
   for i := 0; i < K; i++ {
       order[i] = i
   }
   // sort order by number of zero-cost candidates then one-cost
   // compute zero counts
   type od struct{ k, c0, c1 int }
   ods := make([]od, K)
   for _, k := range order {
       c0, c1 := 0, 0
       for _, e := range cand[k] {
           if costMat[k][e] == 0 {
               c0++
           } else if costMat[k][e] == 1 {
               c1++
           } else {
               break
           }
       }
       ods[k] = od{k, c0, c1}
   }
   // simple sort
   for i := 1; i < K; i++ {
       tmp := ods[i]
       j := i
       for j > 0 {
           o := ods[j-1]
           if tmp.c0 < o.c0 || (tmp.c0 == o.c0 && tmp.c1 < o.c1) {
               ods[j] = o
               j--
           } else {
               break
           }
       }
       ods[j] = tmp
   }
   for i := 0; i < K; i++ {
       order[i] = ods[i].k
   }
   // lower bound sums
   lbSum = make([]int, K+1)
   lb := make([]int, K)
   for i := 0; i < K; i++ {
       k := order[i]
       // minimal cost edge
       mc := costMat[k][cand[k][0]]
       lb[i] = mc
   }
   lbSum[K] = 0
   for i := K - 1; i >= 0; i-- {
       lbSum[i] = lbSum[i+1] + lb[i]
   }
   best = 2 * K
   dfs(0, 0, 0)
   fmt.Fprintln(out, best)
}
