package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// BIT implements a Fenwick tree for prefix sums
type BIT struct {
   n    int
   tree []int
}

// NewBIT creates a BIT for size n (1-based indexing)
func NewBIT(n int) *BIT {
   return &BIT{n: n, tree: make([]int, n+1)}
}

// Update adds v at position i (1-based)
func (b *BIT) Update(i, v int) {
   for ; i <= b.n; i += i & -i {
       b.tree[i] += v
   }
}

// Sum returns the prefix sum [1..i]
func (b *BIT) Sum(i int) int {
   s := 0
   for ; i > 0; i -= i & -i {
       s += b.tree[i]
   }
   return s
}

// FindByPrefix finds smallest i such that sum[1..i] >= k; assumes all values non-negative and k >= 1
func (b *BIT) FindByPrefix(k int) int {
   idx := 0
   bitMask := 1
   for bitMask<<1 <= b.n {
       bitMask <<= 1
   }
   for bitMask > 0 {
       next := idx + bitMask
       if next <= b.n && b.tree[next] < k {
           idx = next
           k -= b.tree[next]
       }
       bitMask >>= 1
   }
   return idx + 1
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   A := make([]int, n)
   B := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &A[i])
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &B[i])
   }
   // sort A ascending
   sort.Ints(A)
   // build BIT for B counts
   bit := NewBIT(m)
   for _, b := range B {
       // b in [0..m-1], map to idx b+1
       bit.Update(b+1, 1)
   }
   C := make([]int, n)
   for i, a := range A {
       // bound for non-wrap: b <= m-1-a
       bound := m - 1 - a
       var idxB int
       if bound >= 0 {
           // sum of counts up to bound
           cnt := bit.Sum(min(bound, m-1) + 1)
           if cnt > 0 {
               // find b index
               idx := bit.FindByPrefix(cnt)
               idxB = idx - 1
           } else {
               // wrap: pick largest b available
               total := bit.Sum(m)
               idx := bit.FindByPrefix(total)
               idxB = idx - 1
           }
       } else {
           // no non-wrap possible, pick largest b
           total := bit.Sum(m)
           idx := bit.FindByPrefix(total)
           idxB = idx - 1
       }
       // remove b
       bit.Update(idxB+1, -1)
       // compute C
       sum := a + idxB
       if sum >= m {
           sum -= m
       }
       C[i] = sum
   }
   // sort C descending
   sort.Sort(sort.Reverse(sort.IntSlice(C)))
   // output
   for i, v := range C {
       if i > 0 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, v)
   }
   out.WriteByte('\n')
}
