package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007
const N = 1000005

var fact = make([]int, N)
var rfact = make([]int, N)
var pwM = make([]int, N)
var pwN = make([]int, N)

// mul returns (a * b) % mod
func mul(a, b int) int {
   return int(int64(a) * int64(b) % mod)
}

// add returns (a + b) % mod
func add(a, b int) int {
   s := a + b
   if s >= mod {
       s -= mod
   }
   return s
}

// binpow computes x^pw % mod
func binpow(x, pw int) int {
   res := 1
   tmp := x
   for pw > 0 {
       if pw&1 == 1 {
           res = mul(res, tmp)
       }
       tmp = mul(tmp, tmp)
       pw >>= 1
   }
   return res
}

// CNK computes combination n choose k
func CNK(n, k int) int {
   if n < k {
       return 0
   }
   res := mul(rfact[k], rfact[n-k])
   return mul(res, fact[n])
}

// downFact computes falling factorial: from!/(from-cnt)! if from>=cnt, else from!.
func downFact(from, cnt int) int {
   if cnt <= 0 {
       return 1
   }
   if from >= cnt {
       return mul(fact[from], rfact[from-cnt])
   }
   return fact[from]
}

// build precomputes factorials, inverse factorials, and powers
func build(n, m int) {
   fact[0] = 1
   for i := 1; i < N; i++ {
       fact[i] = mul(fact[i-1], i)
   }
   rfact[N-1] = binpow(fact[N-1], mod-2)
   for i := N - 2; i >= 0; i-- {
       rfact[i] = mul(rfact[i+1], i+1)
   }
   pwM[0], pwN[0] = 1, 1
   for i := 1; i < N; i++ {
       pwM[i] = mul(pwM[i-1], m)
       pwN[i] = mul(pwN[i-1], n)
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m, a, b int
   fmt.Fscan(reader, &n, &m, &a, &b)

   build(n, m)
   res := 0
   for i := 1; i < n && i <= m; i++ {
       fi := CNK(m-1, m-i)
       k1 := downFact(n-2, i-1)
       // k2 is always 1
       var k3 int
       if i+1 < n {
           k3 = mul(i+1, pwN[n-i-2])
       } else {
           k3 = 1
       }
       k4 := pwM[n-1-i]
       tmp := mul(fi, k1)
       tmp = mul(tmp, k3)
       tmp = mul(tmp, k4)
       res = add(res, tmp)
   }

   fmt.Fprintln(writer, res)
}
