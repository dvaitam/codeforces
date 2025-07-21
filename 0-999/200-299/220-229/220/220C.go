package main

import (
   "bufio"
   "fmt"
   "os"
)

// Seg is a segment tree for predecessor/successor queries on active points
type Seg struct {
   size int
   tree []int
}

// NewSeg creates a segment tree with capacity at least n
func NewSeg(n int) *Seg {
   size := 1
   for size < n {
       size <<= 1
   }
   return &Seg{size: size, tree: make([]int, 2*size)}
}

// Update adds delta at position pos
func (s *Seg) Update(pos, delta int) {
   i := pos + s.size
   s.tree[i] += delta
   for i >>= 1; i > 0; i >>= 1 {
       s.tree[i] = s.tree[2*i] + s.tree[2*i+1]
   }
}

// predRec finds the rightmost active position <= q in node idx covering [l,r]
func (s *Seg) predRec(idx, l, r, q int) int {
   if l > q || s.tree[idx] == 0 {
       return -1
   }
   if l == r {
       return l
   }
   mid := (l + r) >> 1
   // try right child first if it overlaps
   if mid < q {
       if res := s.predRec(2*idx+1, mid+1, r, q); res != -1 {
           return res
       }
       return s.predRec(2*idx, l, mid, q)
   }
   // q <= mid
   return s.predRec(2*idx, l, mid, q)
}

// Pred returns the largest active position <= q, or -1 if none
func (s *Seg) Pred(q int) int {
   return s.predRec(1, 0, s.size-1, q)
}

// succRec finds the leftmost active position >= q in node idx covering [l,r]
func (s *Seg) succRec(idx, l, r, q int) int {
   if r < q || s.tree[idx] == 0 {
       return -1
   }
   if l == r {
       return l
   }
   mid := (l + r) >> 1
   // try left child first if it overlaps
   if mid >= q {
       if res := s.succRec(2*idx, l, mid, q); res != -1 {
           return res
       }
       return s.succRec(2*idx+1, mid+1, r, q)
   }
   // q > mid
   return s.succRec(2*idx+1, mid+1, r, q)
}

// Succ returns the smallest active position >= q, or -1 if none
func (s *Seg) Succ(q int) int {
   return s.succRec(1, 0, s.size-1, q)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   fmt.Fscan(in, &n)
   a := make([]int, n+1)
   b := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &a[i])
   }
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &b[i])
   }
   posA := make([]int, n+1)
   posB := make([]int, n+1)
   for i := 1; i <= n; i++ {
       posA[a[i]] = i - 1 // zero-based
       posB[b[i]] = i - 1
   }
   offset := n
   // events at shift k: values x with posB[x] == k
   // events at shift k: values x with posB[x] + 1 == k (transition point from regime1 to regime2)
   events := make([][]int, n+1)
   pos1 := make([]int, n+1)
   pos2 := make([]int, n+1)
   maxPos := 0
   for x := 1; x <= n; x++ {
       g := posB[x]
       // at k == g+1, the element moves from regime1 to regime2
       if g+1 <= n {
           events[g+1] = append(events[g+1], x)
       }
       c1 := g - posA[x]
       c2 := c1 + n
       p1 := c1 + offset
       p2 := c2 + offset
       pos1[x] = p1
       pos2[x] = p2
       if p2 > maxPos {
           maxPos = p2
       }
   }
   // build segment trees over [0..maxPos]
   st1 := NewSeg(maxPos + 1)
   st2 := NewSeg(maxPos + 1)
   // initialize st1 with all pos1
   for x := 1; x <= n; x++ {
       st1.Update(pos1[x], 1)
   }
   // compute answers
   for k := 0; k < n; k++ {
       // move events at k
       for _, x := range events[k] {
           st1.Update(pos1[x], -1)
           st2.Update(pos2[x], 1)
       }
       center := k + offset
       best := n // max possible distance is n-1
       // check st1
       if p := st1.Pred(center); p != -1 {
           d := center - p
           if d < best {
               best = d
           }
       }
       if s := st1.Succ(center); s != -1 {
           d := s - center
           if d < best {
               best = d
           }
       }
       // check st2
       if p := st2.Pred(center); p != -1 {
           d := center - p
           if d < best {
               best = d
           }
       }
       if s := st2.Succ(center); s != -1 {
           d := s - center
           if d < best {
               best = d
           }
       }
       // best is distance in shifted pos, same as answer
       fmt.Fprintln(out, best)
   }
}
