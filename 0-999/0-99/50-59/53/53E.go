package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
)

// Bareiss algorithm for exact determinant of integer matrix
func det(mat [][]int64) int64 {
   n := len(mat)
   if n == 0 {
       return 1
   }
   // copy matrix
   M := make([][]int64, n)
   for i := range mat {
       M[i] = make([]int64, n)
       copy(M[i], mat[i])
   }
   for k := 0; k < n; k++ {
       if M[k][k] == 0 {
           return 0
       }
       for i := k + 1; i < n; i++ {
           for j := k + 1; j < n; j++ {
               num := M[i][j]*M[k][k] - M[i][k]*M[k][j]
               den := int64(1)
               if k > 0 {
                   den = M[k-1][k-1]
               }
               M[i][j] = num / den
           }
           M[i][k] = 0
       }
   }
   return M[n-1][n-1]
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m, k int
   fmt.Fscan(reader, &n, &m, &k)
   adj := make([][]bool, n)
   for i := 0; i < n; i++ {
       adj[i] = make([]bool, n)
   }
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       u--; v--
       adj[u][v] = true
       adj[v][u] = true
   }
   total := int64(0)
   // iterate subsets L of size k: leaves
   for mask := 0; mask < (1 << n); mask++ {
       if bits.OnesCount(uint(mask)) != k {
           continue
       }
       // compute R = V \ L
       rmask := ((1 << n) - 1) ^ mask
       // degree product for leaves
       prod := int64(1)
       valid := true
       for v := 0; v < n; v++ {
           if mask&(1<<v) != 0 {
               deg := 0
               for u := 0; u < n; u++ {
                   if rmask&(1<<u) != 0 && adj[v][u] {
                       deg++
                   }
               }
               if deg == 0 {
                   valid = false
                   break
               }
               prod *= int64(deg)
           }
       }
       if !valid {
           continue
       }
       // collect R vertices
       var rverts []int
       for v := 0; v < n; v++ {
           if rmask&(1<<v) != 0 {
               rverts = append(rverts, v)
           }
       }
       rlen := len(rverts)
       var tcount int64 = 1
       if rlen > 1 {
           // build Laplacian matrix of size rlen
           L := make([][]int64, rlen)
           for i := range L {
               L[i] = make([]int64, rlen)
           }
           // fill
           for i := 0; i < rlen; i++ {
               for j := i + 1; j < rlen; j++ {
                   u := rverts[i]
                   v := rverts[j]
                   if adj[u][v] {
                       L[i][i]++
                       L[j][j]++
                       L[i][j]--
                       L[j][i]--
                   }
               }
           }
           // build minor by removing last row/col
           sz := rlen - 1
           M := make([][]int64, sz)
           for i := 0; i < sz; i++ {
               M[i] = make([]int64, sz)
               for j := 0; j < sz; j++ {
                   M[i][j] = L[i][j]
               }
           }
           tcount = det(M)
           if tcount == 0 {
               continue
           }
       }
       total += prod * tcount
   }
   fmt.Println(total)
}
