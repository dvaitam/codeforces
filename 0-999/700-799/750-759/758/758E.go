package main

import (
   "bufio"
   "fmt"
   "os"
)

const INF = 100000000000000000 // 1e17

// Edge represents a directed tree edge with weights w and p, and d is auxiliary
type Edge struct {
   v1, v2 int
   w, p   int
   d      int64
}

var (
   n  int
   ed []Edge
   g  [][]int
   s  []int64
   ok bool
)

func min64(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func dfs(v int) {
   for _, i := range g[v] {
       e := &ed[i]
       // initial transfer up to parent
       e.d = min64(int64(e.w-1), int64(e.p))
       e.w -= int(e.d)
       e.p -= int(e.d)
       dfs(e.v2)
       // ensure child requirement
       if s[e.v2] > int64(e.p) {
           t := min64(s[e.v2]-int64(e.p), e.d)
           e.p += int(t)
           e.w += int(t)
           e.d -= t
           if s[e.v2] > int64(e.p) {
               ok = false
           }
       }
       s[v] += s[e.v2] + int64(e.w)
   }
}

// goRedistribute pushes remaining capacity down the tree
func goRedistribute(v int, add int64) int64 {
   var res int64
   for _, i := range g[v] {
       e := &ed[i]
       // use available d to increase w and p
       t := min64(add, e.d)
       e.w += int(t)
       e.p += int(t)
       e.d -= t
       add -= t
       res += t
       // then propagate further, constrained by p - s[child]
       need := int64(e.p) - s[e.v2]
       if need < 0 {
           need = 0
       }
       tt := goRedistribute(e.v2, min64(add, need))
       add -= tt
       res += tt
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   fmt.Fscan(reader, &n)
   ed = make([]Edge, n-1)
   g = make([][]int, n)
   s = make([]int64, n)
   for i := 0; i < n-1; i++ {
       var v1, v2, w, p int
       fmt.Fscan(reader, &v1, &v2, &w, &p)
       v1--
       v2--
       ed[i] = Edge{v1: v1, v2: v2, w: w, p: p, d: 0}
       g[v1] = append(g[v1], i)
   }
   ok = true
   dfs(0)
   if !ok {
       fmt.Fprintln(writer, -1)
       return
   }
   goRedistribute(0, INF)
   // output result
   fmt.Fprintln(writer, n)
   for i := 0; i < n-1; i++ {
       e := ed[i]
       // convert back to 1-based
       fmt.Fprintf(writer, "%d %d %d %d\n", e.v1+1, e.v2+1, e.w, e.p)
   }
}
