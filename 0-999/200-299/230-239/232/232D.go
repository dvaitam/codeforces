package main

import (
   "bufio"
   "fmt"
   "os"
)

// Suffix array with merge sort tree and double rolling hash to count matching segments
func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   fmt.Fscan(in, &n)
   h := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &h[i])
   }
   // build D and F arrays
   Dlen := n - 1
   D := make([]int64, Dlen)
   F := make([]int64, Dlen)
   for i := 0; i < Dlen; i++ {
       D[i] = h[i+1] - h[i]
       F[i] = -D[i]
   }
   // rolling hash params
   const mod1 = 1000000007
   const mod2 = 1000000009
   const base = 91138233
   // precompute powers
   pow1 := make([]int64, Dlen+5)
   pow2 := make([]int64, Dlen+5)
   pow1[0], pow2[0] = 1, 1
   for i := 1; i <= Dlen; i++ {
       pow1[i] = pow1[i-1] * base % mod1
       pow2[i] = pow2[i-1] * base % mod2
   }
   // prefix hashes
   H1 := make([]int64, Dlen+1)
   H2 := make([]int64, Dlen+1)
   HF1 := make([]int64, Dlen+1)
   HF2 := make([]int64, Dlen+1)
   for i := 0; i < Dlen; i++ {
       v := (D[i] % mod1 + mod1) % mod1
       H1[i+1] = (H1[i]*base + v) % mod1
       v2 := (D[i] % mod2 + mod2) % mod2
       H2[i+1] = (H2[i]*base + v2) % mod2
       fv1 := (F[i] % mod1 + mod1) % mod1
       HF1[i+1] = (HF1[i]*base + fv1) % mod1
       fv2 := (F[i] % mod2 + mod2) % mod2
       HF2[i+1] = (HF2[i]*base + fv2) % mod2
   }
   // suffix array of D
   sa := buildSA(D)
   // merge sort tree on sa positions
   mst := newMergeTree(sa)
   // helper to get hash of D[u:u+ln]
   getD := func(u, ln int) (int64, int64) {
       a := (H1[u+ln] - H1[u]*pow1[ln] % mod1 + mod1) % mod1
       b := (H2[u+ln] - H2[u]*pow2[ln] % mod2 + mod2) % mod2
       return a, b
   }
   // hash of F[l:l+ln]
   getF := func(l, ln int) (int64, int64) {
       a := (HF1[l+ln] - HF1[l]*pow1[ln] % mod1 + mod1) % mod1
       b := (HF2[l+ln] - HF2[l]*pow2[ln] % mod2 + mod2) % mod2
       return a, b
   }
   // compare suffix at u with pattern F[l:l+m]
   // returns true if suffix < pattern
   isSufLess := func(u, l, m int) bool {
       // binary search lcp
       lo, hi := 0, m
       for lo < hi {
           mid := (lo + hi + 1) >> 1
           d1, d2 := getD(u, mid)
           f1, f2 := getF(l, mid)
           if d1 == f1 && d2 == f2 {
               lo = mid
           } else {
               hi = mid - 1
           }
       }
       lcp := lo
       if lcp == m {
           return false
       }
       // compare next char
       dv := D[u+lcp]
       fv := F[l+lcp]
       return dv < fv
   }
   // check suffix <= pattern
   isSufLe := func(u, l, m int) bool {
       lo, hi := 0, m
       for lo < hi {
           mid := (lo + hi + 1) >> 1
           d1, d2 := getD(u, mid)
           f1, f2 := getF(l, mid)
           if d1 == f1 && d2 == f2 {
               lo = mid
           } else {
               hi = mid - 1
           }
       }
       lcp := lo
       if lcp == m {
           return true
       }
       dv := D[u+lcp]
       fv := F[l+lcp]
       return dv <= fv
   }
   // process queries
   var q int
   fmt.Fscan(in, &q)
   for qi := 0; qi < q; qi++ {
       var l, r int
       fmt.Fscan(in, &l, &r)
       l--, r--
       w := r - l + 1
       if w == 1 {
           // any other plank
           fmt.Fprintln(out, n-1)
           continue
       }
       m := w - 1
       // find SA interval where suffixes match pattern F[l:l+m]
       // lower bound
       lo, hi := 0, len(sa)-1
       L := len(sa)
       for lo <= hi {
           mid := (lo + hi) >> 1
           if isSufLess(sa[mid], l, m) {
               lo = mid + 1
           } else {
               L = mid
               hi = mid - 1
           }
       }
       // upper bound
       lo, hi = 0, len(sa)-1
       R := -1
       for lo <= hi {
           mid := (lo + hi) >> 1
           if isSufLe(sa[mid], l, m) {
               R = mid
               lo = mid + 1
           } else {
               hi = mid - 1
           }
       }
       if L > R {
           fmt.Fprintln(out, 0)
           continue
       }
       // count starting positions u in [1..l-w] and [r+1..n-w+1]
       // u is D index 0-based, corresponds to h index u
       // u <= l-w => u <= l - w => l-w = l - (m+1) => l - m - 1
       max1 := l - m - 1
       cnt := 0
       if max1 >= 0 {
           cnt += mst.query(0, 0, Dlen-1, L, R, 0, max1)
       }
       // u >= r+1
       min2 := r + 1
       if min2 <= Dlen-m {
           cnt += mst.query(0, 0, Dlen-1, L, R, min2, Dlen-m)
       }
       fmt.Fprintln(out, cnt)
   }
}

