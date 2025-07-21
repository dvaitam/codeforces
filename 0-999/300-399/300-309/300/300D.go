package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

const P = 7340033

func modAdd(a, b int) int {
   a += b
   if a >= P {
       a -= P
   }
   return a
}
func modSub(a, b int) int {
   a -= b
   if a < 0 {
       a += P
   }
   return a
}
func modMul(a, b int) int {
   return (a * b) % P
}
func modPow(a, e int) int {
   res := 1
   for e > 0 {
       if e&1 != 0 {
           res = modMul(res, a)
       }
       a = modMul(a, a)
       e >>= 1
   }
   return res
}
func modInv(a int) int {
   // P is prime
   return modPow((a%P+P)%P, P-2)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   // read q
   line, _ := in.ReadString('\n')
   q, _ := strconv.Atoi(line[:len(line)-1])
   ns := make([]int, q)
   ks := make([]int, q)
   maxK := 0
   for i := 0; i < q; i++ {
       line, _ = in.ReadString('\n')
       // parse n and k
       f := bufio.NewScanner(bufio.NewReaderString(line))
       f.Split(bufio.ScanWords)
       f.Scan()
       n64, _ := strconv.ParseInt(f.Text(), 10, 64)
       f.Scan()
       k, _ := strconv.Atoi(f.Text())
       ns[i] = int((n64 - 1) / 2)
       ks[i] = k
       if k > maxK {
           maxK = k
       }
   }
   // precompute dp up to X = 2*maxK
   X := 2*maxK + 5
   // dp[i][j]
   dp := make([][]int, X+1)
   for i := range dp {
       dp[i] = make([]int, maxK+1)
   }
   dp[0][0] = 1
   for i := 1; i <= X; i++ {
       dp[i][0] = 1
       // s = 2*i-1
       t := modMul(2*i-1, 2*i-1)
       for j := 1; j <= maxK; j++ {
           // dp[i][j] = dp[i-1][j] + t * dp[i-1][j-1]
           dp[i][j] = modAdd(dp[i-1][j], modMul(t, dp[i-1][j-1]))
       }
   }
   // group queries by k
   byK := make(map[int][]int)
   for i, k := range ks {
       byK[k] = append(byK[k], i)
   }
   ans := make([]int, q)
   // process k=0
   if idxs, ok := byK[0]; ok {
       for _, idx := range idxs {
           ans[idx] = 1
       }
   }
   // temp arrays
   for k, idxs := range byK {
       if k == 0 {
           continue
       }
       D := 2 * k
       if D > X {
           D = X
       }
       // prepare y = dp[0..D][k]
       y := dp[:D+1]
       // precompute invDenom
       invDenom := make([]int, D+1)
       // fac for denom
       // denom[i] = fac[i]*fac[D-i] * (-1)^(D-i)
       // precompute small factorials
       fac := make([]int, D+1)
       fac[0] = 1
       for i := 1; i <= D; i++ {
           fac[i] = modMul(fac[i-1], i)
       }
       for i := 0; i <= D; i++ {
           d := modMul(fac[i], fac[D-i])
           if (D-i)&1 == 1 {
               d = modSub(0, d)
           }
           invDenom[i] = modInv(d)
       }
       // process queries for this k
       for _, idx := range idxs {
           m := ns[idx]
           if m < k {
               ans[idx] = 0
           } else if m <= D {
               ans[idx] = dp[m][k]
           } else {
               // compute t = prod_{j=0..D} (m - j)
               tprod := 1
               for j := 0; j <= D; j++ {
                   tprod = modMul(tprod, modSub(m, j))
               }
               // compute inv of (m-j) via prefix-suffix
               a := make([]int, D+1)
               for j := 0; j <= D; j++ {
                   a[j] = modSub(m, j)
               }
               pre := make([]int, D+1)
               pre[0] = a[0]
               for j := 1; j <= D; j++ {
                   pre[j] = modMul(pre[j-1], a[j])
               }
               suf := make([]int, D+2)
               suf[D+1] = 1
               for j := D; j >= 0; j-- {
                   suf[j] = modMul(suf[j+1], a[j])
               }
               // sum
               res := 0
               for i := 0; i <= D; i++ {
                   // numerator without (m-i)
                   var num int
                   if i == 0 {
                       num = suf[1]
                   } else if i == D {
                       num = pre[D-1]
                   } else {
                       num = modMul(pre[i-1], suf[i+1])
                   }
                   term := modMul(y[i][k], modMul(num, invDenom[i]))
                   res = modAdd(res, term)
               }
               // multiply by tprod cancels one extra (m-i) in denom? Actually formula uses num without denom includes(-1), correct
               // but here tprod = (m-0)...(m-D) = num*(m-i) so num = tprod * inv(m-i)
               // we used num = prod_{j!=i}(m-j), so do not multiply tprod
               ans[idx] = res
           }
       }
   }
   // output answers
   for i := 0; i < q; i++ {
       fmt.Fprintln(out, ans[i])
   }
}
