package main

import (
   "bufio"
   "fmt"
   "os"
)

// get computes sum from (x) to (x+y) in triangular numbers difference
func get(x, y int64) int64 {
   p1 := (x + y) * (x + y + 1) / 2
   p2 := (x - 1) * x / 2
   return p1 - p2
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   fmt.Fscan(reader, &t)
   for ; t > 0; t-- {
       var n int
       var x, y, s0 int64
       fmt.Fscan(reader, &n, &x, &y, &s0)
       p := x % y
       // adjust total after subtracting base values
       if s0 < p*int64(n) || (s0-p*int64(n))%y != 0 {
           writer.WriteString("NO\n")
           continue
       }
       rem := (s0 - p*int64(n)) / y
       s := int(rem)
       inf := n + 5
       dp := make([]int, s+1)
       par := make([]struct{ prev, cnt int }, s+1)
       for i := 1; i <= s; i++ {
           dp[i] = inf
       }
       // knapsack-like DP for selecting blocks of size i with cost k = i*(i-1)/2
       dp[0] = 0
       for i := 1; ; i++ {
           k := i * (i - 1) / 2
           if k > s {
               break
           }
           for j := k; j <= s; j++ {
               if dp[j-k]+i <= dp[j] {
                   dp[j] = dp[j-k] + i
                   par[j].prev = j - k
                   par[j].cnt = i
               }
           }
       }
       be := x / y
       found := false
       // try prefix lengths i
       for i := 0; i < n; i++ {
           tGet := get(be, int64(i))
           if tGet > rem {
               break
           }
           if dp[s-int(tGet)] < n-i {
               // construct sequence
               v := make([]int64, 0, n)
               for j := 0; j <= i; j++ {
                   v = append(v, be+int64(j))
               }
               k2 := s - int(tGet)
               for k2 > 0 {
                   pr := par[k2]
                   for g := 0; g < pr.cnt; g++ {
                       v = append(v, int64(g))
                   }
                   k2 = pr.prev
               }
               for len(v) < n {
                   v = append(v, 0)
               }
               writer.WriteString("YES\n")
               for idx, val := range v {
                   writer.WriteString(fmt.Sprintf("%d", val*y+p))
                   if idx+1 < len(v) {
                       writer.WriteByte(' ')
                   }
               }
               writer.WriteString("\n")
               found = true
               break
           }
       }
       if !found {
           writer.WriteString("NO\n")
       }
   }
}
