package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m, k int
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   // domain sizes for inner rectangles
   nx := n - 2
   ny := m - 2
   // if no possible rectangle, result is 0
   if nx <= 0 || ny <= 0 {
       fmt.Println(0)
       return
   }
   // initialize A and B for t=1
   A := make([]int, nx)
   B := make([]int, ny)
   for i := 0; i < nx; i++ {
       // width w = i+1, count positions: (n - w - 1)
       A[i] = (n - (i + 1) - 1) % MOD
   }
   for i := 0; i < ny; i++ {
       B[i] = (m - (i + 1) - 1) % MOD
   }
   // iterate for moves 2..k
   for t := 2; t <= k; t++ {
       // update A
       newA := make([]int, nx)
       // prefix sums from right: P0[i] = sum_{j=i..nx-1} A[j]
       P0 := make([]int, nx+2)
       P1 := make([]int, nx+2)
       for i := nx - 1; i >= 0; i-- {
           P0[i] = (P0[i+1] + A[i]) % MOD
           P1[i] = (P1[i+1] + A[i]* (i+1) % MOD) % MOD
       }
       for i := 0; i < nx; i++ {
           idx := i + 2
           if idx <= nx-1 {
               s0 := P0[idx]
               s1 := P1[idx]
               // A_new[i] = sum A[j]*(j-i-1) = sum A[j]*(j+1) - (i+2)*sum A[j]
               val := (int64(s1) - int64(i+2)*int64(s0)%MOD + MOD) % MOD
               newA[i] = int(val)
           } else {
               newA[i] = 0
           }
       }
       A = newA
       // update B similarly
       newB := make([]int, ny)
       Q0 := make([]int, ny+2)
       Q1 := make([]int, ny+2)
       for i := ny - 1; i >= 0; i-- {
           Q0[i] = (Q0[i+1] + B[i]) % MOD
           Q1[i] = (Q1[i+1] + B[i]* (i+1) % MOD) % MOD
       }
       for i := 0; i < ny; i++ {
           idx := i + 2
           if idx <= ny-1 {
               s0 := Q0[idx]
               s1 := Q1[idx]
               val := (int64(s1) - int64(i+2)*int64(s0)%MOD + MOD) % MOD
               newB[i] = int(val)
           } else {
               newB[i] = 0
           }
       }
       B = newB
   }
   // sum up
   var sumA, sumB int64
   for i := 0; i < nx; i++ {
       sumA = (sumA + int64(A[i])) % MOD
   }
   for i := 0; i < ny; i++ {
       sumB = (sumB + int64(B[i])) % MOD
   }
   result := sumA * sumB % MOD
   fmt.Println(result)
}
