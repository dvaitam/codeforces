package main

import (
   "bufio"
   "fmt"
   "os"
)

// BIT implements a Fenwick tree for prefix sums and order-statistic
type BIT struct {
   n    int
   tree []int
}

// NewBIT creates a BIT of size n (1-based indexing)
func NewBIT(n int) *BIT {
   return &BIT{n: n, tree: make([]int, n+1)}
}

// Add adds delta at position i (1-based)
func (b *BIT) Add(i, delta int) {
   for x := i; x <= b.n; x += x & -x {
       b.tree[x] += delta
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

// FindKth returns smallest i such that Sum(i) >= k, assumes all values non-negative and total sum >= k
func (b *BIT) FindKth(k int) int {
   pos := 0
   bitMask := 1
   // compute highest power of two >= n
   for bitMask<<1 <= b.n {
       bitMask <<= 1
   }
   for d := bitMask; d > 0; d >>= 1 {
       nxt := pos + d
       if nxt <= b.n && b.tree[nxt] < k {
           pos = nxt
           k -= b.tree[nxt]
       }
   }
   return pos + 1
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var k, n int
   var s string
   fmt.Fscan(reader, &k)
   fmt.Fscan(reader, &s)
   fmt.Fscan(reader, &n)
   // total length
   m := len(s)
   total := k * m
   // positions per character
   pos := make([][]int, 26)
   for i := 0; i < total; i++ {
       c := s[i%m] - 'a'
       pos[c] = append(pos[c], i)
   }
   // BIT per character
   bits := make([]*BIT, 26)
   for c := 0; c < 26; c++ {
       bits[c] = NewBIT(len(pos[c]))
       // initialize all occurrences as present
       for i := 1; i <= len(pos[c]); i++ {
           bits[c].Add(i, 1)
       }
   }
   // deleted markers
   deleted := make([]bool, total)
   // process operations
   for i := 0; i < n; i++ {
       var p int
       var cs string
       fmt.Fscan(reader, &p, &cs)
       cidx := int(cs[0] - 'a')
       // find p-th alive occurrence in pos[cidx]
       idx := bits[cidx].FindKth(p) - 1 // zero-based index in pos[cidx]
       gpos := pos[cidx][idx]
       deleted[gpos] = true
       bits[cidx].Add(idx+1, -1)
   }
   // build result
   var res []byte
   res = make([]byte, 0, total)
   for i := 0; i < total; i++ {
       if !deleted[i] {
           res = append(res, s[i%m])
       }
   }
   writer.Write(res)
}
