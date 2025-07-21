package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // collect lucky positions and values
   b := make([]int, 0, 1005)
   v := make([]int, 0, 1005)
   for i, x := range a {
       if isLucky(x) {
           b = append(b, i+1)
           v = append(v, x)
       }
   }
   m := len(b)
   // prev occurrence in lucky list
   prev := make([]int, m)
   last := map[int]int{}
   for i := 0; i < m; i++ {
       if p, ok := last[v[i]]; ok {
           prev[i] = p + 1 // 1-based
       } else {
           prev[i] = 0
       }
       last[v[i]] = i
   }
   // compute blocks c[0..m]
   c := make([]int, m+1)
   if m == 0 {
       // all pure segments
       // count pairs of non-intersecting subarrays
       // total segments = n*(n+1)/2, but we need pairs non-intersecting
       // formula sum_{r=1..n-1} r*(n-r)*(n-r+1)/2
       var ans int64
       for r := 1; r < n; r++ {
           cnt1 := int64(r)
           cnt2 := int64(n-r)
           ans += cnt1 * cnt2 * (cnt2 + 1) / 2
       }
       fmt.Println(ans)
       return
   }
   c[0] = b[0] - 1
   for i := 1; i < m; i++ {
       c[i] = b[i] - b[i-1] - 1
   }
   c[m] = n - b[m-1]
   // A[1..m], B[1..m]
   A := make([]int64, m+1)
   B := make([]int64, m+1)
   for i := 1; i <= m; i++ {
       A[i] = int64(c[i-1] + 1)
       B[i] = int64(c[i] + 1)
   }
   // PA and S1prefix
   PA := make([]int64, m+1)
   S1 := make([]int64, m+1)
   for i := 1; i <= m; i++ {
       PA[i] = PA[i-1] + A[i]
       S1[i] = S1[i-1] + PA[i]*B[i]
   }
   // SB, E, SE
   SB := make([]int64, m+2)
   for i := m; i >= 1; i-- {
       SB[i] = SB[i+1] + B[i]
   }
   E := make([]int64, m+2)
   for i := 1; i <= m; i++ {
       E[i] = A[i] * SB[i]
   }
   SE := make([]int64, m+3)
   for i := m; i >= 1; i-- {
       SE[i] = SE[i+1] + E[i]
   }
   // Case1: both include lucky
   var ans int64
   for i2 := 1; i2 <= m; i2++ {
       minPrev := m + 1
       for j2 := i2; j2 <= m; j2++ {
           // prev index for j2: 1-based if exists
           if prev[j2-1] > 0 && prev[j2-1] < minPrev {
               minPrev = prev[j2-1]
           }
           // if no prev, prev=0 -> minPrev stays
           J := minPrev - 1
           var s1 int64
           if J >= 1 {
               s1 = S1[J]
           }
           ans += s1 * (A[i2] * B[j2])
       }
   }
   // pure segment counts per block
   pureC := make([]int64, m+1)
   for k := 0; k <= m; k++ {
       ck := int64(c[k])
       pureC[k] = ck * (ck + 1) / 2
   }
   // Case2: pure-pure
   // internal in same block
   for k := 0; k <= m; k++ {
       ck := c[k]
       // sum r1=1..ck-1 of (r1*(r1+1)/2)*((ck-r1)*(ck-r1+1)/2)
       for r1 := 1; r1 < ck; r1++ {
           left := int64(r1) * int64(r1+1) / 2
           rem := int64(ck-r1)
           right := rem * (rem + 1) / 2
           ans += left * right
       }
   }
   // inter-block k<l
   var sumPrev int64
   for k := 0; k <= m; k++ {
       ans += sumPrev * pureC[k]
       sumPrev += pureC[k]
   }
   // Case3: seg1 pure, seg2 include lucky
   for k := 0; k <= m; k++ {
       idx := k + 2
       if idx <= m {
           ans += pureC[k] * SE[idx]
       }
   }
   // Case4: seg1 include lucky, seg2 pure
   for k := 0; k <= m; k++ {
       if k-1 >= 1 {
           ans += pureC[k] * S1[k-1]
       }
   }
   fmt.Println(ans)
}

func isLucky(x int) bool {
   if x <= 0 {
       return false
   }
   for x > 0 {
       d := x % 10
       if d != 4 && d != 7 {
           return false
       }
       x /= 10
   }
   return true
}
