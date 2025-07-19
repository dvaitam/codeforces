package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, q int
   fmt.Fscan(reader, &n, &q)
   var s string
   fmt.Fscan(reader, &s)
   sa, inv, lcpArr := buildSuffixArray(s)
   rmq := newRMQ(lcpArr)
   getLCP := func(i, j int) int {
       if i == j {
           return n - sa[i]
       }
       if i > j {
           i, j = j, i
       }
       return rmq.query(i, j-1)
   }
   // process queries
   // prepare segment tree once to reuse
   st := newSegTree(n + 1)
   for qi := 0; qi < q; qi++ {
       var k, l int
       fmt.Fscan(reader, &k, &l)
       A := make([]int, k)
       for i := 0; i < k; i++ {
           fmt.Fscan(reader, &A[i])
           A[i]--
       }
       B := make([]int, l)
       for i := 0; i < l; i++ {
           fmt.Fscan(reader, &B[i])
           B[i]--
       }
       res := solveQuery(A, B, inv, getLCP, n, st)
       fmt.Fprintln(writer, res)
   }
}

// solveQuery implements the two-sweep algorithm using a segment tree over LCP values
// solveQuery implements the two-sweep algorithm using a segment tree over LCP values
func solveQuery(A, B, inv []int, getLCP func(int, int) int, maxLCP int, st *segTree) int64 {
   // clear any previous state
   st.clearRange(0, maxLCP)
   var sum int64
   var res int64
   pos := 0
   // first sweep
   for i := 0; i < len(A); i++ {
       now := inv[A[i]]
       if i > 0 {
           x := getLCP(inv[A[i-1]], now)
           if x < maxLCP {
               cnt, s := st.queryRange(x+1, maxLCP)
               if cnt > 0 {
                   sum -= (s - int64(x)*cnt)
                   st.clearRange(x+1, maxLCP)
                   st.updatePoint(x, cnt)
               }
           }
       }
       for pos < len(B) && inv[B[pos]] <= now {
           x := getLCP(inv[B[pos]], now)
           sum += int64(x)
           st.updatePoint(x, 1)
           pos++
       }
       res += sum
   }
   // second sweep
   st.clearRange(0, maxLCP)
   sum = 0
   pos = len(B) - 1
   for i := len(A) - 1; i >= 0; i-- {
       now := inv[A[i]]
       if i < len(A)-1 {
           x := getLCP(now, inv[A[i+1]])
           if x < maxLCP {
               cnt, s := st.queryRange(x+1, maxLCP)
               if cnt > 0 {
                   sum -= (s - int64(x)*cnt)
                   st.clearRange(x+1, maxLCP)
                   st.updatePoint(x, cnt)
               }
           }
       }
       for pos >= 0 && inv[B[pos]] > now {
           x := getLCP(now, inv[B[pos]])
           sum += int64(x)
           st.updatePoint(x, 1)
           pos--
       }
       res += sum
   }
   return res
}

// suffix array and LCP
func buildSuffixArray(s string) (sa, inv, lcp []int) {
   n := len(s)
   sa = make([]int, n)
   ranks := make([]int, n)
   tmp := make([]int, n)
   for i := 0; i < n; i++ {
       sa[i] = i
       ranks[i] = int(s[i])
   }
   for k := 1; k < n; k <<= 1 {
       // radix sort by second key then first key
       sa = radixPass(sa, ranks, k)
       sa = radixPass(sa, ranks, 0)
       tmp[sa[0]] = 0
       for i := 1; i < n; i++ {
           prev, cur := sa[i-1], sa[i]
           if ranks[prev] == ranks[cur] && getRank(ranks, prev+k, n) == getRank(ranks, cur+k, n) {
               tmp[cur] = tmp[prev]
           } else {
               tmp[cur] = tmp[prev] + 1
           }
       }
       copy(ranks, tmp)
       if ranks[sa[n-1]] == n-1 {
           break
       }
   }
   // build inv
   inv = make([]int, n)
   for i := 0; i < n; i++ {
       inv[sa[i]] = i
   }
   // build LCP
   lcp = make([]int, n-1)
   h := 0
   for i := 0; i < n; i++ {
       if inv[i] > 0 {
           j := sa[inv[i]-1]
           for i+h < n && j+h < n && s[i+h] == s[j+h] {
               h++
           }
           lcp[inv[i]-1] = h
           if h > 0 {
               h--
           }
       }
   }
   return
}

