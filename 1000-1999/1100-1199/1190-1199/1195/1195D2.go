package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 998244353

func add(a, b int) int {
   a += b
   if a >= mod {
       a -= mod
   }
   return a
}

func mul(a, b int) int {
   return int((int64(a) * int64(b)) % mod)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   // sum of digits by length and position
   // sum of digits by length and position: sumDig[L][pos]
   var sumDig [11][11]int
   var cnt [11]int
   for i := 0; i < n; i++ {
       var x int
       fmt.Fscan(reader, &x)
       // extract digits
       digits := make([]int, 0, 11)
       for x > 0 {
           digits = append(digits, x%10)
           x /= 10
       }
       L := len(digits)
       cnt[L]++
       for pos, d := range digits {
           sumDig[L][pos] = add(sumDig[L][pos], d)
       }
   }
   // precompute powers of 10 up to 20
   maxP := 20
   pow10 := make([]int, maxP+1)
   pow10[0] = 1
   for i := 1; i <= maxP; i++ {
       pow10[i] = mul(pow10[i-1], 10)
   }
   // compute answer
   var ans int
   // partX: contributions from x's digits
   for L1 := 1; L1 <= 10; L1++ {
       if cnt[L1] == 0 {
           continue
       }
       for L2 := 1; L2 <= 10; L2++ {
           if cnt[L2] == 0 {
               continue
           }
           for dx := 0; dx < L1; dx++ {
               sd := sumDig[L1][dx]
               if sd == 0 {
                   continue
               }
               // position of x's digit
               var pos int
               if dx < L2 {
                   pos = 2*dx + 1
               } else {
                   pos = 2*L2 + (dx - L2)
               }
               contrib := mul(mul(sd, cnt[L2]), pow10[pos])
               ans = add(ans, contrib)
           }
       }
   }
   // partY: contributions from y's digits
   for L2 := 1; L2 <= 10; L2++ {
       if cnt[L2] == 0 {
           continue
       }
       for L1 := 1; L1 <= 10; L1++ {
           if cnt[L1] == 0 {
               continue
           }
           for dy := 0; dy < L2; dy++ {
               sd := sumDig[L2][dy]
               if sd == 0 {
                   continue
               }
               // position of y's digit
               var pos int
               if dy < L1 {
                   pos = 2*dy
               } else {
                   pos = 2*L1 + (dy - L1)
               }
               contrib := mul(mul(sd, cnt[L1]), pow10[pos])
               ans = add(ans, contrib)
           }
       }
   }
   fmt.Fprintln(writer, ans)
}
