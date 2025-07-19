package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const NSEG = 1 << 19
const INF64 = 1 << 60

// Node represents a segment tree node
type Node struct {
   pt, x, maxx, tot int64
}

var tree [2 * NSEG]Node

// add combines two nodes
func add(a, b Node) Node {
   var res Node
   // total prefix sum combine
   if a.tot < b.tot + a.pt {
       res.pt = a.pt + b.pt
       res.x = b.x
       res.maxx = b.maxx
       if a.tot > b.tot + a.pt {
           res.tot = a.tot
       } else {
           res.tot = b.tot + a.pt
       }
   } else {
       res.pt = a.pt + b.pt
       res.x = b.x
       res.maxx = a.maxx
       if a.tot > b.tot + a.pt {
           res.tot = a.tot
       } else {
           res.tot = b.tot + a.pt
       }
   }
   return res
}

// setV initializes leaf at p with value v and coordinate x
func setV(p int, v, x int64) {
   idx := p + NSEG
   if tree[idx].x != x {
       tree[idx].x = x
       tree[idx].maxx = x
       tree[idx].tot -= x
   }
   tree[idx].pt += v
   tree[idx].tot += v
}

// build constructs the segment tree
func build() {
   for i := NSEG - 1; i >= 1; i-- {
       tree[i] = add(tree[i<<1], tree[i<<1|1])
   }
}

// update subtracts v from leaf at p and updates ancestors
func update(p int, v int64) {
   idx := p + NSEG
   tree[idx].pt -= v
   tree[idx].tot -= v
   for idx >>= 1; idx >= 1; idx >>= 1 {
       tree[idx] = add(tree[idx<<1], tree[idx<<1|1])
   }
}

// query returns combined node over [l, r]
func query(l, r int) Node {
   l += NSEG
   r += NSEG
   lsum := Node{pt: 0, x: 0, maxx: 0, tot: -INF64}
   rsum := Node{pt: 0, x: INF64, maxx: 0, tot: -INF64}
   for l <= r {
       if l&1 == 1 {
           lsum = add(lsum, tree[l])
           l++
       }
       if r&1 == 0 {
           rsum = add(tree[r], rsum)
           r--
       }
       l >>= 1
       r >>= 1
   }
   return add(lsum, rsum)
}

// Ent holds a triple a, b-index, c-value
type Ent struct{ a, b, c int64 }

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var N int
   fmt.Fscan(in, &N)
   ents := make([]Ent, N)
   X := make([]int64, 0, N)
   for i := 0; i < N; i++ {
       var a, b, c int64
       fmt.Fscan(in, &a, &b, &c)
       if a > b {
           a, b = b, a
       }
       ents[i] = Ent{a: a, b: b, c: c}
       X = append(X, b)
   }
   sort.Slice(X, func(i, j int) bool { return X[i] < X[j] })
   // unique
   m := 0
   for i := 0; i < len(X); i++ {
       if m == 0 || X[i] != X[m-1] {
           X[m] = X[i]
           m++
       }
   }
   X = X[:m]
   // init segment tree leaves
   for i := 0; i < N; i++ {
       bi := sort.Search(len(X), func(i int) bool { return X[i] >= ents[i].b })
       ents[i].b = int64(bi)
       setV(bi, ents[i].c, X[bi])
   }
   build()
   // sort by a
   sort.Slice(ents, func(i, j int) bool { return ents[i].a < ents[j].a })

   var resVal, resL, resR int64
   resVal = 0
   resL = INF64
   resR = INF64
   // process by groups of equal a
   for i := 0; i < N; {
       a0 := ents[i].a
       // query
       pi := sort.Search(len(X), func(j int) bool { return X[j] >= a0 })
       r := query(pi, len(X)-1)
       r.tot += a0
       if r.tot > resVal {
           resVal = r.tot
           resL = a0
           resR = r.maxx
       }
       // update all in group
       j := i
       for j < N && ents[j].a == a0 {
           update(int(ents[j].b), ents[j].c)
           j++
       }
       i = j
   }
   // output
   fmt.Fprintln(out, resVal)
   // print a a b b
   fmt.Fprintf(out, "%d %d %d %d\n", resL, resL, resR, resR)
}
