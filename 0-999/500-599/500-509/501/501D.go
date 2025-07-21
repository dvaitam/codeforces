package main

import (
   "bufio"
   "fmt"
   "os"
)

// Fenwick implements a Fenwick tree (BIT) for point updates and prefix sums.
type Fenwick struct {
   n    int
   tree []int
}

// NewFenwick creates a Fenwick tree of size n (1-based internally).
func NewFenwick(n int) *Fenwick {
   return &Fenwick{n: n, tree: make([]int, n+1)}
}

// Add adds delta at position i (1-based).
func (f *Fenwick) Add(i, delta int) {
   for ; i <= f.n; i += i & -i {
       f.tree[i] += delta
   }
}

// Sum returns the prefix sum [1..i] (1-based).
func (f *Fenwick) Sum(i int) int {
   s := 0
   for ; i > 0; i -= i & -i {
       s += f.tree[i]
   }
   return s
}

// FindByOrder finds the smallest index idx (1-based) such that Sum(idx) > order-1
// i.e., the (order)-th one, where order is 1-based.
func (f *Fenwick) FindByOrder(order int) int {
   idx := 0
   bitMask := 1
   for bitMask<<1 <= f.n {
       bitMask <<= 1
   }
   // binary search on tree
   for bit := bitMask; bit > 0; bit >>= 1 {
       next := idx + bit
       if next <= f.n && f.tree[next] < order {
           order -= f.tree[next]
           idx = next
       }
   }
   return idx + 1
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   p := make([]int, n)
   q := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &p[i])
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &q[i])
   }
   // compute Lehmer codes c and d
   c := make([]int, n)
   d := make([]int, n)
   bitP := NewFenwick(n)
   bitQ := NewFenwick(n)
   // initialize BITs with 1s
   for i := 1; i <= n; i++ {
       bitP.Add(i, 1)
       bitQ.Add(i, 1)
   }
   for i := 0; i < n; i++ {
       // for p
       x := p[i] + 1
       c[i] = bitP.Sum(x - 1)
       bitP.Add(x, -1)
       // for q
       y := q[i] + 1
       d[i] = bitQ.Sum(y - 1)
       bitQ.Add(y, -1)
   }
   // add factorial base digits
   e := make([]int, n)
   carry := 0
   for i := n - 1; i >= 0; i-- {
       base := n - i
       sum := c[i] + d[i] + carry
       e[i] = sum % base
       carry = sum / base
   }
   // reconstruct permutation from factorial code e
   bitR := NewFenwick(n)
   for i := 1; i <= n; i++ {
       bitR.Add(i, 1)
   }
   res := make([]int, n)
   for i := 0; i < n; i++ {
       // e[i] is 0-based order among remaining
       idx := bitR.FindByOrder(e[i] + 1)
       res[i] = idx - 1
       bitR.Add(idx, -1)
   }
   // output
   for i, v := range res {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   writer.WriteByte('\n')
}
