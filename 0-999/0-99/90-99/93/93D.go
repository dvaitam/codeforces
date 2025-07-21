package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var L, R int64
   fmt.Fscan(reader, &L, &R)
   const MOD = 1000000007
   inv2 := int64(500000004)
   // Colors: 0=W,1=Y,2=R,3=B
   allowedAdj := func(i, j int) bool {
       if i == j {
           return false
       }
       if (i == 0 && j == 1) || (i == 1 && j == 0) {
           return false
       }
       if (i == 2 && j == 3) || (i == 3 && j == 2) {
           return false
       }
       return true
   }
   forbiddenTriple := map[[3]int]bool{
       {3, 0, 2}: true, // B,W,R
       {2, 0, 3}: true, // R,W,B
   }
   // Enumerate allowed adjacency pairs
   var pair [][2]int
   idx := make(map[[2]int]int)
   for i := 0; i < 4; i++ {
       for j := 0; j < 4; j++ {
           if allowedAdj(i, j) {
               id := len(pair)
               pair = append(pair, [2]int{i, j})
               idx[[2]int{i, j}] = id
           }
       }
   }
   m := len(pair) // 8
   // Build transition matrix T (m x m)
   T := make([][]int64, m)
   for i := range T {
       T[i] = make([]int64, m)
   }
   for a := 0; a < m; a++ {
       i1, i2 := pair[a][0], pair[a][1]
       for k := 0; k < 4; k++ {
           if !allowedAdj(i2, k) {
               continue
           }
           if forbiddenTriple[[3]int{i1, i2, k}] {
               continue
           }
           b := idx[[2]int{i2, k}]
           T[a][b] = 1
       }
   }
   // Build augmented matrix M of size (m+1)x(m+1)
   size := m + 1
   M := make([][]int64, size)
   for i := 0; i < size; i++ {
       M[i] = make([]int64, size)
   }
   // Top-left block = T
   for i := 0; i < m; i++ {
       for j := 0; j < m; j++ {
           M[i][j] = T[i][j]
       }
   }
   // Last row: sums
   for j := 0; j < m; j++ {
       M[m][j] = 1
   }
   M[m][m] = 1
   // Initial vector w2 = [v2 (n=2), S2]
   v2 := make([]int64, size)
   for i := 0; i < m; i++ {
       v2[i] = 1
   }
   v2[m] = int64(m)
   // Matrix multiplication (square)
   mul := func(A, B [][]int64) [][]int64 {
       n := len(A)
       C := make([][]int64, n)
       for i := range C {
           C[i] = make([]int64, n)
       }
       for i := 0; i < n; i++ {
           for k := 0; k < n; k++ {
               if A[i][k] == 0 {
                   continue
               }
               aik := A[i][k]
               for j := 0; j < n; j++ {
                   C[i][j] = (C[i][j] + aik*B[k][j]) % MOD
               }
           }
       }
       return C
   }
   // Matrix exponentiation
   pow := func(mat [][]int64, e int64) [][]int64 {
       n := len(mat)
       res := make([][]int64, n)
       for i := 0; i < n; i++ {
           res[i] = make([]int64, n)
           res[i][i] = 1
       }
       base := mat
       for e > 0 {
           if e&1 == 1 {
               res = mul(res, base)
           }
           base = mul(base, base)
           e >>= 1
       }
       return res
   }
   // sumA(b): sum of A_n for n=2..b
   sumA := func(b int64) int64 {
       if b < 2 {
           return 0
       }
       Mexp := pow(M, b-2)
       var s int64
       for j := 0; j < size; j++ {
           s = (s + v2[j]*Mexp[j][m]) % MOD
       }
       return s
   }
   // Compute S1 = sum_{n=L..R} A_n
   var S1 int64
   if L <= 1 && R >= 1 {
       S1 = (S1 + 4) % MOD
   }
   l2 := L
   if l2 < 2 {
       l2 = 2
   }
   if R >= 2 && l2 <= R {
       S1 = (S1 + (sumA(R)-sumA(l2-1)+MOD)%MOD) % MOD
   }
   // Compute S2 = sum_{odd n in [L..R]} A_{(n+1)/2}
   k1 := (L + 2) / 2
   if (L%2)==0 {
       k1 = (L/2)+1
   }
   k2 := (R + 1) / 2
   var S2 int64
   if k1 <= k2 {
       if k1 <= 1 && k2 >= 1 {
           S2 = (S2 + 4) % MOD
       }
       kk1 := k1
       if kk1 < 2 {
           kk1 = 2
       }
       if k2 >= 2 && kk1 <= k2 {
           S2 = (S2 + (sumA(k2)-sumA(kk1-1)+MOD)%MOD) % MOD
       }
   }
   ans := (S1 + S2) % MOD * inv2 % MOD
   fmt.Println(ans)
}