func getRank(ranks []int, i, n int) int {
   if i < n {
       return ranks[i]
   }
   return -1
}

// radixPass sorts sa by ranks at offset k
func radixPass(sa, ranks []int, k int) []int {
   n := len(sa)
   maxv := max(256, n) + 1
   cnt := make([]int, maxv)
   tmp := make([]int, n)
   for _, v := range sa {
       r := getRank(ranks, v+k, n) + 1
       cnt[r]++
   }
   for i := 1; i < maxv; i++ {
       cnt[i] += cnt[i-1]
   }
   for i := n - 1; i >= 0; i-- {
       v := sa[i]
       r := getRank(ranks, v+k, n) + 1
       cnt[r]--
       tmp[cnt[r]] = v
   }
   return tmp
}

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

// RMQ sparse table for min
type RMQ struct {
   log []int
   st  [][]int
}

func newRMQ(a []int) *RMQ {
   n := len(a)
   log := make([]int, n+1)
   for i := 2; i <= n; i++ {
       log[i] = log[i/2] + 1
   }
   k := log[n] + 1
   st := make([][]int, k)
   st[0] = make([]int, n)
   copy(st[0], a)
   for i := 1; i < k; i++ {
       st[i] = make([]int, n-(1<<i)+1)
       for j := 0; j+ (1<<i) <= n; j++ {
           x, y := st[i-1][j], st[i-1][j+(1<<(i-1))]
           if x < y {
               st[i][j] = x
           } else {
               st[i][j] = y
           }
       }
   }
   return &RMQ{log, st}
}

func (r *RMQ) query(l, rgt int) int {
   j := r.log[rgt-l+1]
   x, y := r.st[j][l], r.st[j][rgt-(1<<j)+1]
   if x < y {
       return x
   }
   return y
}

// segment tree with range clear and point update
type segTree struct {
   n    int
   cnt  []int64
   sum  []int64
   lazy []bool
}

func newSegTree(n int) *segTree {
   size := 1
   for size < n {
       size <<= 1
   }
   cnt := make([]int64, 2*size)
   sum := make([]int64, 2*size)
   lazy := make([]bool, 2*size)
   return &segTree{n: size, cnt: cnt, sum: sum, lazy: lazy}
}

func (st *segTree) push(node int) {
   if st.lazy[node] {
       for _, ch := range []int{node * 2, node*2 + 1} {
           st.cnt[ch] = 0
           st.sum[ch] = 0
           st.lazy[ch] = true
       }
       st.lazy[node] = false
   }
}

func (st *segTree) updatePoint(idx int, v int64) {
   st.update(1, 0, st.n-1, idx, v)
}

func (st *segTree) update(node, l, r, idx int, v int64) {
   if l == r {
       st.cnt[node] += v
       st.sum[node] += int64(idx) * v
       return
   }
   st.push(node)
   mid := (l + r) >> 1
   if idx <= mid {
       st.update(node*2, l, mid, idx, v)
   } else {
       st.update(node*2+1, mid+1, r, idx, v)
   }
   st.cnt[node] = st.cnt[node*2] + st.cnt[node*2+1]
   st.sum[node] = st.sum[node*2] + st.sum[node*2+1]
}

func (st *segTree) queryRange(L, R int) (int64, int64) {
   return st.query(1, 0, st.n-1, L, R)
}

func (st *segTree) query(node, l, r, L, R int) (int64, int64) {
   if R < l || r < L {
       return 0, 0
   }
   if L <= l && r <= R {
       return st.cnt[node], st.sum[node]
   }
   st.push(node)
   mid := (l + r) >> 1
   c1, s1 := st.query(node*2, l, mid, L, R)
   c2, s2 := st.query(node*2+1, mid+1, r, L, R)
   return c1 + c2, s1 + s2
}

func (st *segTree) clearRange(L, R int) {
   st.clear(1, 0, st.n-1, L, R)
}

func (st *segTree) clear(node, l, r, L, R int) {
   if R < l || r < L {
       return
   }
   if L <= l && r <= R {
       st.cnt[node] = 0
       st.sum[node] = 0
       st.lazy[node] = true
       return
   }
   st.push(node)
   mid := (l + r) >> 1
   st.clear(node*2, l, mid, L, R)
   st.clear(node*2+1, mid+1, r, L, R)
   st.cnt[node] = st.cnt[node*2] + st.cnt[node*2+1]
   st.sum[node] = st.sum[node*2] + st.sum[node*2+1]
}
