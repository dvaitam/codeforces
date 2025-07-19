package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   N   int
   A   []int
   g   [][]bool
   vl  [][]bool
   vr  [][]bool
   lt  [][]bool
   rt  [][]bool
)

func gcd(a, b int) int {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func lrecur(i, j int) bool {
   if i > j {
       return true
   }
   if i == j {
       return g[i-1][i]
   }
   if !vl[i][j] {
       vl[i][j] = true
       for k := i; k <= j; k++ {
           if g[i-1][k] && rrecur(i, k-1) && lrecur(k+1, j) {
               lt[i][j] = true
               break
           }
       }
   }
   return lt[i][j]
}

func rrecur(i, j int) bool {
   if i > j {
       return true
   }
   if i == j {
       return g[j][j+1]
   }
   if !vr[i][j] {
       vr[i][j] = true
       for k := i; k <= j; k++ {
           if g[k][j+1] && rrecur(i, k-1) && lrecur(k+1, j) {
               rt[i][j] = true
               break
           }
       }
   }
   return rt[i][j]
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   fmt.Fscan(reader, &N)
   A = make([]int, N)
   for i := 0; i < N; i++ {
       fmt.Fscan(reader, &A[i])
   }

   g = make([][]bool, N)
   for i := range g {
       g[i] = make([]bool, N)
   }
   for i := 0; i < N; i++ {
       for j := i; j < N; j++ {
           good := gcd(A[i], A[j]) > 1
           g[i][j] = good
           g[j][i] = good
       }
   }

   vl = make([][]bool, N)
   vr = make([][]bool, N)
   lt = make([][]bool, N)
   rt = make([][]bool, N)
   for i := 0; i < N; i++ {
       vl[i] = make([]bool, N)
       vr[i] = make([]bool, N)
       lt[i] = make([]bool, N)
       rt[i] = make([]bool, N)
   }

   ok := false
   if lrecur(1, N-1) || rrecur(0, N-2) {
       ok = true
   } else {
       for k := 1; k < N-1; k++ {
           if rrecur(0, k-1) && lrecur(k+1, N-1) {
               ok = true
               break
           }
       }
   }

   if ok {
       fmt.Fprint(writer, "Yes")
   } else {
       fmt.Fprint(writer, "No")
   }
}
