package main

import (
   "bufio"
   "fmt"
   "io"
   "os"
)

func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   data, _ := io.ReadAll(os.Stdin)
   n := 0
   pos := 0
   // read integer from data
   nextInt := func() int {
       // skip non-number
       for pos < len(data) && (data[pos] < '0' || data[pos] > '9') && data[pos] != '-' {
           pos++
       }
       sign := 1
       if pos < len(data) && data[pos] == '-' {
           sign = -1
           pos++
       }
       val := 0
       for pos < len(data) && data[pos] >= '0' && data[pos] <= '9' {
           val = val*10 + int(data[pos]-'0')
           pos++
       }
       return val * sign
   }
   // parse n
   n = nextInt()
   a := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       a[i] = int64(nextInt())
   }
   p := make([]int, n+1)
   for i := 2; i <= n; i++ {
       p[i] = nextInt()
   }
   s := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       s[i] = a[i]
   }
   for i := n; i >= 2; i-- {
       s[p[i]] += s[i]
   }
   cnt := make([]int, n+1)
   tot := s[1]
   for i := 1; i <= n; i++ {
       g := gcd(tot, s[i])
       tmp := tot / g
       if tmp <= int64(n) {
           cnt[tmp]++
       }
   }
   // accumulate counts for multiples
   for i := n; i >= 1; i-- {
       for j := i + i; j <= n; j += i {
           cnt[j] += cnt[i]
       }
   }
   const mod = 1000000007
   dp := make([]int, n+1)
   dp[1] = 1
   ans := 0
   for i := 1; i <= n; i++ {
       if cnt[i] == i {
           di := dp[i]
           for j := i + i; j <= n; j += i {
               dp[j] += di
               if dp[j] >= mod {
                   dp[j] -= mod
               }
           }
           ans += di
           if ans >= mod {
               ans -= mod
           }
       }
   }
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   fmt.Fprintln(w, ans)
}
