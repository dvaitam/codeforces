package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   // fast integer reader
   readInt := func() int {
       var c byte
       var err error
       // skip non-digits
       for {
           c, err = reader.ReadByte()
           if err != nil {
               return 0
           }
           if c >= '0' && c <= '9' {
               break
           }
       }
       x := int(c - '0')
       for {
           c, err = reader.ReadByte()
           if err != nil || c < '0' || c > '9' {
               break
           }
           x = x*10 + int(c-'0')
       }
       return x
   }
   n := readInt()
   s := readInt()
   // adjacency list
   head := make([]int, n+1)
   type Edge struct{ to, val, next int }
   edges := make([]Edge, 2*(n+1))
   edgeCnt := 0
   for i := 0; i < n-1; i++ {
       u := readInt()
       v := readInt()
       w := readInt()
       edgeCnt++
       edges[edgeCnt] = Edge{v, w, head[u]}
       head[u] = edgeCnt
       edgeCnt++
       edges[edgeCnt] = Edge{u, w, head[v]}
       head[v] = edgeCnt
   }
   // distance and parent arrays
   dis := make([]int64, n+1)
   fa := make([]int, n+1)
   // farthest node from a start
   type nodePair struct{ x, p int }
   var farthest func(int) int
   farthest = func(start int) int {
       // reset
       for i := 1; i <= n; i++ {
           dis[i] = 0
           fa[i] = 0
       }
       stack := []nodePair{{start, 0}}
       mx := start
       for len(stack) > 0 {
           cur := stack[len(stack)-1]
           stack = stack[:len(stack)-1]
           x, p := cur.x, cur.p
           if dis[x] > dis[mx] {
               mx = x
           }
           for e := head[x]; e != 0; e = edges[e].next {
               y := edges[e].to
               if y == p {
                   continue
               }
               dis[y] = dis[x] + int64(edges[e].val)
               fa[y] = x
               stack = append(stack, nodePair{y, x})
           }
       }
       return mx
   }
   // find diameter endpoints
   u := farthest(1)
   v := farthest(u)
   // record path from v to u
   path := make([]int, 0, n)
   onTrunk := make([]bool, n+1)
   for x := v; x != 0; x = fa[x] {
       path = append(path, x)
       onTrunk[x] = true
   }
   m := len(path)
   // distances from u
   disU := make([]int64, m)
   for i, x := range path {
       disU[i] = dis[x]
   }
   // max subtree distance for each trunk node
   mxdPath := make([]int64, m)
   type subInfo struct{ x, p int; dist int64 }
   for i, x := range path {
       var branchMax int64
       for e := head[x]; e != 0; e = edges[e].next {
           y := edges[e].to
           if onTrunk[y] {
               continue
           }
           // DFS this subtree
           st := []subInfo{{y, x, int64(edges[e].val)}}
           for len(st) > 0 {
               cur := st[len(st)-1]
               st = st[:len(st)-1]
               if cur.dist > branchMax {
                   branchMax = cur.dist
               }
               for ee := head[cur.x]; ee != 0; ee = edges[ee].next {
                   yy := edges[ee].to
                   if yy == cur.p {
                       continue
                   }
                   st = append(st, subInfo{yy, cur.x, cur.dist + int64(edges[ee].val)})
               }
           }
       }
       mxdPath[i] = branchMax
   }
   // sliding window over path
   INF := int64(4e18)
   best := INF
   dq := make([]int, 0, m)
   r := 0
   // helper
   max := func(a, b int64) int64 {
       if a > b {
           return a
       }
       return b
   }
   for l := 0; l < m; l++ {
       // pop outdated
       if len(dq) > 0 && dq[0] < l {
           dq = dq[1:]
       }
       // extend r
       for r < m && r-l < s {
           for len(dq) > 0 && mxdPath[r] >= mxdPath[dq[len(dq)-1]] {
               dq = dq[:len(dq)-1]
           }
           dq = append(dq, r)
           r++
       }
       // compute candidate
       // distance to v = disU[0] - disU[l]
       // distance to u = disU[r-1]
       cand := max(max(disU[0]-disU[l], disU[r-1]), mxdPath[dq[0]])
       if cand < best {
           best = cand
       }
   }
   fmt.Println(best)
}
