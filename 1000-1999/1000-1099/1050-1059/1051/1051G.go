package main

import (
   "bufio"
   "fmt"
   "os"
)

// BIT supports point updates and prefix sums
type BIT struct {
   n    int
   tree []int64
}

// NewBIT initializes a BIT for size n
func NewBIT(n int) *BIT {
   return &BIT{n: n, tree: make([]int64, n+1)}
}

// Add increments index i by value v
func (b *BIT) Add(i int, v int64) {
   for x := i; x <= b.n; x += x & -x {
       b.tree[x] += v
   }
}

// Sum returns the prefix sum up to i
func (b *BIT) Sum(i int) int64 {
   var s int64
   for x := i; x > 0; x -= x & -x {
       s += b.tree[x]
   }
   return s
}

// RangeSum returns sum of [l..r]
func (b *BIT) RangeSum(l, r int) int64 {
   if l > r {
       return 0
   }
   return b.Sum(r) - b.Sum(l-1)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   fmt.Fscan(in, &n)
   // b_i in [1..n]
   bitCnt := NewBIT(n)
   bitSum := NewBIT(n)
   var totalB, S, sumAb int64
   low := int64(1e18)
   res := make([]int64, n)
   for i := 0; i < n; i++ {
       var a, b int
       fmt.Fscan(in, &a, &b)
       if int64(a) < low {
           low = int64(a)
       }
       // count of b_j > b
       cntGt := bitCnt.RangeSum(b+1, n)
       sumGt := bitSum.RangeSum(b+1, n)
       sumLower := totalB - sumGt
       // update S and totalB
       S += int64(b)*cntGt + sumLower
       totalB += int64(b)
       bitCnt.Add(b, 1)
       bitSum.Add(b, int64(b))
       sumAb += int64(b) * int64(a)
       // f = low*totalB + S - sumAb
       res[i] = low*totalB + S - sumAb
   }
   for i, v := range res {
       if i > 0 {
           fmt.Fprint(out, " ")
       }
       fmt.Fprint(out, v)
   }
   fmt.Fprintln(out)
}
