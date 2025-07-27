package main

import (
   "bufio"
   "fmt"
   "math/big"
   "os"
)

const mod = 998244353

func gcd(a, b int) int {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   N := new(big.Int)
   N, _ = N.SetString(s, 10)
   ans := 0
   // iterate reduced fractions a/b with gcd(a,b)=1, 1<=a,b<=9
   for a := 1; a <= 9; a++ {
       for b := 1; b <= 9; b++ {
           if gcd(a, b) != 1 {
               continue
           }
           m := a
           if b > a {
               m = b
           }
           // K = floor(N / m)
           K := new(big.Int).Div(N, big.NewInt(int64(m)))
           if K.Sign() == 0 {
               continue
           }
           ks := K.String()
           L := len(ks)
           // dp[tight][carry1][carry2][seen1][seen2]
           var dp [2][10][10][2][2]int
           var ndp [2][10][10][2][2]int
           dp[1][0][0][0][0] = 1
           for i := 0; i < L; i++ {
               // clear ndp
               for t := 0; t < 2; t++ {
                   for c1 := 0; c1 < 10; c1++ {
                       for c2 := 0; c2 < 10; c2++ {
                           for s1 := 0; s1 < 2; s1++ {
                               for s2 := 0; s2 < 2; s2++ {
                                   ndp[t][c1][c2][s1][s2] = 0
                               }
                           }
                       }
                   }
               }
               digit := int(ks[i] - '0')
               for t := 0; t < 2; t++ {
                   for c1 := 0; c1 < 10; c1++ {
                       for c2 := 0; c2 < 10; c2++ {
                           for s1 := 0; s1 < 2; s1++ {
                               for s2 := 0; s2 < 2; s2++ {
                                   v := dp[t][c1][c2][s1][s2]
                                   if v == 0 {
                                       continue
                                   }
                                   maxd := 9
                                   if t == 1 {
                                       maxd = digit
                                   }
                                   for d := 0; d <= maxd; d++ {
                                       nt := 0
                                       if t == 1 && d == digit {
                                           nt = 1
                                       }
                                       x := c1 + d*a
                                       d1 := x % 10
                                       nc1 := x / 10
                                       y := c2 + d*b
                                       d2 := y % 10
                                       nc2 := y / 10
                                       ns1 := s1
                                       if d1 == a {
                                           ns1 = 1
                                       }
                                       ns2 := s2
                                       if d2 == b {
                                           ns2 = 1
                                       }
                                       ndp[nt][nc1][nc2][ns1][ns2] += v
                                       if ndp[nt][nc1][nc2][ns1][ns2] >= mod {
                                           ndp[nt][nc1][nc2][ns1][ns2] -= mod
                                       }
                                   }
                               }
                           }
                       }
                   }
               }
               // swap dp and ndp
               dp = ndp
           }
           // sum valid
           total := 0
           for t := 0; t < 2; t++ {
               for c1 := 0; c1 < 10; c1++ {
                   for c2 := 0; c2 < 10; c2++ {
                       for s1 := 0; s1 < 2; s1++ {
                           for s2 := 0; s2 < 2; s2++ {
                               v := dp[t][c1][c2][s1][s2]
                               if v == 0 {
                                   continue
                               }
                               fs1 := (s1 == 1) || (c1 == a)
                               fs2 := (s2 == 1) || (c2 == b)
                               if fs1 && fs2 {
                                   total += v
                                   if total >= mod {
                                       total -= mod
                                   }
                               }
                           }
                       }
                   }
               }
           }
           ans += total
           if ans >= mod {
               ans -= mod
           }
       }
   }
   fmt.Println(ans)
}
