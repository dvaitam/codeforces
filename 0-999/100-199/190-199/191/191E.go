package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Fenwick tree for frequencies
type BIT struct {
   n    int
   tree []int
}

func NewBIT(n int) *BIT {
   return &BIT{n: n, tree: make([]int, n+1)}
}

func (b *BIT) Update(i int, v int) {
   for ; i <= b.n; i += i & -i {
       b.tree[i] += v
   }
}

// Query sum of [1..i]
func (b *BIT) Query(i int) int {
   s := 0
   for ; i > 0; i -= i & -i {
       s += b.tree[i]
   }
   return s
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   var k int64
   fmt.Fscan(in, &n, &k)
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // prefix sums
   S := make([]int64, n+1)
   for i := 0; i < n; i++ {
       S[i+1] = S[i] + a[i]
   }
   // compress prefix sums
   P := make([]int64, n+1)
   copy(P, S)
   sort.Slice(P, func(i, j int) bool { return P[i] < P[j] })
   P = uniqueInt64(P)

   // count number of subarrays with sum >= x
   count := func(x int64) int64 {
       bit := NewBIT(len(P))
       var cnt int64
       for _, sj := range S {
           // number of si <= sj - x
           v := sj - x
           u := sort.Search(len(P), func(i int) bool { return P[i] > v })
           cnt += int64(bit.Query(u))
           // add current sj
           pos := sort.Search(len(P), func(i int) bool { return P[i] >= sj })
           bit.Update(pos+1, 1)
       }
       return cnt
   }

   // binary search for k-th largest sum
   var low, high int64 = -100000000000000, 100000000000000
   for low < high {
       mid := (low + high + 1) >> 1
       if count(mid) >= k {
           low = mid
       } else {
           high = mid - 1
       }
   }
   fmt.Fprintln(out, low)
}

// uniqueInt64 returns sorted unique slice
func uniqueInt64(a []int64) []int64 {
   if len(a) == 0 {
       return a
   }
   j := 0
   for i := 1; i < len(a); i++ {
       if a[i] != a[j] {
           j++
           a[j] = a[i]
       }
   }
   return a[:j+1]
}
