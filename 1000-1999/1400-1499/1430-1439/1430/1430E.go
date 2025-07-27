package main

import (
   "bufio"
   "fmt"
   "os"
)

// BIT implements a Fenwick tree for prefix sums.
type BIT struct {
   n    int
   tree []int
}

// NewBIT creates a BIT of size n (0..n-1 indices).
func NewBIT(n int) *BIT {
   return &BIT{n: n, tree: make([]int, n+1)}
}

// Add adds v at index i.
func (b *BIT) Add(i, v int) {
   // internal tree is 1-based
   for x := i + 1; x <= b.n; x += x & -x {
       b.tree[x] += v
   }
}

// Sum returns sum of [0..i]
func (b *BIT) Sum(i int) int {
   s := 0
   for x := i + 1; x > 0; x -= x & -x {
       s += b.tree[x]
   }
   return s
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   // Read n
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // Read string s
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   // positions of each character in s
   pos := make([][]int, 26)
   for i := 0; i < n; i++ {
       c := s[i] - 'a'
       pos[c] = append(pos[c], i)
   }
   // build target mapping P by scanning s in reverse
   P := make([]int, n)
   ptr := make([]int, 26)
   idx := 0
   for i := n - 1; i >= 0; i-- {
       c := s[i] - 'a'
       P[idx] = pos[c][ptr[c]]
       ptr[c]++
       idx++
   }
   // count inversions in P
   bit := NewBIT(n)
   var inv int64
   for i, v := range P {
       // count of previous <= v
       cnt := bit.Sum(v)
       // previous elements count = i, so greater = i - cnt
       inv += int64(i - cnt)
       bit.Add(v, 1)
   }
   fmt.Fprintln(writer, inv)
}
