package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

type weight struct {
   rCount int
   length int
}

// compare returns true if w is better (less) than o
func (w weight) less(o weight) bool {
   if w.rCount != o.rCount {
       return w.rCount < o.rCount
   }
   return w.length < o.length
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var m int
   if _, err := fmt.Fscan(reader, &m); err != nil {
       return
   }
   essay := make([]string, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &essay[i])
       essay[i] = strings.ToLower(essay[i])
   }
   var n int
   fmt.Fscan(reader, &n)
   // build mapping
   id := make(map[string]int, m+2*n)
   var words []string
   addWord := func(s string) int {
       if v, ok := id[s]; ok {
           return v
       }
       idx := len(words)
       id[s] = idx
       words = append(words, s)
       return idx
   }
   // add essay words
   for _, w := range essay {
       addWord(w)
   }
   edges := make([][]int, 0)
   // placeholder, will resize after mapping
   var edgePairs [][2]int
   edgePairs = make([][2]int, n)
   // read synonyms
   for i := 0; i < n; i++ {
       var x, y string
       fmt.Fscan(reader, &x, &y)
       x = strings.ToLower(x)
       y = strings.ToLower(y)
       u := addWord(x)
       v := addWord(y)
       edgePairs[i] = [2]int{u, v}
   }
   N := len(words)
   edges = make([][]int, N)
   rev := make([][]int, N)
   for _, p := range edgePairs {
       u, v := p[0], p[1]
       edges[u] = append(edges[u], v)
       rev[v] = append(rev[v], u)
   }
   // compute weights for nodes
   nodeW := make([]weight, N)
   for i, s := range words {
       cnt := 0
       for _, ch := range s {
           if ch == 'r' {
               cnt++
           }
       }
       nodeW[i] = weight{rCount: cnt, length: len(s)}
   }
   // first pass DFS order (iterative)
   visited := make([]bool, N)
   order := make([]int, 0, N)
   type frame struct{ v, idx int; pre bool }
   for v := 0; v < N; v++ {
       if visited[v] {
           continue
       }
       stack := []frame{{v, 0, false}}
       for len(stack) > 0 {
           fr := &stack[len(stack)-1]
           if !fr.pre {
               visited[fr.v] = true
               fr.pre = true
           }
           if fr.idx < len(edges[fr.v]) {
               u := edges[fr.v][fr.idx]
               fr.idx++
               if !visited[u] {
                   stack = append(stack, frame{u, 0, false})
               }
           } else {
               order = append(order, fr.v)
               stack = stack[:len(stack)-1]
           }
       }
   }
   // second pass assign components
   comp := make([]int, N)
   for i := range comp {
       comp[i] = -1
   }
   cid := 0
   for i := len(order) - 1; i >= 0; i-- {
       v := order[i]
       if comp[v] != -1 {
           continue
       }
       // DFS on rev
       stk := []int{v}
       for len(stk) > 0 {
           u := stk[len(stk)-1]
           stk = stk[:len(stk)-1]
           if comp[u] != -1 {
               continue
           }
           comp[u] = cid
           for _, w := range rev[u] {
               if comp[w] == -1 {
                   stk = append(stk, w)
               }
           }
       }
       cid++
   }
   C := cid
   // comp weights
   inf := 1_000_000_000
   compW := make([]weight, C)
   for i := 0; i < C; i++ {
       compW[i] = weight{rCount: inf, length: inf}
   }
   for i := 0; i < N; i++ {
       c := comp[i]
       if nodeW[i].less(compW[c]) {
           compW[c] = nodeW[i]
       }
   }
   // build condensed graph
   cdg := make([][]int, C)
   indeg := make([]int, C)
   for u := 0; u < N; u++ {
       cu := comp[u]
       for _, v := range edges[u] {
           cv := comp[v]
           if cu != cv {
               cdg[cu] = append(cdg[cu], cv)
           }
       }
   }
   // remove duplicate edges? optional, but indeg count over duplicates is fine
   for u := 0; u < C; u++ {
       for _, v := range cdg[u] {
           indeg[v]++
       }
   }
   // topo sort
   q := make([]int, 0, C)
   for i := 0; i < C; i++ {
       if indeg[i] == 0 {
           q = append(q, i)
       }
   }
   topo := make([]int, 0, C)
   for i := 0; i < len(q); i++ {
       u := q[i]
       topo = append(topo, u)
       for _, v := range cdg[u] {
           indeg[v]--
           if indeg[v] == 0 {
               q = append(q, v)
           }
       }
   }
   // dp on reversed topo
   dp := make([]weight, C)
   copy(dp, compW)
   for i := len(topo) - 1; i >= 0; i-- {
       u := topo[i]
       for _, v := range cdg[u] {
           if dp[v].less(dp[u]) {
               dp[u] = dp[v]
           }
       }
   }
   // sum result
   totalR, totalLen := 0, 0
   for _, w := range essay {
       c := comp[id[w]]
       totalR += dp[c].rCount
       totalLen += dp[c].length
   }
   fmt.Println(totalR, totalLen)
}
