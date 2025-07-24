package main

import (
   "bufio"
   "fmt"
   "os"
)

type Op struct {
   t, i, j int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m, q int
   fmt.Fscan(reader, &n, &m, &q)
   ops := make([]Op, q+1)
   parent := make([]int, q+1)
   children := make([][]int, q+1)
   last := 0
   for id := 1; id <= q; id++ {
       var t int
       fmt.Fscan(reader, &t)
       ops[id].t = t
       if t == 1 || t == 2 {
           fmt.Fscan(reader, &ops[id].i, &ops[id].j)
           parent[id] = last
       } else if t == 3 {
           fmt.Fscan(reader, &ops[id].i)
           parent[id] = last
       } else if t == 4 {
           var k int
           fmt.Fscan(reader, &k)
           parent[id] = k
       }
       children[parent[id]] = append(children[parent[id]], id)
       last = id
   }
   // bitsets
   W := (m + 63) >> 6
   bits := make([][]uint64, n+1)
   for i := 1; i <= n; i++ {
       bits[i] = make([]uint64, W)
   }
   lastMask := uint64(0)
   if m%64 != 0 {
       lastMask = (uint64(1) << uint(m%64)) - 1
   } else {
       lastMask = ^uint64(0)
   }
   cnt := make([]int, n+1)
   ans := make([]int, q+1)
   tot := 0
   // iterative DFS
   type Frame struct { id, next int; changed bool }
   stack := make([]Frame, 0, q+1)
   stack = append(stack, Frame{0, 0, false})
   for len(stack) > 0 {
       f := &stack[len(stack)-1]
       id := f.id
       if f.next == 0 {
           // apply op for id
           if id != 0 {
               op := ops[id]
               switch op.t {
               case 1:
                   i, j := op.i, op.j-1
                   w, b := j>>6, uint(j&63)
                   if ((bits[i][w]>>b)&1) == 0 {
                       bits[i][w] |= 1 << b
                       cnt[i]++
                       tot++
                       f.changed = true
                   }
               case 2:
                   i, j := op.i, op.j-1
                   w, b := j>>6, uint(j&63)
                   if ((bits[i][w]>>b)&1) == 1 {
                       bits[i][w] &^= 1 << b
                       cnt[i]--
                       tot--
                       f.changed = true
                   }
               case 3:
                   i := op.i
                   tot -= cnt[i]
                   cnt[i] = m - cnt[i]
                   tot += cnt[i]
                   for w := 0; w < W; w++ {
                       bits[i][w] = ^bits[i][w]
                   }
                   // mask last
                   bits[i][W-1] &= lastMask
               }
           }
           ans[id] = tot
       }
       if f.next < len(children[id]) {
           child := children[id][f.next]
           f.next++
           stack = append(stack, Frame{child, 0, false})
           continue
       }
       // done children, undo and pop
       if id != 0 {
           op := ops[id]
           switch op.t {
           case 1:
               if f.changed {
                   i, j := op.i, op.j-1
                   w, b := j>>6, uint(j&63)
                   bits[i][w] &^= 1 << b
                   cnt[i]--
                   tot--
               }
           case 2:
               if f.changed {
                   i, j := op.i, op.j-1
                   w, b := j>>6, uint(j&63)
                   bits[i][w] |= 1 << b
                   cnt[i]++
                   tot++
               }
           case 3:
               i := op.i
               tot -= cnt[i]
               cnt[i] = m - cnt[i]
               tot += cnt[i]
               for w := 0; w < W; w++ {
                   bits[i][w] = ^bits[i][w]
               }
               bits[i][W-1] &= lastMask
           }
       }
       stack = stack[:len(stack)-1]
   }
   // output answers for ops 1..q
   for id := 1; id <= q; id++ {
       fmt.Fprintln(writer, ans[id])
   }
}
