package main

import (
   "bufio"
   "fmt"
   "math/big"
   "os"
   "sort"
)

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}

// calc builds dp and tr bitsets for subset sum to detect sums reachable and trace picks
func calc(n, sum int, arr []int) ([]*big.Int, []*big.Int) {
   dp := make([]*big.Int, n+1)
   tr := make([]*big.Int, n+1)
   dp[0] = big.NewInt(1)
   for i := 1; i <= n; i++ {
       dp[i] = new(big.Int)
       tr[i] = new(big.Int)
   }
   for i := 1; i <= n; i++ {
       // shift previous by arr[i]
       shift := new(big.Int).Lsh(dp[i-1], uint(arr[i]))
       // dp[i] = dp[i-1] OR shift
       dp[i].Or(dp[i-1], shift)
       // record shift positions for trace
       tr[i].Set(shift)
   }
   return dp, tr
}

// rev reconstructs choices: positive arr[i] if bit set in tr, else negative
func rev(n, sum int, arr []int, tr []*big.Int) []int {
   res := make([]int, 0, n)
   for i := n; i >= 1; i-- {
       if tr[i].Bit(sum) == 1 {
           res = append(res, arr[i])
           sum -= arr[i]
       } else {
           res = append(res, -arr[i])
       }
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var t int
   if _, err := fmt.Fscan(in, &t); err != nil {
       return
   }
   for t > 0 {
       t--
       var n, m int
       fmt.Fscan(in, &n)
       a := make([]int, n+1)
       for i := 1; i <= n; i++ {
           fmt.Fscan(in, &a[i])
       }
       fmt.Fscan(in, &m)
       b := make([]int, m+1)
       for i := 1; i <= m; i++ {
           fmt.Fscan(in, &b[i])
       }
       if n != m {
           fmt.Fprintln(out, "No")
           continue
       }
       sa, sb := 0, 0
       for i := 1; i <= n; i++ {
           sa += a[i]
           sb += b[i]
       }
       if sa%2 != 0 || sb%2 != 0 {
           fmt.Fprintln(out, "No")
           continue
       }
       sa /= 2
       sb /= 2
       // subset sum on a
       dpA, trA := calc(n, sa, a)
       if dpA[n].Bit(sa) == 0 {
           fmt.Fprintln(out, "No")
           continue
       }
       h := rev(n, sa, a, trA)
       // subset sum on b
       dpB, trB := calc(n, sb, b)
       if dpB[n].Bit(sb) == 0 {
           fmt.Fprintln(out, "No")
           continue
       }
       v := rev(n, sb, b, trB)
       nh, nv := 0, 0
       for _, x := range h {
           if x < 0 {
               nh++
           }
       }
       for _, x := range v {
           if x < 0 {
               nv++
           }
       }
       if nv < nh {
           for i := range h {
               h[i] = -h[i]
               v[i] = -v[i]
           }
       }
       // sort h: positives first by increasing abs, then negatives by increasing abs
       sort.Slice(h, func(i, j int) bool {
           a1, b1 := h[i], h[j]
           if a1*b1 < 0 {
               return a1 > 0
           }
           return abs(a1) < abs(b1)
       })
       // sort v: positives first by decreasing abs, then negatives by decreasing abs
       sort.Slice(v, func(i, j int) bool {
           a1, b1 := v[i], v[j]
           if a1*b1 < 0 {
               return a1 > 0
           }
           return abs(a1) > abs(b1)
       })
       fmt.Fprintln(out, "Yes")
       x, y := 0, 0
       for i := 0; i < n; i++ {
           y += v[i]
           fmt.Fprintln(out, x, y)
           x += h[i]
           fmt.Fprintln(out, x, y)
       }
   }
}
