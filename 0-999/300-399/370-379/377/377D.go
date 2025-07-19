package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type Segment struct {
   l, x, r, idx int
}

// Segment tree for range add and range max with index
type SegTree struct {
   n    int
   mx   []int
   idx  []int
   add  []int
}

func NewSegTree(n int) *SegTree {
   size := 1
   for size < n {
       size <<= 1
   }
   mx := make([]int, 2*size)
   idx := make([]int, 2*size)
   add := make([]int, 2*size)
   // initialize indices at leaves
   for i := 0; i < size; i++ {
       if i < n {
           idx[size+i] = i + 1 // positions 1..n
       } else {
           idx[size+i] = 1
       }
       mx[size+i] = 0
   }
   for i := size - 1; i > 0; i-- {
       // combine children
       li, ri := 2*i, 2*i+1
       if mx[li] >= mx[ri] {
           mx[i] = mx[li]
           idx[i] = idx[li]
       } else {
           mx[i] = mx[ri]
           idx[i] = idx[ri]
       }
   }
   return &SegTree{n: size, mx: mx, idx: idx, add: add}
}

func (st *SegTree) apply(p, v int) {
   st.mx[p] += v
   st.add[p] += v
}

func (st *SegTree) push(p int) {
   if st.add[p] != 0 {
       st.apply(p<<1, st.add[p])
       st.apply(p<<1|1, st.add[p])
       st.add[p] = 0
   }
}

func (st *SegTree) pull(p int) {
   li, ri := p<<1, p<<1|1
   if st.mx[li] >= st.mx[ri] {
       st.mx[p] = st.mx[li]
       st.idx[p] = st.idx[li]
   } else {
       st.mx[p] = st.mx[ri]
       st.idx[p] = st.idx[ri]
   }
}

// update range [l..r] add v, current node p covers [L..R]
func (st *SegTree) update(p, L, R, l, r, v int) {
   if l > R || r < L {
       return
   }
   if l <= L && R <= r {
       st.apply(p, v)
       return
   }
   st.push(p)
   M := (L + R) >> 1
   st.update(p<<1, L, M, l, r, v)
   st.update(p<<1|1, M+1, R, l, r, v)
   st.pull(p)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   fmt.Fscan(in, &n)
   segs := make([]Segment, n)
   m := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &segs[i].l, &segs[i].x, &segs[i].r)
       segs[i].idx = i + 1
       if segs[i].r > m {
           m = segs[i].r
       }
   }
   // sort by x and by r
   px := make([]Segment, n)
   pr := make([]Segment, n)
   copy(px, segs)
   copy(pr, segs)
   sort.Slice(px, func(i, j int) bool { return px[i].x < px[j].x })
   sort.Slice(pr, func(i, j int) bool { return pr[i].r < pr[j].r })
   st := NewSegTree(m)
   cpx, cpr := 0, 0
   ans, al, ar := 0, 1, 1
   // iterate through possible right endpoints
   for i := 1; i <= m; i++ {
       // add segments with x == i
       for cpx < n && px[cpx].x == i {
           l := px[cpx].l
           x := px[cpx].x
           if l <= x {
               st.update(1, 1, st.n, l, x, 1)
           }
           cpx++
       }
       // query max
       if st.mx[1] > ans {
           ans = st.mx[1]
           al = st.idx[1]
           ar = i
       }
       // remove segments with r == i
       for cpr < n && pr[cpr].r == i {
           l := pr[cpr].l
           x := pr[cpr].x
           if l <= x {
               st.update(1, 1, st.n, l, x, -1)
           }
           cpr++
       }
   }
   // output
   fmt.Fprintln(out, ans)
   // print indices of chosen segments
   for i := 0; i < n; i++ {
       if segs[i].x >= al && segs[i].x <= ar && segs[i].l <= al && segs[i].r >= ar {
           fmt.Fprintf(out, "%d ", segs[i].idx)
       }
   }
   fmt.Fprintln(out)
}
