package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// BIT implements a Fenwick tree for prefix sums and find_kth
type BIT struct {
   n int
   t []int
}

// NewBIT creates a BIT of size n
func NewBIT(n int) *BIT {
   return &BIT{n: n, t: make([]int, n+1)}
}

// Add adds v at position i (1-based)
func (b *BIT) Add(i, v int) {
   for ; i <= b.n; i += i & -i {
       b.t[i] += v
   }
}

// Sum returns prefix sum up to i (1-based)
func (b *BIT) Sum(i int) int {
   s := 0
   for ; i > 0; i -= i & -i {
       s += b.t[i]
   }
   return s
}

// FindKth returns the smallest index pos such that Sum(pos) >= k
func (b *BIT) FindKth(k int) int {
   pos := 0
   // largest power of two <= n
   bitMask := 1
   for bitMask<<1 <= b.n {
       bitMask <<= 1
   }
   for d := bitMask; d > 0; d >>= 1 {
       np := pos + d
       if np <= b.n && b.t[np] < k {
           k -= b.t[np]
           pos = np
       }
   }
   return pos + 1
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   a := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // process events
   bits := make([]byte, n+1)
   alive := make([]bool, n+1)
   bit := NewBIT(n)
   total := 0
   for i := 0; i < n; i++ {
       var x int
       fmt.Fscan(reader, &x)
       if x >= 0 {
           total++
           if x == 0 {
               bits[total] = '0'
           } else {
               bits[total] = '1'
           }
           alive[total] = true
           bit.Add(total, 1)
       } else {
           // hit event
           L := bit.Sum(total)
           // find k = number of a[j] <= L
           k := sort.SearchInts(a, L+1)
           // remove for j = 0..k-1
           for j := 0; j < k; j++ {
               // adjusted position in current sequence
               pos := a[j] - j
               if pos <= 0 || pos > L {
                   break
               }
               idx := bit.FindKth(pos)
               // remove idx
               bit.Add(idx, -1)
               alive[idx] = false
           }
       }
   }
   // output remaining sequence
   var has bool
   for i := 1; i <= total; i++ {
       if alive[i] {
           writer.WriteByte(bits[i])
           has = true
       }
   }
   if !has {
       writer.WriteString("Poor stack!")
   }
}
