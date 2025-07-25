package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 998244353

// BIT supports prefix sum and lower_bound on prefix sums (first index with sum >= target)
type BIT struct {
   n      int
   tree   []int
   maxPow int
}

func NewBIT(n int) *BIT {
   b := &BIT{n: n, tree: make([]int, n+1)}
   // compute highest power of two <= n
   pow := 1
   for pow<<1 <= n {
       pow <<= 1
   }
   b.maxPow = pow
   return b
}

// Add value v at index i (0-based)
func (b *BIT) Add(i, v int) {
   i++
   for ; i <= b.n; i += i & -i {
       b.tree[i] += v
   }
}

// Sum returns prefix sum [0..i] inclusive
func (b *BIT) Sum(i int) int {
   if i < 0 {
       return 0
   }
   if i >= b.n {
       i = b.n - 1
   }
   i++
   s := 0
   for ; i > 0; i -= i & -i {
       s += b.tree[i]
   }
   return s
}

// LowerBound returns smallest index p such that Sum(p) >= target
func (b *BIT) LowerBound(target int) int {
   idx := 0
   bitMask := b.maxPow
   for bitMask > 0 {
       t := idx + bitMask
       if t <= b.n && b.tree[t] < target {
           idx = t
           target -= b.tree[t]
       }
       bitMask >>= 1
   }
   // idx is internal tree index with prefix sum < target
   // element index = idx (0-based)
   return idx
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m int
   fmt.Fscan(in, &n, &m)
   c := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &c[i])
   }
   // positions for each color
   pos := make([][]int, n+1)
   for i, v := range c {
       if v >= 1 && v <= n {
           pos[v] = append(pos[v], i+1) // segments numbered 1..m
       }
   }
   // compute L and R for each color
   L := make([]int, n+1)
   R := make([]int, n+1)
   for i := 1; i <= n; i++ {
       if len(pos[i]) == 0 {
           fmt.Fprintln(out, 0)
           return
       }
       L[i] = pos[i][0]
       R[i] = pos[i][len(pos[i])-1]
   }
   // BIT over positions 0..m+1 (size m+2)
   size := m + 2
   bit := NewBIT(size)
   // add sentinels at 0 and m+1
   bit.Add(0, 1)
   bit.Add(m+1, 1)
   ans := 1
   for i := 1; i <= n; i++ {
       li := L[i]
       ri := R[i]
       // find last blocker before li
       leftCount := bit.Sum(li - 1)
       lpos := bit.LowerBound(leftCount)
       // ensure li and ri are in the same good component (no blocker between them)
       leftCountRi := bit.Sum(ri - 1)
       lposRi := bit.LowerBound(leftCountRi)
       if lposRi != lpos {
           fmt.Fprintln(out, 0)
           return
       }
       // find first blocker after ri
       rightCount := bit.Sum(ri)
       rpos := bit.LowerBound(rightCount + 1)
       // compute valid a_i in [boundL..li] and b_i in [ri..boundR]
       boundL := lpos + 1
       boundR := rpos - 1
       ways := (li - boundL + 1) * (boundR - ri + 1) % mod
       if ways < 0 {
           ways += mod
       }
       ans = int(int64(ans) * int64(ways) % mod)
       // insert all positions of color i as blockers for next
       for _, p := range pos[i] {
           bit.Add(p, 1)
       }
   }
   fmt.Fprintln(out, ans)
}
