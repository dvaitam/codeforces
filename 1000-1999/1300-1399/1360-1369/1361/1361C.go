package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

// Edge represents an adjacency entry: destination val, edge id num, next pointer
type Edge struct { val, num, next int }

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   mas := make([]int, 2*n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &mas[2*i], &mas[2*i+1])
   }
   const maxBits = 20
   const maxMask = 1 << maxBits
   last := make([]int, maxMask)
   deg := make([]int, maxMask)
   used := make([]bool, n)
   // preallocate for up to 2 edges per input edge (bidirectional)
   vnum := make([]Edge, 0, 2*n*2)
   path := make([]int, n+1)
   cur := make([]int, 2*n)
   st := make([][2]int, 2*n+2)
   var kRes int

   // check if graph with mask of bitLen has Eulerian trail; if need, build cur
   check := func(bitLen int, need bool) bool {
       if bitLen > maxBits {
           return false
       }
       m := 0
       u := (1<<bitLen) - 1
       // reset adjacency and degrees
       for i := 0; i <= u; i++ {
           last[i] = -1
           deg[i] = 0
       }
       // build graph edges
       vnum = vnum[:0]
       for i := 0; i < 2*n; i += 2 {
           va := mas[i] & u
           vb := mas[i+1] & u
           // add edge va->vb
           vnum = append(vnum, Edge{vb, i / 2, last[va]})
           last[va] = m
           m++
           // add edge vb->va
           vnum = append(vnum, Edge{va, i / 2, last[vb]})
           last[vb] = m
           m++
           deg[va]++
           deg[vb]++
       }
       // check even degree
       for i := 0; i <= u; i++ {
           if deg[i]&1 != 0 {
               return false
           }
       }
       // prepare for Eulerian traversal
       for i := 0; i < n; i++ {
           used[i] = false
       }
       kRes = 0
       // iterative Euler: stack of (vertex, incoming edge id)
       s := 0
       start := mas[0] & u
       st[s][0], st[s][1] = start, -1
       s++
       for s > 0 {
           v := st[s-1][0]
           w := st[s-1][1]
           if deg[v] == 0 {
               path[kRes] = w
               kRes++
               s--
               continue
           }
           // skip used edges
           for last[v] != -1 && used[vnum[last[v]].num] {
               last[v] = vnum[last[v]].next
           }
           edge := last[v]
           used[vnum[edge].num] = true
           deg[v]--
           to := vnum[edge].val
           deg[to]--
           // push next
           st[s][0], st[s][1] = to, vnum[edge].num
           s++
           last[v] = vnum[edge].next
       }
       kRes--
       if kRes < n {
           return false
       }
       if !need {
           return true
       }
       // reconstruct actual sequence of half-edge indices
       s2 := path[0] * 2
       if mas[s2]&u != mas[0]&u {
           s2 ^= 1
       }
       k2 := 0
       for i := 0; i < n; i++ {
           cur[k2] = s2
           k2++
           cur[k2] = s2 ^ 1
           k2++
           if i == n-1 {
               break
           }
           t := path[i+1]
           if mas[2*t]&u == mas[s2]&u {
               s2 = 2*t
           } else {
               s2 = 2*t + 1
           }
       }
       kRes = k2
       return true
   }

   // binary search maximum bit length
   l, r := 0, maxBits+1
   for r-l > 1 {
       mid := (l + r) >> 1
       if !check(mid, false) {
           r = mid
       } else {
           l = mid
       }
   }
   // build result for optimal l
   check(l, true)
   // output
   writer.WriteString(strconv.Itoa(l))
   writer.WriteByte('\n')
   for i := 0; i < kRes; i++ {
       writer.WriteString(strconv.Itoa(cur[i] + 1))
       if i+1 < kRes {
           writer.WriteByte(' ')
       }
   }
   writer.WriteByte('\n')
}
