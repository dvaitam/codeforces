package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

const MOD = 1000000007

func main() {
   reader := bufio.NewReader(os.Stdin)
   line, _ := reader.ReadString('\n')
   m, _ := strconv.ParseInt(line[:len(line)-1], 10, 64)
   // digit DP to count numbers 1..m by number of lucky digits
   s := strconv.FormatInt(m, 10)
   L := len(s)
   // dp[pos][count][tight]
   var dp [11][11][2]int64
   dp[0][0][1] = 1
   for i := 0; i < L; i++ {
       digit := int(s[i] - '0')
       for cnt := 0; cnt <= i; cnt++ {
           for tight := 0; tight < 2; tight++ {
               v := dp[i][cnt][tight]
               if v == 0 {
                   continue
               }
               maxd := 9
               if tight == 1 {
                   maxd = digit
               }
               for d := 0; d <= maxd; d++ {
                   nt := 0
                   if tight == 1 && d == maxd {
                       nt = 1
                   }
                   nc := cnt
                   if d == 4 || d == 7 {
                       nc++
                   }
                   dp[i+1][nc][nt] += v
               }
           }
       }
   }
   // bucket counts
   c := make([]int64, L+1)
   for k := 0; k <= L; k++ {
       c[k] = dp[L][k][0] + dp[L][k][1]
   }
   // remove zero
   c[0]--

   // precompute factorials and inverse factorials up to 6
   fact := make([]int64, 7)
   invfact := make([]int64, 7)
   fact[0] = 1
   for i := 1; i <= 6; i++ {
       fact[i] = fact[i-1] * int64(i) % MOD
   }
   invfact[6] = modInv(fact[6])
   for i := 6; i > 0; i-- {
       invfact[i-1] = invfact[i] * int64(i) % MOD
   }
   // generate tVectors: counts for buckets
   type tv struct {
       t      []int
       sum    int
       multin int64
   }
   var tvs []tv
   t := make([]int, L+1)
   var gen func(idx, rem int)
   gen = func(idx, rem int) {
       if idx == L {
           t[idx] = rem
           // compute sum and multinomial
           sumw := 0
           for i, v := range t {
               sumw += i * v
           }
           multin := fact[6]
           for _, v := range t {
               multin = multin * invfact[v] % MOD
           }
           // store a copy of t
           tc := make([]int, len(t))
           copy(tc, t)
           tvs = append(tvs, tv{t: tc, sum: sumw, multin: multin})
           return
       }
       for x := 0; x <= rem; x++ {
           t[idx] = x
           gen(idx+1, rem-x)
       }
       t[idx] = 0
   }
   gen(0, 6)

   // compute answer
   var answer int64
   // helper for falling factorial
   for k := 0; k <= L; k++ {
       if c[k] <= 0 {
           continue
       }
       // prepare c' modulo
       cp := make([]int64, L+1)
       for i := 0; i <= L; i++ {
           cp[i] = c[i]
       }
       cp[k]--
       for i := 0; i <= L; i++ {
           cp[i] %= MOD
           if cp[i] < 0 {
               cp[i] += MOD
           }
       }
       var T int64
       for _, tvv := range tvs {
           if tvv.sum >= k {
               continue
           }
           prod := tvv.multin
           for i, cnt := range tvv.t {
               if cnt == 0 {
                   continue
               }
               // multiply P(cp[i], cnt)
               for p := 0; p < cnt; p++ {
                   prod = prod * (cp[i] - int64(p) + MOD) % MOD
               }
           }
           T += prod
           if T >= MOD {
               T -= MOD
           }
       }
       answer = (answer + (c[k]%MOD)*T) % MOD
   }
   fmt.Println(answer)
}

// modInv returns modular inverse of a under MOD
func modInv(a int64) int64 {
   return modPow(a, MOD-2)
}

// modPow computes a^e mod MOD
func modPow(a int64, e int64) int64 {
   res := int64(1)
   a %= MOD
   for e > 0 {
       if e&1 != 0 {
           res = res * a % MOD
       }
       a = a * a % MOD
       e >>= 1
   }
   return res
}
