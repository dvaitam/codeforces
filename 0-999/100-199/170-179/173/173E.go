package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Fenwick tree for sum
type Fenwick struct {
   n    int
   tree []int
}

func NewFenwick(n int) *Fenwick {
   return &Fenwick{n: n, tree: make([]int, n+1)}
}

func (f *Fenwick) Update(i, v int) {
   i++
   for ; i <= f.n; i += i & -i {
       f.tree[i] += v
   }
}

func (f *Fenwick) Query(i int) int {
   // sum [0..i]
   i++
   s := 0
   for ; i > 0; i -= i & -i {
       s += f.tree[i]
   }
   return s
}

// Segment tree for range max query
type SegTree struct {
   n    int
   size int
   tree []int
}

func NewSegTree(n int) *SegTree {
   size := 1
   for size < n {
       size <<= 1
   }
   return &SegTree{n: n, size: size, tree: make([]int, 2*size)}
}

func (s *SegTree) Update(i, v int) {
   i += s.size
   if s.tree[i] >= v {
       return
   }
   s.tree[i] = v
   for i >>= 1; i > 0; i >>= 1 {
       if s.tree[2*i] > s.tree[2*i+1] {
           s.tree[i] = s.tree[2*i]
       } else {
           s.tree[i] = s.tree[2*i+1]
       }
   }
}

func (s *SegTree) Query(l, r int) int {
   if l > r {
       return 0
   }
   l += s.size
   r += s.size
   res := 0
   for l <= r {
       if l&1 == 1 {
           if s.tree[l] > res {
               res = s.tree[l]
           }
           l++
       }
       if r&1 == 0 {
           if s.tree[r] > res {
               res = s.tree[r]
           }
           r--
       }
       l >>= 1
       r >>= 1
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   var ageK int
   fmt.Fscan(in, &n, &ageK)
   r := make([]int, n)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &r[i])
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // compress ages
   ages := make([]int, n)
   copy(ages, a)
   sort.Ints(ages)
   ages = uniqueInts(ages)
   m := len(ages)
   idxAge := make([]int, n)
   for i := 0; i < n; i++ {
       idxAge[i] = sort.SearchInts(ages, a[i])
   }
   // compute c_i
   c := make([]int, n)
   fenw := NewFenwick(m)
   // order by responsibility asc
   ord := make([]int, n)
   for i := range ord {
       ord[i] = i
   }
   sort.Slice(ord, func(i, j int) bool {
       return r[ord[i]] < r[ord[j]]
   })
   for i := 0; i < n; {
       j := i
       // same r
       for j < n && r[ord[j]] == r[ord[i]] {
           j++
       }
       // insert all in [i,j)
       for k := i; k < j; k++ {
           fenw.Update(idxAge[ord[k]], 1)
       }
       // compute c for each
       for k := i; k < j; k++ {
           u := ord[k]
           lo := a[u] - ageK
           hi := a[u] + ageK
           lidx := sort.SearchInts(ages, lo)
           ridx := sort.Search(len(ages), func(x int) bool { return ages[x] > hi }) - 1
           if lidx < 0 {
               lidx = 0
           }
           if ridx >= 0 && lidx <= ridx {
               c[u] = fenw.Query(ridx)
               if lidx > 0 {
                   c[u] -= fenw.Query(lidx - 1)
               }
           } else {
               c[u] = 0
           }
       }
       i = j
   }
   // read queries
   var q int
   fmt.Fscan(in, &q)
   type Query struct {
       R0, lidx, ridx, idx int
   }
   qs := make([]Query, q)
   for i := 0; i < q; i++ {
       var x, y int
       fmt.Fscan(in, &x, &y)
       x--
       y--
       R0 := r[x]
       if r[y] > R0 {
           R0 = r[y]
       }
       lowAge := max(a[x]-ageK, a[y]-ageK)
       highAge := min(a[x]+ageK, a[y]+ageK)
       lidx := sort.SearchInts(ages, lowAge)
       ridx := sort.Search(len(ages), func(u int) bool { return ages[u] > highAge }) - 1
       qs[i] = Query{R0: R0, lidx: lidx, ridx: ridx, idx: i}
   }
   // sort members by R desc
   ordDesc := make([]int, n)
   copy(ordDesc, ord)
   sort.Slice(ordDesc, func(i, j int) bool {
       return r[ordDesc[i]] > r[ordDesc[j]]
   })
   // sort queries by R0 desc
   sort.Slice(qs, func(i, j int) bool {
       return qs[i].R0 > qs[j].R0
   })
   // segtree init
   seg := NewSegTree(m)
   ans := make([]int, q)
   pi := 0
   for _, qu := range qs {
       for pi < n && r[ordDesc[pi]] >= qu.R0 {
           u := ordDesc[pi]
           seg.Update(idxAge[u], c[u])
           pi++
       }
       if qu.lidx > qu.ridx {
           ans[qu.idx] = -1
       } else {
           mv := seg.Query(qu.lidx, qu.ridx)
           if mv == 0 {
               ans[qu.idx] = -1
           } else {
               ans[qu.idx] = mv
           }
       }
   }
   // output answers
   for i := 0; i < q; i++ {
       fmt.Fprintln(out, ans[i])
   }
}

func uniqueInts(a []int) []int {
   n := 0
   for i := 0; i < len(a); i++ {
       if n == 0 || a[i] != a[n-1] {
           a[n] = a[i]
           n++
       }
   }
   return a[:n]
}

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}
