package main

import (
   "bufio"
   "fmt"
   "os"
)

// segment tree for range add and range max query
type SegTree struct {
   n    int
   max  []int64
   lazy []int64
}

func NewSegTree(arr []int64) *SegTree {
   n := len(arr)
   size := 1
   for size < n {
       size <<= 1
   }
   maxArr := make([]int64, 2*size)
   lazy := make([]int64, 2*size)
   st := &SegTree{n: size, max: maxArr, lazy: lazy}
   // build leaves
   for i := 0; i < n; i++ {
       st.max[size+i] = arr[i]
   }
   // build internal
   for i := size - 1; i >= 1; i-- {
       left, right := st.max[2*i], st.max[2*i+1]
       if left > right {
           st.max[i] = left
       } else {
           st.max[i] = right
       }
   }
   return st
}

func (st *SegTree) apply(node int, v int64) {
   st.max[node] += v
   st.lazy[node] += v
}

func (st *SegTree) push(node int) {
   if st.lazy[node] != 0 {
       st.apply(2*node, st.lazy[node])
       st.apply(2*node+1, st.lazy[node])
       st.lazy[node] = 0
   }
}

// update adds v to [l..r]
func (st *SegTree) update(node, nl, nr, l, r int, v int64) {
   if r < nl || nr < l {
       return
   }
   if l <= nl && nr <= r {
       st.apply(node, v)
       return
   }
   st.push(node)
   mid := (nl + nr) >> 1
   st.update(2*node, nl, mid, l, r, v)
   st.update(2*node+1, mid+1, nr, l, r, v)
   // pull
   if st.max[2*node] > st.max[2*node+1] {
       st.max[node] = st.max[2*node]
   } else {
       st.max[node] = st.max[2*node+1]
   }
}

// Max returns max over all
func (st *SegTree) Max() int64 {
   return st.max[1]
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, m, k int
   fmt.Fscan(reader, &n, &m, &k)
   w := make([][]int64, n)
   for i := 0; i < n; i++ {
       w[i] = make([]int64, m)
       for j := 0; j < m; j++ {
           fmt.Fscan(reader, &w[i][j])
       }
   }
   // prefix sums per day
   pref := make([][]int64, n)
   for i := 0; i < n; i++ {
       pref[i] = make([]int64, m+1)
       for j := 0; j < m; j++ {
           pref[i][j+1] = pref[i][j] + w[i][j]
       }
   }
   M := m - k + 1
   dpPrev := make([]int64, M)
   // initial dp for day 0 (segment starting at day 1)
   if n == 1 {
       for x := 0; x < M; x++ {
           dpPrev[x] = pref[0][x+k] - pref[0][x]
       }
   } else {
       for x := 0; x < M; x++ {
           dpPrev[x] = (pref[0][x+k] - pref[0][x]) + (pref[1][x+k] - pref[1][x])
       }
   }
   // DP over segments i=1..n-1
   for i := 1; i < n; i++ {
       // build initial array with overlap at x=0
       arr := make([]int64, M)
       // pref[i] is weights for overlap day = i
       for y := 0; y < M; y++ {
           var overlap int64
           if y < k {
               overlap = pref[i][k] - pref[i][y]
           }
           arr[y] = dpPrev[y] - overlap
       }
       st := NewSegTree(arr)
       dpCurr := make([]int64, M)
       var nextPref []int64
       if i < n-1 {
           nextPref = pref[i+1]
       }
       for x := 0; x < M; x++ {
           best := st.Max()
           // compute val_i(x)
           var val int64
           if i < n-1 {
               val = (pref[i][x+k] - pref[i][x]) + (nextPref[x+k] - nextPref[x])
           } else {
               val = pref[i][x+k] - pref[i][x]
           }
           dpCurr[x] = best + val
           // update for next x
           if x == M-1 {
               break
           }
           jRem := x
           jAdd := x + k
           // remove jRem
           wRem := w[i][jRem]
           L1 := jRem - k + 1
           if L1 < 0 {
               L1 = 0
           }
           R1 := jRem
           if R1 >= M {
               R1 = M - 1
           }
           if L1 <= R1 && wRem != 0 {
               st.update(1, 0, st.n-1, L1, R1, wRem)
           }
           // add jAdd
           if jAdd < m {
               wAdd := w[i][jAdd]
               L2 := jAdd - k + 1
               if L2 < 0 {
                   L2 = 0
               }
               R2 := jAdd
               if R2 >= M {
                   R2 = M - 1
               }
               if L2 <= R2 && wAdd != 0 {
                   st.update(1, 0, st.n-1, L2, R2, -wAdd)
               }
           }
       }
       dpPrev = dpCurr
   }
   // answer
   var ans int64
   for _, v := range dpPrev {
       if v > ans {
           ans = v
       }
   }
   fmt.Fprint(writer, ans)
}
