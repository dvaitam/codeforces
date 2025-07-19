package main

import (
   "bufio"
   "fmt"
   "os"
)

// Edge holds a neighbor and edge weight
type Edge struct {
   to, w int
}

// Frame for iterative DFS
type Frame struct {
   u, parent, idx, wUp int
}

func main() {
   rdr := bufio.NewReader(os.Stdin)
   defer rdr.Reset(nil)
   readInt := func() int {
       var x int
       var c byte
       var neg bool
       for {
           b, err := rdr.ReadByte()
           if err != nil {
               return 0
           }
           c = b
           if (c >= '0' && c <= '9') || c == '-' {
               break
           }
       }
       if c == '-' {
           neg = true
       } else {
           x = int(c - '0')
       }
       for {
           b, err := rdr.ReadByte()
           if err != nil {
               break
           }
           c = b
           if c < '0' || c > '9' {
               break
           }
           x = x*10 + int(c-'0')
       }
       if neg {
           return -x
       }
       return x
   }

   n := readInt()
   val := make([]int, n+1)
   f := make([]int64, n+1)
   g := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       val[i] = readInt()
       f[i] = int64(val[i])
       g[i] = int64(val[i])
   }
   adj := make([][]Edge, n+1)
   for i := 2; i <= n; i++ {
       u := readInt()
       v := readInt()
       w := readInt()
       adj[u] = append(adj[u], Edge{to: v, w: w})
       adj[v] = append(adj[v], Edge{to: u, w: w})
   }
   // iterative DFS
   stack := make([]Frame, 0, n)
   stack = append(stack, Frame{u: 1, parent: 0, idx: 0, wUp: 0})
   for len(stack) > 0 {
       top := &stack[len(stack)-1]
       u := top.u
       if top.idx < len(adj[u]) {
           e := adj[u][top.idx]
           top.idx++
           v := e.to
           w := e.w
           if v == top.parent {
               continue
           }
           // propagate g downward
           if g[u] >= int64(w) {
               cand := g[u] - int64(w) + int64(val[v])
               if cand > g[v] {
                   g[v] = cand
               }
           }
           // go deeper
           stack = append(stack, Frame{u: v, parent: u, idx: 0, wUp: w})
       } else {
           // post-order
           cur := stack[len(stack)-1]
           stack = stack[:len(stack)-1]
           if cur.parent != 0 {
               // update f and g on parent
               if f[cur.u] >= int64(cur.wUp) {
                   cand := f[cur.u] - int64(cur.wUp) + int64(val[cur.parent])
                   if cand > f[cur.parent] {
                       f[cur.parent] = cand
                   }
               }
               if f[cur.parent] > g[cur.parent] {
                   g[cur.parent] = f[cur.parent]
               }
           }
       }
   }
   var ans int64
   for i := 1; i <= n; i++ {
       if g[i] > ans {
           ans = g[i]
       }
   }
   fmt.Println(ans)
}
