package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// BIT implements a Fenwick tree for int64 values.
type BIT struct {
   n int
   t []int64
}

// NewBIT creates a BIT of size n.
func NewBIT(n int) *BIT {
   return &BIT{n, make([]int64, n+1)}
}

// Add adds value v at position i (1-based).
func (b *BIT) Add(i int, v int64) {
   for ; i <= b.n; i += i & -i {
       b.t[i] += v
   }
}

// Sum returns the prefix sum of [1..i].
func (b *BIT) Sum(i int) int64 {
   var s int64
   for ; i > 0; i -= i & -i {
       s += b.t[i]
   }
   return s
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   pairs := make([][2]int64, n)
   all := make([]int64, 0, 2*n)
   for i := 0; i < n; i++ {
       var a, b int64
       fmt.Fscan(in, &a, &b)
       pairs[i][0], pairs[i][1] = a, b
       all = append(all, a, b)
   }
   // sort and unique positions
   sort.Slice(all, func(i, j int) bool { return all[i] < all[j] })
   uniq := all[:0]
   for i, v := range all {
       if i == 0 || v != all[i-1] {
           uniq = append(uniq, v)
       }
   }
   all = uniq
   k := len(all)
   // map position to index
   pos2idx := make(map[int64]int, k)
   for i, v := range all {
       pos2idx[v] = i
   }
   // initialize mapping M: at each position index i, value index M[i]
   M := make([]int, k)
   for i := 0; i < k; i++ {
       M[i] = i
   }
   // apply swaps
   for _, p := range pairs {
       ia := pos2idx[p[0]]
       ib := pos2idx[p[1]]
       M[ia], M[ib] = M[ib], M[ia]
   }
   // count inversions among swapped positions
   bit := NewBIT(k)
   var invSS int64
   for i := 0; i < k; i++ {
       v := M[i]
       seen := int64(i)
       // number of previous values <= v
       sum := bit.Sum(v + 1)
       // those > v are inversions
       invSS += seen - sum
       bit.Add(v+1, 1)
   }
   var sumA, sumB int64
   // count inversions between swapped and unswapped
   for i := 0; i < k; i++ {
       pi := all[i]
       vi := all[M[i]]
       // case: i in S, j not in S, i<j, value greater
       if vi > pi+1 {
           a := pi + 1
           b := vi - 1
           tot := b - a + 1
           l := sort.Search(k, func(j int) bool { return all[j] >= a })
           r := sort.Search(k, func(j int) bool { return all[j] > b }) - 1
           var cnt int64
           if r >= l {
               cnt = int64(r - l + 1)
           }
           sumA += tot - cnt
       }
       // case: i not in S, j in S, i<j, unswapped value greater
       if vi < pi-1 {
           a := vi + 1
           b := pi - 1
           tot := b - a + 1
           l := sort.Search(k, func(j int) bool { return all[j] >= a })
           r := sort.Search(k, func(j int) bool { return all[j] > b }) - 1
           var cnt int64
           if r >= l {
               cnt = int64(r - l + 1)
           }
           sumB += tot - cnt
       }
   }
   ans := invSS + sumA + sumB
   fmt.Println(ans)
}
