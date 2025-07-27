package main

import (
   "bufio"
   "fmt"
   "math/rand"
   "os"
   "sort"
)

// Seg represents a road segment endpoints
type Seg struct {
   id       int
   lx, ly   int     // left endpoint (min x, then min y)
   hx, hy   int     // other endpoint
   slope    float64 // used for y interpolation
}

// node for treap
type node struct {
   seg       *Seg
   pri       int
   left, right *node
}

// global for comparator
var currX float64

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   // read segments
   segs := make([]*Seg, n)
   for i := 0; i < n; i++ {
       var x1, y1, x2, y2 int
       fmt.Fscan(in, &x1, &y1, &x2, &y2)
       var lx, ly, hx, hy int
       if x1 < x2 || (x1 == x2 && y1 < y2) {
           lx, ly, hx, hy = x1, y1, x2, y2
       } else {
           lx, ly, hx, hy = x2, y2, x1, y1
       }
       s := &Seg{id: i, lx: lx, ly: ly, hx: hx, hy: hy}
       if hx != lx {
           s.slope = float64(hy-ly) / float64(hx-lx)
       }
       segs[i] = s
   }
   // sort by left endpoint
   sort.Slice(segs, func(i, j int) bool {
       if segs[i].lx != segs[j].lx {
           return segs[i].lx < segs[j].lx
       }
       return segs[i].ly < segs[j].ly
   })
   // map id to segment pointer for output
   id2seg := make([]*Seg, n)
   for _, s := range segs {
       id2seg[s.id] = s
   }
   // treap root
   var root *node
   // fixed seed rand
   rnd := rand.New(rand.NewSource(42))
   // store connections (u,v)
   conns := make([][2]int, 0, n-1)
   // sweep insert
   for _, s := range segs {
       currX = float64(s.lx)
       pred, succ := findPredSucc(root, s)
       if pred != nil {
           conns = append(conns, [2]int{pred.id, s.id})
       } else if succ != nil {
           conns = append(conns, [2]int{succ.id, s.id})
       }
       root = treapInsert(root, s, rnd)
   }
   // output new roads
   for _, p := range conns {
       u := p[0]
       v := p[1]
       su := id2seg[u]
       sv := id2seg[v]
       fmt.Fprintf(out, "%d %d %d %d\n", su.hx, su.hy, sv.lx, sv.ly)
   }
}

// findPredSucc finds predecessor and successor of s in treap by y-order at currX
func findPredSucc(root *node, s *Seg) (*Seg, *Seg) {
   var pred, succ *Seg
   for n := root; n != nil; {
       if keyLess(s, n.seg) {
           succ = n.seg
           n = n.left
       } else {
           pred = n.seg
           n = n.right
       }
   }
   return pred, succ
}

// treapInsert inserts s into treap rooted at root, returns new root
func treapInsert(root *node, s *Seg, rnd *rand.Rand) *node {
   if root == nil {
       return &node{seg: s, pri: rnd.Int(), left: nil, right: nil}
   }
   if keyLess(s, root.seg) {
       root.left = treapInsert(root.left, s, rnd)
       if root.left.pri < root.pri {
           root = rotateRight(root)
       }
   } else {
       root.right = treapInsert(root.right, s, rnd)
       if root.right.pri < root.pri {
           root = rotateLeft(root)
       }
   }
   return root
}

func rotateRight(t *node) *node {
   l := t.left
   t.left = l.right
   l.right = t
   return l
}

func rotateLeft(t *node) *node {
   r := t.right
   t.right = r.left
   r.left = t
   return r
}

// keyLess reports whether a is below b at current currX
func keyLess(a, b *Seg) bool {
   ya := yAt(a, currX)
   yb := yAt(b, currX)
   if ya != yb {
       return ya < yb
   }
   return a.id < b.id
}

// yAt computes y coordinate of segment s at x = currX
func yAt(s *Seg, x float64) float64 {
   if s.hx == s.lx {
       return float64(s.ly)
   }
   return float64(s.ly) + s.slope*(x-float64(s.lx))
}
