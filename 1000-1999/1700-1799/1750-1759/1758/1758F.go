package main

import (
   "bufio"
   "fmt"
   "math/rand"
   "os"
)

const N = 400000

// Segment tree supporting sum and minimum prefix
type SegTree struct {
   n    int
   sum  []int
   pre  []int
}

func NewSegTree(n int) *SegTree {
   size := 4 * n
   st := &SegTree{n: n, sum: make([]int, size), pre: make([]int, size)}
   // build with initial value -1
   var build func(p, l, r int)
   build = func(p, l, r int) {
       if r-l == 1 {
           st.sum[p] = -1
           st.pre[p] = -1
           return
       }
       m := (l + r) >> 1
       build(2*p, l, m)
       build(2*p+1, m, r)
       st.pull(p)
   }
   build(1, 0, n)
   return st
}

func (st *SegTree) pull(p int) {
   left, right := 2*p, 2*p+1
   st.sum[p] = st.sum[left] + st.sum[right]
   a, b := st.pre[left], st.pre[right]
   st.pre[p] = a
   if st.sum[left]+b < st.pre[p] {
       st.pre[p] = st.sum[left] + b
   }
}

func (st *SegTree) modifyRec(p, l, r, idx, v int) {
   if r-l == 1 {
       st.sum[p] = v
       st.pre[p] = v
       return
   }
   m := (l + r) >> 1
   if idx < m {
       st.modifyRec(2*p, l, m, idx, v)
   } else {
       st.modifyRec(2*p+1, m, r, idx, v)
   }
   st.pull(p)
}

// Modify position idx to value v
func (st *SegTree) Modify(idx, v int) {
   st.modifyRec(1, 0, st.n, idx, v)
}

// findRec searches in [x,y) for first position where prefix sum condition fails
func (st *SegTree) findRec(p, l, r, x, y int, v *int) int {
   if l >= y || r <= x {
       return -1
   }
   if l >= x && r <= y && st.pre[p] > *v {
       *v -= st.sum[p]
       return -1
   }
   if r-l == 1 {
       return r
   }
   m := (l + r) >> 1
   res := st.findRec(2*p, l, m, x, y, v)
   if res == -1 {
       res = st.findRec(2*p+1, m, r, x, y, v)
   }
   return res
}

// Find in [l,r) with initial v
func (st *SegTree) Find(l, r, v int) int {
   return st.findRec(1, 0, st.n, l, r, &v)
}

// Treap for ordered map from int to int
type node struct {
   key, val, prio int
   left, right    *node
}

func split(t *node, key int) (l, r *node) {
   if t == nil {
       return nil, nil
   }
   if t.key < key {
       ll, rr := split(t.right, key)
       t.right = ll
       return t, rr
   }
   ll, rr := split(t.left, key)
   t.left = rr
   return ll, t
}

func merge(a, b *node) *node {
   if a == nil {
       return b
   }
   if b == nil {
       return a
   }
   if a.prio > b.prio {
       a.right = merge(a.right, b)
       return a
   }
   b.left = merge(a, b.left)
   return b
}

func insert(t *node, key, val int) *node {
   n := &node{key: key, val: val, prio: rand.Int()}
   l, r := split(t, key)
   return merge(merge(l, n), r)
}

func erase(t *node, key int) *node {
   l, m := split(t, key)
   m, r := split(m, key+1)
   return merge(l, r)
}

func lowerBound(t *node, key int) *node {
   var res *node
   for t != nil {
       if t.key >= key {
           res = t
           t = t.left
       } else {
           t = t.right
       }
   }
   return res
}

func upperBound(t *node, key int) *node {
   var res *node
   for t != nil {
       if t.key > key {
           res = t
           t = t.left
       } else {
           t = t.right
       }
   }
   return res
}

func predecessor(t *node, key int) *node {
   var res *node
   for t != nil {
       if t.key < key {
           res = t
           t = t.right
       } else {
           t = t.left
       }
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   fmt.Fscan(in, &n)
   tarr := make([]int, N)
   seg := NewSegTree(N)
   var ranges *node
   for i := 0; i < n; i++ {
       var a int
       fmt.Fscan(in, &a)
       a--
       var X [][2]int
       var Y [][2]int
       if tarr[a] == 1 {
           // delete
           tarr[a] = 0
           seg.Modify(a, -1)
           it := predecessor(ranges, a+1)
           l := it.key
           r := it.val
           X = append(X, [2]int{l, r})
           ranges = erase(ranges, l)
           j := seg.Find(l, r, -2)
           if l < j-2 {
               Y = append(Y, [2]int{l, j - 2})
               ranges = insert(ranges, l, j-2)
           }
           if j < r {
               Y = append(Y, [2]int{j, r})
               ranges = insert(ranges, j, r)
           }
       } else {
           // insert
           tarr[a] = 1
           seg.Modify(a, 1)
           pred := predecessor(ranges, a+1)
           if pred != nil && a < pred.val {
               l := pred.key
               r := pred.val
               X = append(X, [2]int{l, r})
               ranges = erase(ranges, l)
               r += 2
               nxt := lowerBound(ranges, r)
               if nxt != nil && nxt.key == r {
                   X = append(X, [2]int{nxt.key, nxt.val})
                   ranges = erase(ranges, nxt.key)
                   r = nxt.val
               }
               Y = append(Y, [2]int{l, r})
               ranges = insert(ranges, l, r)
           } else {
               if a%2 == 1 {
                   a--
               }
               l := a
               r := a + 2
               pred2 := predecessor(ranges, l)
               if pred2 != nil && pred2.val == l {
                   l = pred2.key
                   X = append(X, [2]int{pred2.key, pred2.val})
                   ranges = erase(ranges, pred2.key)
               }
               nxt2 := lowerBound(ranges, r)
               if nxt2 != nil && nxt2.key == r {
                   X = append(X, [2]int{nxt2.key, nxt2.val})
                   ranges = erase(ranges, nxt2.key)
                   r = nxt2.val
               }
               Y = append(Y, [2]int{l, r})
               ranges = insert(ranges, l, r)
           }
       }
       // output
       fmt.Fprintln(out, len(X))
       for _, p := range X {
           fmt.Fprintf(out, "%d %d\n", p[0]+1, p[1])
       }
       fmt.Fprintln(out, len(Y))
       for _, p := range Y {
           fmt.Fprintf(out, "%d %d\n", p[0]+1, p[1])
       }
   }
}
