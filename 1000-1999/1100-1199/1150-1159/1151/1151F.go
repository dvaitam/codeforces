package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

// modPow computes a^e % mod
func modPow(a, e int64) int64 {
   res := int64(1)
   a %= mod
   for e > 0 {
       if e&1 == 1 {
           res = res * a % mod
       }
       a = a * a % mod
       e >>= 1
   }
   return res
}

// modInv computes modular inverse of a mod mod
func modInv(a int64) int64 {
   return modPow(a, mod-2)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var k int64
   if _, err := fmt.Fscan(in, &n, &k); err != nil {
       return
   }
   a := make([]int, n)
   zeros := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
       if a[i] == 0 {
           zeros++
       }
   }
   // number of zeros = zeros, left block size = zeros
   // initial misplacements = number of ones in first zeros positions
   leftSize := zeros
   onesLeft := 0
   for i := 0; i < leftSize; i++ {
       if a[i] == 1 {
           onesLeft++
       }
   }
   // maximum misplacements
   maxMis := zeros
   if n-zeros < maxMis {
       maxMis = n - zeros
   }
   d := maxMis + 1
   // total pairs = n*(n-1)/2
   totalPairs := int64(n) * int64(n-1) / 2
   invPairs := modInv(totalPairs % mod)
   // build transition matrix mat[i][j]: from state i to j
   mat := make([][]int64, d)
   for i := range mat {
       mat[i] = make([]int64, d)
   }
   for x := 0; x <= maxMis; x++ {
       // increase misplacements: swap 0 in left with 1 in right
       inc := int64(leftSize-x) * int64((n-leftSize)-x) % mod
       // decrease misplacements: swap 1 in left with 0 in right
       dec := int64(x) * int64(x) % mod
       if x+1 <= maxMis {
           mat[x][x+1] = inc * invPairs % mod
       }
       if x-1 >= 0 {
           mat[x][x-1] = dec * invPairs % mod
       }
       // stay
       sum := int64(0)
       if x+1 <= maxMis {
           sum = (sum + mat[x][x+1]) % mod
       }
       if x-1 >= 0 {
           sum = (sum + mat[x][x-1]) % mod
       }
       mat[x][x] = (mod + 1 - sum) % mod
   }
   // matrix exponentiation mat^k
   // res = identity
   res := make([][]int64, d)
   for i := range res {
       res[i] = make([]int64, d)
       res[i][i] = 1
   }
   // power
   for k > 0 {
       if k&1 == 1 {
           res = mul(res, mat, d)
       }
       mat = mul(mat, mat, d)
       k >>= 1
   }
   // answer is res[initial][0]
   ans := res[onesLeft][0]
   fmt.Println(ans)
}

// mul computes a * b for dxd matrices
func mul(a, b [][]int64, d int) [][]int64 {
   c := make([][]int64, d)
   for i := range c {
       c[i] = make([]int64, d)
   }
   for i := 0; i < d; i++ {
       for k := 0; k < d; k++ {
           if a[i][k] == 0 {
               continue
           }
           aik := a[i][k]
           for j := 0; j < d; j++ {
               c[i][j] = (c[i][j] + aik*b[k][j]) % mod
           }
       }
   }
   return c
}