// build suffix array of int64 slice D
func buildSA(D []int64) []int {
   n := len(D)
   sa := make([]int, n)
   rnk := make([]int, n)
   tmp := make([]int, n)
   // compress values
   vals := append([]int64(nil), D...)
   // copy and sort unique
   m := len(vals)
   // map values to rank
   mp := make(map[int64]int, m)
   vs := make([]int64, m)
   copy(vs, vals)
   // sort vs
   for i := 0; i < m; i++ {
       for j := i + 1; j < m; j++ {
           if vs[j] < vs[i] {
               vs[i], vs[j] = vs[j], vs[i]
           }
       }
   }
   uniq := make([]int64, 0, m)
   for i, v := range vs {
       if i == 0 || v != vs[i-1] {
           uniq = append(uniq, v)
       }
   }
   for i, v := range uniq {
       mp[v] = i
   }
   for i := 0; i < n; i++ {
       rnk[i] = mp[D[i]]
       sa[i] = i
   }
   k := 1
   tmpSA := make([]int, n)
   for k < n {
       // sort by second key then first key via counting sort
       // second key: rnk[i+k]
       maxv := n
       cnt := make([]int, maxv+1)
       for i := 0; i < n; i++ {
           key := 0
           if sa[i]+k < n {
               key = rnk[sa[i]+k] + 1
           }
           cnt[key]++
       }
       for i := 1; i <= maxv; i++ {
           cnt[i] += cnt[i-1]
       }
       for i := n - 1; i >= 0; i-- {
           key := 0
           if sa[i]+k < n {
               key = rnk[sa[i]+k] + 1
           }
           cnt[key]--
           tmpSA[cnt[key]] = sa[i]
       }
       // sort by first key
       cnt = make([]int, maxv+1)
       for i := 0; i < n; i++ {
           key := rnk[tmpSA[i]] + 1
           cnt[key]++
       }
       for i := 1; i <= maxv; i++ {
           cnt[i] += cnt[i-1]
       }
       for i := n - 1; i >= 0; i-- {
           key := rnk[tmpSA[i]] + 1
           cnt[key]--
           sa[cnt[key]] = tmpSA[i]
       }
       // update ranks
       tmp[sa[0]] = 0
       p := 0
       for i := 1; i < n; i++ {
           a, b := sa[i-1], sa[i]
           if rnk[a] != rnk[b] ||
               (a+k < n && b+k < n && rnk[a+k] != rnk[b+k]) ||
               (a+k >= n && b+k < n) || (a+k < n && b+k >= n) {
               p++
           }
           tmp[b] = p
       }
       copy(rnk, tmp)
       if p == n-1 {
           break
       }
       k <<= 1
   }
   return sa
}

// merge sort tree
type mergeTree struct {
   sa []int
   t  [][]int
}

func newMergeTree(sa []int) *mergeTree {
   n := len(sa)
   size := 1
   for size < n*4 {
       size <<= 1
   }
   t := make([][]int, n*4)
   mt := &mergeTree{sa: sa, t: t}
   mt.build(0, 0, n-1)
   return mt
}

func (mt *mergeTree) build(node, l, r int) {
   if l == r {
       mt.t[node] = []int{mt.sa[l]}
       return
   }
   mid := (l + r) >> 1
   lc, rc := node*2+1, node*2+2
   mt.build(lc, l, mid)
   mt.build(rc, mid+1, r)
   a, b := mt.t[lc], mt.t[rc]
   m1, m2 := len(a), len(b)
   merged := make([]int, m1+m2)
   i, j := 0, 0
   for i < m1 || j < m2 {
       if j == m2 || (i < m1 && a[i] < b[j]) {
           merged[i+j] = a[i]
           i++
       } else {
           merged[i+j] = b[j]
           j++
       }
   }
   mt.t[node] = merged
}

// query values in [L..R] of suffix array indices, counting sa value in [lo..hi]
func (mt *mergeTree) query(node, l, r, L, R, lo, hi int) int {
   if r < L || l > R {
       return 0
   }
   if L <= l && r <= R {
       arr := mt.t[node]
       // count in arr values between lo and hi
       // lower bound lo
       cntL := lower(arr, lo)
       cntR := lower(arr, hi+1)
       return cntR - cntL
   }
   mid := (l + r) >> 1
   return mt.query(node*2+1, l, mid, L, R, lo, hi) + mt.query(node*2+2, mid+1, r, L, R, lo, hi)
}

// lower returns first idx with arr[idx]>=x
func lower(arr []int, x int) int {
   lo, hi := 0, len(arr)
   for lo < hi {
       mid := (lo + hi) >> 1
       if arr[mid] < x {
           lo = mid + 1
       } else {
           hi = mid
       }
   }
   return lo
}
