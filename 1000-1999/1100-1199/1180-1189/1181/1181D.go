package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// BIT implements a Fenwick tree for prefix sums and k-th order statistic
type BIT struct {
   n    int
   tree []int
   logn int
}

// NewBIT initializes a BIT of size n
func NewBIT(n int) *BIT {
   bit := &BIT{n: n, tree: make([]int, n+1)}
   // compute highest power of two <= n
   lg := 1
   for (lg << 1) <= n {
       lg <<= 1
   }
   bit.logn = lg
   return bit
}

// Update adds v at index i (1-based)
func (b *BIT) Update(i, v int) {
   for x := i; x <= b.n; x += x & -x {
       b.tree[x] += v
   }
}

// Sum returns prefix sum up to i (1-based)
func (b *BIT) Sum(i int) int {
   s := 0
   for x := i; x > 0; x -= x & -x {
       s += b.tree[x]
   }
   return s
}

// Kth finds smallest index i such that Sum(i) >= k
func (b *BIT) Kth(k int) int {
   pos := 0
   bitMask := b.logn
   for bitMask > 0 {
       nxt := pos + bitMask
       if nxt <= b.n && b.tree[nxt] < k {
           k -= b.tree[nxt]
           pos = nxt
       }
       bitMask >>= 1
   }
   return pos + 1
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m, q int
   fmt.Fscan(reader, &n, &m, &q)
   a := make([]int, n)
   ci := make([]int, m+1)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       ci[a[i]]++
   }
   ks := make([]uint64, q)
   for i := 0; i < q; i++ {
       fmt.Fscan(reader, &ks[i])
   }
   // initial answer array
   ans := make([]int, q)
   // compute max count
   C := 0
   for i := 1; i <= m; i++ {
       if ci[i] > C {
           C = ci[i]
       }
   }
   // count cities by initial count
   cntCi := make([]int, C+1)
   for i := 1; i <= m; i++ {
       cntCi[ci[i]]++
   }
   // prefix sum of cntCi to get number of cities with ci <= h
   prefCi := make([]int, C+1)
   prefCi[0] = cntCi[0]
   for h := 1; h <= C; h++ {
       prefCi[h] = prefCi[h-1] + cntCi[h]
   }
   // levels[h] = number of cities with ci <= h
   levels := make([]uint64, C)
   for h := 0; h < C; h++ {
       levels[h] = uint64(prefCi[h])
   }
   // prefix sum of levels
   prefLevels := make([]uint64, C+1)
   prefLevels[0] = 0
   for i := 1; i <= C; i++ {
       prefLevels[i] = prefLevels[i-1] + levels[i-1]
   }
   D := prefLevels[C]
   // prepare buckets for leveling queries
   type qSpec struct{ idx, offset int }
   buckets := make([][]qSpec, C)
   // answer cyclic and direct queries
   for i := 0; i < q; i++ {
       k := ks[i]
       if k <= uint64(n) {
           ans[i] = a[k-1]
       } else {
           t := k - uint64(n)
           if uint64(C) == 0 || t > D {
               // cyclic part
               v := (t - D - 1) % uint64(m)
               ans[i] = int(v) + 1
           } else {
               // leveling part: find layer
               x := sort.Search(len(prefLevels), func(x int) bool {
                   return prefLevels[x] >= t
               })
               h := x - 1
               offset := int(t - prefLevels[x-1])
               buckets[h] = append(buckets[h], qSpec{i, offset})
           }
       }
   }
   // bucket cities by initial count
   ciBuckets := make([][]int, C)
   for i := 1; i <= m; i++ {
       if ci[i] < C {
           ciBuckets[ci[i]] = append(ciBuckets[ci[i]], i)
       }
   }
   // process leveling with BIT
   bit := NewBIT(m)
   for h := 0; h < C; h++ {
       // add cities with ci == h
       for _, city := range ciBuckets[h] {
           bit.Update(city, 1)
       }
       // answer queries at this layer
       for _, qs := range buckets[h] {
           ans[qs.idx] = bit.Kth(qs.offset)
       }
   }
   // output answers
   for i := 0; i < q; i++ {
       fmt.Fprintln(writer, ans[i])
   }
}
