package main

import (
   "bufio"
   "fmt"
   "os"
)

// DSUNext supports finding next unvisited element >= x
type DSUNext struct {
   parent []int
   n      int
}

func NewDSUNext(n int) *DSUNext {
   p := make([]int, n+2)
   for i := 0; i <= n+1; i++ {
       p[i] = i
   }
   return &DSUNext{parent: p, n: n}
}

func (d *DSUNext) find(x int) int {
   if x > d.n+1 {
       return d.n + 1
   }
   if d.parent[x] != x {
       d.parent[x] = d.find(d.parent[x])
   }
   return d.parent[x]
}

// remove x, so next find(x) will return first >x
func (d *DSUNext) remove(x int) {
   d.parent[x] = d.find(x + 1)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var N, M int
   var x, y int64
   if _, err := fmt.Fscan(in, &N, &x, &M, &y); err != nil {
       return
   }
   d := x - y
   if d < 0 {
       d = -d
   }
   // compute intervals for set1 (i:1..N) over j:1..M
   Li := make([]int, N+1)
   Ri := make([]int, N+1)
   cntI := 0
   var P int64
   for i := 1; i <= N; i++ {
       // L = max(1, |i-d|+1), R = min(M, i+d-1)
       dd := int64(i) - d
       if dd < 0 {
           dd = -dd
       }
       l := int(dd) + 1
       if l < 1 {
           l = 1
       }
       r := int64(i) + d - 1
       if r > int64(M) {
           r = int64(M)
       }
       if l <= int(r) {
           Li[i] = l
           Ri[i] = int(r)
           cntI++
           P += int64(r - int64(l) + 1)
       } else {
           Li[i] = 1
           Ri[i] = 0
       }
   }
   // intervals for set2 (j:1..M) over i:1..N
   Lj := make([]int, M+1)
   Rj := make([]int, M+1)
   cntJ := 0
   for j := 1; j <= M; j++ {
       dd := int64(j) - d
       if dd < 0 {
           dd = -dd
       }
       l := int(dd) + 1
       if l < 1 {
           l = 1
       }
       r := int64(j) + d - 1
       if r > int64(N) {
           r = int64(N)
       }
       if l <= int(r) {
           Lj[j] = l
           Rj[j] = int(r)
           cntJ++
       } else {
           Lj[j] = 1
           Rj[j] = 0
       }
   }
   // prepare DSUNext for BFS
   dsuI := NewDSUNext(N)
   dsuJ := NewDSUNext(M)
   // remove i with no neighbors
   for i := 1; i <= N; i++ {
       if Ri[i] < Li[i] {
           dsuI.remove(i)
       }
   }
   for j := 1; j <= M; j++ {
       if Rj[j] < Lj[j] {
           dsuJ.remove(j)
       }
   }
   // BFS to count connected components among intersecting circles
   CCbig := 0
   // queues
   type node struct{ left bool; idx int }
   var queue []node
   for i := 1; i <= N; i++ {
       // find next unvisited i
       ii := dsuI.find(i)
       if ii > N {
           break
       }
       // start new component
       CCbig++
       queue = append(queue, node{true, ii})
       dsuI.remove(ii)
       // BFS
       for qi := 0; qi < len(queue); qi++ {
           nd := queue[qi]
           if nd.left {
               // expand to js in [Li, Ri]
               i0 := nd.idx
               for j := dsuJ.find(Li[i0]); j <= Ri[i0]; j = dsuJ.find(j) {
                   queue = append(queue, node{false, j})
                   dsuJ.remove(j)
               }
           } else {
               j0 := nd.idx
               for i2 := dsuI.find(Lj[j0]); i2 <= Rj[j0]; i2 = dsuI.find(i2) {
                   queue = append(queue, node{true, i2})
                   dsuI.remove(i2)
               }
           }
       }
       queue = queue[:0]
   }
   // Z = isolated nodes
   Z := int64((N - cntI) + (M - cntJ))
   Cc := int64(Z + CCbig)
   // total regions F = 2*P + Cc + 1
   F := 2*P + Cc + 1
   // output
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   fmt.Fprintln(w, F)
}
