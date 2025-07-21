package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

// multiply two matrices a and b of size m x m
func mul(a, b [][]int64, m int) [][]int64 {
   c := make([][]int64, m)
   for i := 0; i < m; i++ {
       c[i] = make([]int64, m)
       for k := 0; k < m; k++ {
           if a[i][k] == 0 {
               continue
           }
           aik := a[i][k]
           for j := 0; j < m; j++ {
               c[i][j] = (c[i][j] + aik*b[k][j]) % mod
           }
       }
   }
   return c
}

// matrix exponentiation a^e
func matPow(a [][]int64, e int64, m int) [][]int64 {
   // initialize result as identity
   res := make([][]int64, m)
   for i := 0; i < m; i++ {
       res[i] = make([]int64, m)
       res[i][i] = 1
   }
   for e > 0 {
       if e&1 == 1 {
           res = mul(res, a, m)
       }
       a = mul(a, a, m)
       e >>= 1
   }
   return res
}

// charToIndex maps 'a'-'z'->0-25, 'A'-'Z'->26-51
func charToIndex(c byte) int {
   if c >= 'a' && c <= 'z' {
       return int(c - 'a')
   }
   return int(c - 'A' + 26)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int64
   var m, k int
   if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
       return
   }
   // adjacency matrix
   a := make([][]int64, m)
   for i := 0; i < m; i++ {
       a[i] = make([]int64, m)
       for j := 0; j < m; j++ {
           a[i][j] = 1
       }
   }
   // apply forbidden pairs
   for i := 0; i < k; i++ {
       var s string
       fmt.Fscan(in, &s)
       if len(s) != 2 {
           continue
       }
       u := charToIndex(s[0])
       v := charToIndex(s[1])
       if u < m && v < m {
           a[u][v] = 0
       }
   }
   var result int64
   if n == 1 {
       result = int64(m) % mod
   } else {
       // total sequences = sum of all entries of a^(n-1)
       p := matPow(a, n-1, m)
       var sum int64
       for i := 0; i < m; i++ {
           for j := 0; j < m; j++ {
               sum = (sum + p[i][j]) % mod
           }
       }
       result = sum
   }
   fmt.Println(result)
}
