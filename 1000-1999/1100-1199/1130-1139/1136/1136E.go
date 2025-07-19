package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// BIT implements a Fenwick tree for sum queries
type BIT struct {
   n int
   t []int64
}

// NewBIT creates a BIT for indices 1..n
func NewBIT(n int) *BIT {
   return &BIT{n: n, t: make([]int64, n+1)}
}

// Add adds v at index i (1-based)
func (b *BIT) Add(i int, v int64) {
   for ; i <= b.n; i += i & -i {
       b.t[i] += v
   }
}

// Sum returns prefix sum up to i (1-based)
func (b *BIT) Sum(i int) int64 {
   var s int64
   for ; i > 0; i -= i & -i {
       s += b.t[i]
   }
   return s
}

// getv returns sum of bit in [a..b] where a,b are 0-based
func getv(b *BIT, a, bIdx int) int64 {
   if bIdx < a {
       return 0
   }
   return b.Sum(bIdx+1) - b.Sum(a)
}

// set applies delta v at 0-based index idx for two BITs
func set(bit1, bit2 *BIT, idx int, v int64) {
   bit1.Add(idx+1, v)
   bit2.Add(idx+1, int64(idx)*v)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var N int
   fmt.Fscan(reader, &N)
   a := make([]int64, N+1)
   k := make([]int64, N+1)
   for i := 1; i <= N; i++ {
       fmt.Fscan(reader, &a[i])
   }
   for i := 1; i < N; i++ {
       fmt.Fscan(reader, &k[i])
   }
   // v1 tracks surplus at each boundary
   v1 := make([]int64, N+1)
   // BIT size N+1 to allow updates at idx N
   bit1 := NewBIT(N + 1)
   bit2 := NewBIT(N + 1)
   // sorted slice of indices with nonzero v1
   idxs := make([]int, 0, N)
   // initialize
   for i := 0; i < N; i++ {
       var ki int64
       if i >= 1 {
           ki = k[i]
       }
       v1[i] = a[i+1] - a[i] - ki
       diff := a[i+1] - a[i]
       set(bit1, bit2, i, diff)
       if v1[i] != 0 {
           idxs = append(idxs, i)
       }
   }
   sort.Ints(idxs)
   var Q int
   fmt.Fscan(reader, &Q)
   for Q > 0 {
       Q--
       var op byte
       var x1, x2 int
       fmt.Fscan(reader, &op, &x1, &x2)
       if op == 's' {
           l, r := x1, x2
           ar := getv(bit1, 0, r-1)
           al := getv(bit1, 0, l-1)
           mid := getv(bit2, l, r-1)
           res := int64(r)*ar - mid - int64(l-1)*al
           fmt.Fprintln(writer, res)
       } else {
           id := x1
           delta := int64(x2)
           if delta == 0 {
               continue
           }
           // increase v1[id-1]
           i0 := id - 1
           pos := sort.SearchInts(idxs, i0)
           if v1[i0] == 0 {
               idxs = append(idxs, 0)
               copy(idxs[pos+1:], idxs[pos:])
               idxs[pos] = i0
           }
           v1[i0] += delta
           set(bit1, bit2, i0, delta)
           // decrease v1[id]
           pos = sort.SearchInts(idxs, id)
           if id <= N && v1[id] == 0 {
               idxs = append(idxs, 0)
               copy(idxs[pos+1:], idxs[pos:])
               idxs[pos] = id
           }
           if id <= N {
               v1[id] -= delta
               set(bit1, bit2, id, -delta)
               if v1[id] > 0 {
                   continue
               }
               if v1[id] == 0 {
                   if pos < len(idxs) && idxs[pos] == id {
                       idxs = append(idxs[:pos], idxs[pos+1:]...)
                   }
                   continue
               }
               // v1[id] < 0
               remain := -v1[id]
               // adjust by removing full deficit
               set(bit1, bit2, id, -v1[id])
               v1[id] = 0
               if pos < len(idxs) && idxs[pos] == id {
                   idxs = append(idxs[:pos], idxs[pos+1:]...)
               }
               // propagate to next indices
               it := sort.SearchInts(idxs, id)
               for remain > 0 && it < len(idxs) {
                   i := idxs[it]
                   v := v1[i]
                   if v > remain {
                       v = remain
                   }
                   remain -= v
                   v1[i] -= v
                   set(bit1, bit2, i, -v)
                   if v1[i] == 0 {
                       idxs = append(idxs[:it], idxs[it+1:]...)
                   } else {
                       it++
                   }
               }
           }
       }
   }
}
