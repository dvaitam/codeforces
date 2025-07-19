package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

const mod = 998244353

type Mat struct {
   a [][]int
}

func mulMat(x, y Mat, s int) Mat {
   res := Mat{a: make([][]int, s)}
   for i := 0; i < s; i++ {
       res.a[i] = make([]int, s)
   }
   mod1 := mod - 1
   for i := 0; i < s; i++ {
       for k := 0; k < s; k++ {
           if x.a[i][k] == 0 {
               continue
           }
           xv := int64(x.a[i][k])
           for j := 0; j < s; j++ {
               res.a[i][j] = int((int64(res.a[i][j]) + xv*int64(y.a[k][j])) % int64(mod1))
           }
       }
   }
   return res
}

func powMat(x Mat, exp int64, s int) Mat {
   // identity
   res := Mat{a: make([][]int, s)}
   for i := 0; i < s; i++ {
       res.a[i] = make([]int, s)
       res.a[i][i] = 1
   }
   for exp > 0 {
       if exp&1 == 1 {
           res = mulMat(res, x, s)
       }
       x = mulMat(x, x, s)
       exp >>= 1
   }
   return res
}

func powMod(x, e, m int64) int64 {
   res := int64(1)
   x %= m
   for e > 0 {
       if e&1 == 1 {
           res = res * x % m
       }
       x = x * x % m
       e >>= 1
   }
   return res
}

// baby-step giant-step to solve 3^k = y mod mod, returns k or -1
func bsgs(y int) int64 {
   m := int(math.Sqrt(mod)) + 1
   tbl := make(map[int]int)
   var cur int64 = 1
   for j := 0; j < m; j++ {
       tbl[int(cur)] = j
       cur = cur * 3 % mod
   }
   factor := powMod(3, int64(mod-1-m), mod) // 3^{-m} mod mod
   // since 3^{mod-1} =1, so 3^{-m} = 3^{mod-1-m}
   cur = int64(y)
   for i := 0; i <= m; i++ {
       if j, ok := tbl[int(cur)]; ok {
           return int64(i)*int64(m) + int64(j)
       }
       cur = cur * factor % mod
   }
   return -1
}

// extended gcd: returns (g, x, y) s.t. a*x + b*y = g
func exgcd(a, b int64) (g, x, y int64) {
   if b == 0 {
       return a, 1, 0
   }
   g, x1, y1 := exgcd(b, a%b)
   x = y1
   y = x1 - (a/b)*y1
   return
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var k int
   fmt.Fscan(in, &k)
   b := make([]int, k)
   for i := 0; i < k; i++ {
       fmt.Fscan(in, &b[i])
   }
   var n int64
   var mval int
   fmt.Fscan(in, &n, &mval)
   // build matrix
   s := k
   base := Mat{a: make([][]int, s)}
   for i := 0; i < s; i++ {
       base.a[i] = make([]int, s)
   }
   for i := 0; i < s; i++ {
       base.a[0][i] = b[i] % (mod - 1)
   }
   for i := 1; i < s; i++ {
       base.a[i][i-1] = 1
   }
   // get exponent multiplier pw for f_k
   var pw int64
   if n > int64(s) {
       mtx := powMat(base, n-int64(s), s)
       pw = int64(mtx.a[0][0])
   } else {
       // n <= k: f_n = 1 for n<k, or f_k is unknown
       if n < int64(s) {
           // f_n =1, so 1 = mval? if mval!=1 no solution, else any f_k works, print 1
           if mval != 1 {
               fmt.Println(-1)
           } else {
               fmt.Println(1)
           }
           return
       }
       // n == k: f_k^1 = mval
       pw = 1
   }
   // discrete log: find exp k s.t. 3^k = mval
   kexp := bsgs(mval)
   if kexp < 0 {
       fmt.Println(-1)
       return
   }
   // solve pw * x ≡ kexp (mod mod-1)
   mod1 := int64(mod - 1)
   d := gcd64(pw, mod1)
   if kexp%d != 0 {
       fmt.Println(-1)
       return
   }
   pw1 := pw / d
   k1 := kexp / d
   mod1d := mod1 / d
   // inv = inverse of pw1 mod mod1d
   _, inv, _ := exgcd(pw1, mod1d)
   inv = (inv%mod1d + mod1d) % mod1d
   x := inv * k1 % mod1d
   // result = 3^x mod mod
   res := powMod(3, x, mod)
   fmt.Println(res)
}

func gcd64(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}
