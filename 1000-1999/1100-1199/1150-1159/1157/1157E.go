package main

import (
   "bufio"
   "fmt"
   "os"
)

// Fenwick tree (Binary Indexed Tree) for prefix sums over ints.
type Fenwick struct {
   n    int
   tree []int
}

// NewFenwick creates a Fenwick tree for n elements (0..n-1).
func NewFenwick(n int) *Fenwick {
   return &Fenwick{n: n, tree: make([]int, n+1)}
}

// Add adds delta to element at index i (0-based).
func (f *Fenwick) Add(i, delta int) {
   // internal indices are 1-based
   for j := i + 1; j <= f.n; j += j & -j {
       f.tree[j] += delta
   }
}

// Sum returns the sum of elements in [0..i] (0-based). If i < 0, returns 0.
func (f *Fenwick) Sum(i int) int {
   if i < 0 {
       return 0
   }
   res := 0
   for j := i + 1; j > 0; j -= j & -j {
       res += f.tree[j]
   }
   return res
}

// FindKth finds the smallest index idx (0-based) such that the prefix sum [0..idx] >= k (1-based k).
// Requires 1 <= k <= total sum.
func (f *Fenwick) FindKth(k int) int {
   idx := 0
   // compute highest power of two <= n
   bit := 1
   for bit<<1 <= f.n {
       bit <<= 1
   }
   for ; bit > 0; bit >>= 1 {
       nxt := idx + bit
       if nxt <= f.n && f.tree[nxt] < k {
           idx = nxt
           k -= f.tree[nxt]
       }
   }
   // idx is the largest position with prefix sum < original k, so result is idx (0-based)
   return idx
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   fenw := NewFenwick(n)
   for i := 0; i < n; i++ {
       var x int
       fmt.Fscan(reader, &x)
       fenw.Add(x, 1)
   }
   c := make([]int, n)
   for i := 0; i < n; i++ {
       want := (n - a[i]) % n
       // count of elements < want
       s := fenw.Sum(want - 1)
       total := fenw.Sum(n - 1)
       var k int
       if total-s > 0 {
           k = s + 1
       } else {
           k = 1
       }
       idx := fenw.FindKth(k)
       c[i] = (a[i] + idx) % n
       fenw.Add(idx, -1)
   }
   for i := 0; i < n; i++ {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, c[i])
   }
   writer.WriteByte('\n')
}
