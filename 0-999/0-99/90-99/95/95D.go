package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

func main() {
   reader := bufio.NewReader(os.Stdin)
   var t, k int
   fmt.Fscan(reader, &t, &k)
   L := make([]string, t)
   R := make([]string, t)
   maxlen := 0
   for i := 0; i < t; i++ {
       fmt.Fscan(reader, &L[i], &R[i])
       if len(R[i]) > maxlen {
           maxlen = len(R[i])
       }
       if len(L[i]) > maxlen {
           maxlen = len(L[i])
       }
   }
   // Precompute dp0[len][s][last]: number of sequences of length len from state (s,last) ending with state 2
   lastDim := k + 2
   dp0 := make([][][]int, maxlen+1)
   for i := 0; i <= maxlen; i++ {
       dp0[i] = make([][]int, 3)
       for s := 0; s < 3; s++ {
           dp0[i][s] = make([]int, lastDim)
       }
   }
   // Base: len=0
   for s := 0; s < 3; s++ {
       for last := 1; last < lastDim; last++ {
           if s == 2 {
               dp0[0][s][last] = 1
           }
       }
   }
   // Transitions
   transNon := func(s, last int) (int, int) {
       // non-lucky digit
       nl := last + 1
       if nl >= lastDim {
           nl = lastDim - 1
       }
       return s, nl
   }
   transLucky := func(s, last int) (int, int) {
       var ns int
       if s == 0 {
           ns = 1
       } else if s == 1 {
           if last <= k {
               ns = 2
           } else {
               ns = 1
           }
       } else {
           ns = 2
       }
       return ns, 1
   }
   // Build dp0
   for length := 1; length <= maxlen; length++ {
       for s := 0; s < 3; s++ {
           for last := 1; last < lastDim; last++ {
               var sum int64
               ns, nl := transNon(s, last)
               sum += int64(8) * int64(dp0[length-1][ns][nl])
               ls, ll := transLucky(s, last)
               sum += int64(2) * int64(dp0[length-1][ls][ll])
               dp0[length][s][last] = int(sum % MOD)
           }
       }
   }
   // Precompute full-length counts
   countFull := make([]int, maxlen+1)
   for l := 1; l <= maxlen; l++ {
       if l-1 >= 0 {
           a := dp0[l-1][0][k+1]
           b := dp0[l-1][1][1]
           countFull[l] = int((int64(7)*int64(a) + int64(2)*int64(b)) % MOD)
       }
   }
   // function to decrement string number by 1
   decStr := func(s string) string {
       b := []byte(s)
       i := len(b) - 1
       for i >= 0 && b[i] == '0' {
           b[i] = '9'
           i--
       }
       if i >= 0 {
           b[i]--
       }
       // trim leading zeros
       j := 0
       for j < len(b) && b[j] == '0' {
           j++
       }
       if j == len(b) {
           return "0"
       }
       return string(b[j:])
   }
   // compute f(X)
   f := func(X string) int {
       n := len(X)
       var ans int64
       // full lengths < n
       for l := 1; l < n; l++ {
           ans += int64(countFull[l])
       }
       ans %= MOD
       s0, last0 := 0, k+1
       for i := 0; i < n; i++ {
           d := int(X[i] - '0')
           low := 0
           if i == 0 {
               low = 1
           }
           if d > low {
               total := d - low
               luckyLt := 0
               if 4 >= low && 4 < d {
                   luckyLt++
               }
               if 7 >= low && 7 < d {
                   luckyLt++
               }
               nonLuckyLt := total - luckyLt
               // non-lucky group
               ns, nl := transNon(s0, last0)
               ans += int64(nonLuckyLt) * int64(dp0[n-i-1][ns][nl])
               // lucky group
               ls, ll := transLucky(s0, last0)
               ans += int64(luckyLt) * int64(dp0[n-i-1][ls][ll])
               ans %= MOD
           }
           // process actual digit
           if d == 4 || d == 7 {
               s0, last0 = transLucky(s0, last0)
           } else {
               s0, last0 = transNon(s0, last0)
           }
           if i == n-1 && s0 == 2 {
               ans = (ans + 1) % MOD
           }
       }
       return int(ans)
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i := 0; i < t; i++ {
       ldec := decStr(L[i])
       ans := f(R[i]) - f(ldec)
       if ans < 0 {
           ans += MOD
       }
       fmt.Fprintln(writer, ans)
   }
}
