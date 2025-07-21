package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// Rolling hash parameters
const (
   mod1 = 1000000007
   mod2 = 1000000009
   base = 91138233
)

// Candidate for square: diff x, start b, letter key k
type Candidate struct {
   x, b int
   k    int
}
// Min-heap of Candidate by x, then b
type CHeap []Candidate
func (h CHeap) Len() int           { return len(h) }
func (h CHeap) Less(i, j int) bool {
   if h[i].x != h[j].x {
       return h[i].x < h[j].x
   }
   return h[i].b < h[j].b
}
func (h CHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *CHeap) Push(x interface{}) { *h = append(*h, x.(Candidate)) }
func (h *CHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[:n-1]
   return x
}

// Letter group: positions and current index
type Group struct {
   pos []int
   i   int // index in pos for candidate
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   s := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &s[i])
   }
   // build groups
   groups := make(map[int]*Group, n)
   for i, v := range s {
       g, ok := groups[v]
       if !ok {
           g = &Group{pos: make([]int, 0, 4)}
           groups[v] = g
       }
       g.pos = append(g.pos, i)
   }
   // precompute hashes and powers
   H1 := make([]int, n+1)
   H2 := make([]int, n+1)
   P1 := make([]int, n+1)
   P2 := make([]int, n+1)
   P1[0], P2[0] = 1, 1
   for i := 0; i < n; i++ {
       P1[i+1] = int((int64(P1[i]) * base) % mod1)
       P2[i+1] = int((int64(P2[i]) * base) % mod2)
       x := s[i]
       H1[i+1] = int(( (int64(H1[i])*base)%mod1 + int64(x+1) ) % mod1)
       H2[i+1] = int(( (int64(H2[i])*base)%mod2 + int64(x+1) ) % mod2)
   }
   // function to get hash of [l..r]
   getHash := func(l, r int) (int, int) {
       a1 := H1[r+1] - int(int64(H1[l])*int64(P1[r-l+1])%mod1)
       if a1 < 0 {
           a1 += mod1
       }
       a2 := H2[r+1] - int(int64(H2[l])*int64(P2[r-l+1])%mod2)
       if a2 < 0 {
           a2 += mod2
       }
       return a1, a2
   }
   // initialize heap
   h := &CHeap{}
   heap.Init(h)
   // current offset
   d := 0
   // helper to advance group candidate
   pushCandidate := func(k int, g *Group) {
       // find first i such that pos[i] and pos[i+1] >= d
       for g.i < len(g.pos)-1 && (g.pos[g.i] < d || g.pos[g.i+1] < d) {
           g.i++
       }
       if g.i < len(g.pos)-1 {
           x := g.pos[g.i+1] - g.pos[g.i]
           heap.Push(h, Candidate{x: x, b: g.pos[g.i], k: k})
       }
   }
   // initial fill
   for k, g := range groups {
       if len(g.pos) >= 2 {
           pushCandidate(k, g)
       }
   }
   // process
   for h.Len() > 0 {
       can := heap.Pop(h).(Candidate)
       g := groups[can.k]
       // check stale or out of date
       if g.i >= len(g.pos)-1 || g.pos[g.i] != can.b || (g.pos[g.i+1]-g.pos[g.i]) != can.x || can.b < d {
           // stale, try next for this group
           pushCandidate(can.k, g)
           continue
       }
       x, b := can.x, can.b
       if b+2*x > n {
           // too long, no more for this group
           g.i++
           pushCandidate(can.k, g)
           continue
       }
       // check hash equality for the two halves
       ha1, ha2 := getHash(b, b+x-1)
       hb1, hb2 := getHash(b+x, b+2*x-1)
       if ha1 != hb1 || ha2 != hb2 {
           // not a square, skip this candidate
           g.i++
           pushCandidate(can.k, g)
           continue
       }
       // found shortest-leftmost repeat: delete prefix up to b+x-1
       d = b + x
       // continue with new d
   }
   // output remaining from d to n-1
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   rem := n - d
   // print result length and remaining elements (possibly empty line)
   fmt.Fprintln(out, rem)
   if rem == 0 {
       fmt.Fprintln(out)
   } else {
       for i := d; i < n; i++ {
           fmt.Fprint(out, s[i])
           if i+1 < n {
               fmt.Fprint(out, " ")
           }
       }
       fmt.Fprintln(out)
   }
}
