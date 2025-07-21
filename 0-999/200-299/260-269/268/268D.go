package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000009

func modPow(a, e int) int {
   res := 1
   x := a % MOD
   for e > 0 {
       if e&1 == 1 {
           res = int((int64(res) * x) % MOD)
       }
       x = int((int64(x) * x) % MOD)
       e >>= 1
   }
   return res
}

// compute A(k): number of strings of length M over alphabet {0..3}
// such that for k specified directions, no run of zeros (i.e., non-occurrence) of length >= h
// Here letters 0..k-1 are constrained, others free
func computeA(M, h, k int) int {
   // special for k==0: no constraints
   if k == 0 {
       return modPow(4, M)
   }
   // special for k==4: if M < h, vacuously true, else impossible
   if k == 4 {
       if M < h {
           return modPow(4, M)
       }
       return 0
   }
   // state dimension = h^k
   dim := 1
   for i := 0; i < k; i++ {
       dim *= h
   }
   // precompute transitions
   // for each state index s in [0,dim), we compute trans for k constrained letters, and for other letters
   trans := make([][]int, k+1) // 0..k-1 for each constrained letter, k for others
   for j := 0; j <= k; j++ {
       trans[j] = make([]int, dim)
   }
   // decode and encode functions
   decode := func(s int) []int {
       z := make([]int, k)
       for i := 0; i < k; i++ {
           z[i] = s % h
           s /= h
       }
       return z
   }
   encode := func(z []int) int {
       s := 0
       mul := 1
       for i := 0; i < k; i++ {
           s += z[i] * mul
           mul *= h
       }
       return s
   }
   for s := 0; s < dim; s++ {
       z := decode(s)
       // transitions for each constrained letter j
       for j := 0; j < k; j++ {
           ok := true
           nz := make([]int, k)
           for i := 0; i < k; i++ {
               if i == j {
                   nz[i] = 0
               } else {
                   if z[i]+1 >= h {
                       ok = false
                       break
                   }
                   nz[i] = z[i] + 1
               }
           }
           if ok {
               trans[j][s] = encode(nz)
           } else {
               trans[j][s] = -1
           }
       }
       // transition for other letters (4-k choices)
       ok := true
       nz := make([]int, k)
       for i := 0; i < k; i++ {
           if z[i]+1 >= h {
               ok = false
               break
           }
           nz[i] = z[i] + 1
       }
       if ok {
           trans[k][s] = encode(nz)
       } else {
           trans[k][s] = -1
       }
   }
   // DP arrays
   dp := make([]int, dim)
   ndp := make([]int, dim)
   // initial state: zeros since last occurrence = 0
   dp[0] = 1
   others := 4 - k
   for pos := 0; pos < M; pos++ {
       // clear ndp
       for i := range ndp {
           ndp[i] = 0
       }
       for s := 0; s < dim; s++ {
           v := dp[s]
           if v == 0 {
               continue
           }
           // constrained letters
           for j := 0; j < k; j++ {
               t := trans[j][s]
               if t >= 0 {
                   ndp[t] = (ndp[t] + v) % MOD
               }
           }
           // other letters
           t := trans[k][s]
           if t >= 0 {
               ndp[t] = (ndp[t] + int((int64(v) * int64(others)) % MOD)) % MOD
           }
       }
       dp, ndp = ndp, dp
   }
   // sum dp
   sum := 0
   for _, v := range dp {
       sum = (sum + v) % MOD
   }
   return sum
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, h int
   fmt.Fscan(in, &n, &h)
   // M = n-h+1
   M := n - h + 1
   // precompute A(k)
   A := make([]int, 5)
   for k := 0; k <= 4; k++ {
       A[k] = computeA(M, h, k)
   }
   // inclusion-exclusion f(M) = sum_{k=0..4} (-1)^k * C(4,k) * A[k]
   C4 := []int{1, 4, 6, 4, 1}
   f := 0
   for k := 0; k <= 4; k++ {
       term := int((int64(C4[k]) * int64(A[k])) % MOD)
       if k%2 == 1 {
           f = (f - term + MOD) % MOD
       } else {
           f = (f + term) % MOD
       }
   }
   // g(n) = f * 4^(n-M)
   powTail := modPow(4, n-M)
   g := int((int64(f) * int64(powTail)) % MOD)
   total := modPow(4, n)
   ans := (total - g + MOD) % MOD
   fmt.Println(ans)
}
