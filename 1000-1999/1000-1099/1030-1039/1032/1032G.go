package main

import (
   "bufio"
   "fmt"
   "os"
)

// segTree supports range minimum and maximum queries on int32 arrays
type segTree struct {
   n    int
   size int
   minv []int32
   maxv []int32
}

// newSegTree creates a segment tree for n elements
func newSegTree(n int) *segTree {
   size := 1
   for size < n {
       size <<= 1
   }
   return &segTree{
       n:    n,
       size: size,
       minv: make([]int32, 2*size),
       maxv: make([]int32, 2*size),
   }
}

// build constructs the tree using leaves from larr and rarr
func (st *segTree) build(larr, rarr []int32) {
   const INF = int32(1e9)
   // initialize leaves
   for i := 0; i < st.size; i++ {
       if i < st.n {
           st.minv[st.size+i] = larr[i]
           st.maxv[st.size+i] = rarr[i]
       } else {
           st.minv[st.size+i] = INF
           st.maxv[st.size+i] = -INF
       }
   }
   // build internal nodes
   for i := st.size - 1; i > 0; i-- {
       // min
       if st.minv[2*i] < st.minv[2*i+1] {
           st.minv[i] = st.minv[2*i]
       } else {
           st.minv[i] = st.minv[2*i+1]
       }
       // max
       if st.maxv[2*i] > st.maxv[2*i+1] {
           st.maxv[i] = st.maxv[2*i]
       } else {
           st.maxv[i] = st.maxv[2*i+1]
       }
   }
}

// queryMin returns the minimum value in [l, r]
func (st *segTree) queryMin(l, r int) int32 {
   res := int32(1e9)
   l += st.size
   r += st.size
   for l <= r {
       if l&1 == 1 {
           if st.minv[l] < res {
               res = st.minv[l]
           }
           l++
       }
       if r&1 == 0 {
           if st.minv[r] < res {
               res = st.minv[r]
           }
           r--
       }
       l >>= 1
       r >>= 1
   }
   return res
}

// queryMax returns the maximum value in [l, r]
func (st *segTree) queryMax(l, r int) int32 {
   res := int32(-1e9)
   l += st.size
   r += st.size
   for l <= r {
       if l&1 == 1 {
           if st.maxv[l] > res {
               res = st.maxv[l]
           }
           l++
       }
       if r&1 == 0 {
           if st.maxv[r] > res {
               res = st.maxv[r]
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
   fmt.Fscan(in, &n)
   r := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &r[i])
   }
   if n == 1 {
       fmt.Fprintln(out, 0)
       return
   }
   // unroll circle 3 times
   M := 3 * n
   r2 := make([]int, M)
   for i := 0; i < M; i++ {
       r2[i] = r[i%n]
   }
   // initial reach intervals
   L := make([][]int32, 0)
   R := make([][]int32, 0)
   // determine levels K for doubling
   K := 0
   for (1 << K) < n {
       K++
   }
   K++
   // level 0
   L0 := make([]int32, M)
   R0 := make([]int32, M)
   for i := 0; i < M; i++ {
       l := i - r2[i]
       if l < 0 {
           l = 0
       }
       rr := i + r2[i]
       if rr >= M {
           rr = M - 1
       }
       L0[i] = int32(l)
       R0[i] = int32(rr)
   }
   L = append(L, L0)
   R = append(R, R0)
   // build doubling tables
   for k := 1; k < K; k++ {
       prevL := L[k-1]
       prevR := R[k-1]
       st := newSegTree(M)
       st.build(prevL, prevR)
       curL := make([]int32, M)
       curR := make([]int32, M)
       for i := 0; i < M; i++ {
           l := int(prevL[i])
           rr := int(prevR[i])
           curL[i] = st.queryMin(l, rr)
           curR[i] = st.queryMax(l, rr)
       }
       L = append(L, curL)
       R = append(R, curR)
   }
   // initialize current intervals and answers for each start
   curLpos := make([]int, n)
   curRpos := make([]int, n)
   ans := make([]int, n)
   for i := 0; i < n; i++ {
       pos := i + n
       curLpos[i] = pos
       curRpos[i] = pos
       ans[i] = 0
   }
   // binary lifting from highest bit
   for k := K - 1; k >= 0; k-- {
       st := newSegTree(M)
       st.build(L[k], R[k])
       step := 1 << k
       for i := 0; i < n; i++ {
           if curRpos[i]-curLpos[i]+1 >= n {
               continue
           }
           l := curLpos[i]
           rr := curRpos[i]
           nl := int(st.queryMin(l, rr))
           nr := int(st.queryMax(l, rr))
           if nr-nl+1 < n {
               curLpos[i] = nl
               curRpos[i] = nr
               ans[i] += step
           }
       }
   }
   // one final step if needed
   for i := 0; i < n; i++ {
       if curRpos[i]-curLpos[i]+1 < n {
           ans[i]++
       }
   }
   // output answer
   for i := 0; i < n; i++ {
       if i > 0 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, ans[i])
   }
   out.WriteByte('\n')
}
