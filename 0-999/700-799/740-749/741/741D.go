package main

import (
   "bufio"
   "fmt"
   "os"
)

const BITS = 22
const INF = 1000000000

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n int
   fmt.Fscan(reader, &n)
   children := make([][]struct{to, c int}, n+1)
   for i := 2; i <= n; i++ {
       var p int
       var cc byte
       fmt.Fscan(reader, &p, &cc)
       children[p] = append(children[p], struct{to, c int}{i, int(cc - 'a')})
   }
   depth := make([]int, n+1)
   mask := make([]int, n+1)
   tin := make([]int, n+1)
   tout := make([]int, n+1)
   ver := make([]int, n)
   sz := make([]int, n+1)
   heavy := make([]int, n+1)
   // DFS to compute tin, tout, depth, mask, sz, heavy
   type nodeState struct{ v, ci int }
   timer := 0
   stack := []nodeState{{1, 0}}
   for len(stack) > 0 {
       top := &stack[len(stack)-1]
       v := top.v
       if top.ci == 0 {
           tin[v] = timer
           ver[timer] = v
           timer++
           sz[v] = 1
       }
       if top.ci < len(children[v]) {
           ch := children[v][top.ci]
           top.ci++
           u := ch.to
           depth[u] = depth[v] + 1
           mask[u] = mask[v] ^ (1 << ch.c)
           stack = append(stack, nodeState{u, 0})
       } else {
           // exiting v
           // compute heavy
           maxsz := 0
           for _, ch := range children[v] {
               u := ch.to
               if sz[u] > maxsz {
                   maxsz = sz[u]
                   heavy[v] = u
               }
               sz[v] += sz[u]
           }
           tout[v] = timer
           stack = stack[:len(stack)-1]
       }
   }
   // prepare mp
   M := 1 << BITS
   mp := make([]int, M)
   for i := 0; i < M; i++ {
       mp[i] = -INF
   }
   ans := make([]int, n+1)
   for i := 1; i <= n; i++ {
       ans[i] = 1
   }
   // helper functions
   query := func(x, root int) {
       mx := mp[mask[x]]
       for b := 0; b < BITS; b++ {
           v := mask[x] ^ (1 << b)
           if mp[v] > mx {
               mx = mp[v]
           }
       }
       if mx > -INF/2 {
           // path length in nodes
           length := depth[x] + mx - 2*depth[root] + 1
           if length > ans[root] {
               ans[root] = length
           }
       }
   }
   add := func(x int) {
       m := mask[x]
       if depth[x] > mp[m] {
           mp[m] = depth[x]
       }
   }
   // iterative DSU on tree
   type frame struct{ v, phase, idx int; keep bool }
   st := make([]frame, 0, n*2)
   st = append(st, frame{1, 0, 0, true})
   for len(st) > 0 {
       f := st[len(st)-1]
       st = st[:len(st)-1]
       v := f.v
       switch f.phase {
       case 0:
           // process small children
           f.phase = 1; f.idx = 0
           st = append(st, f)
       case 1:
           // handle each small child
           for f.idx < len(children[v]) {
               u := children[v][f.idx].to
               f.idx++
               if u == heavy[v] {
                   continue
               }
               // recurse on u
               st = append(st, f)
               st = append(st, frame{u, 0, 0, false})
               goto Next
           }
           // done small
           f.phase = 2; st = append(st, f)
       case 2:
           // heavy child
           if heavy[v] != 0 {
               // after heavy, go to merge
               st = append(st, frame{v, 3, 0, f.keep})
               st = append(st, frame{heavy[v], 0, 0, true})
           } else {
               f.phase = 3; st = append(st, f)
           }
       case 3:
           // merge small children into mp
           for _, ch := range children[v] {
               u := ch.to
               if u == heavy[v] {
                   continue
               }
               for i := tin[u]; i < tout[u]; i++ {
                   x := ver[i]
                   query(x, v)
                   add(x)
               }
           }
           f.phase = 4; st = append(st, f)
       case 4:
           // include v itself
           query(v, v)
           add(v)
           f.phase = 5; st = append(st, f)
       case 5:
           // cleanup if needed
           if !f.keep {
               for i := tin[v]; i < tout[v]; i++ {
                   mp[mask[ver[i]]] = -INF
               }
           }
       }
   Next:
       continue
   }
   // output
   for i := 1; i <= n; i++ {
       fmt.Fprint(writer, ans[i])
       if i < n {
           fmt.Fprint(writer, ' ')
       }
   }
   fmt.Fprintln(writer)
}
