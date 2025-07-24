package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, m int
   fmt.Fscan(reader, &n, &m)
   edges := make([][2]int, m)
   for i := range edges {
       fmt.Fscan(reader, &edges[i][0], &edges[i][1])
       edges[i][0]--
       edges[i][1]--
   }
   blocks := (2*n + 63) >> 6
   mat := make([][]uint64, n)
   for i := 0; i < n; i++ {
       mat[i] = make([]uint64, blocks)
       // identity on right half
       idx := n + i
       mat[i][idx>>6] |= 1 << (idx & 63)
   }
   // set adjacency on left half
   for _, e := range edges {
       a := e[0]
       b := e[1]
       mat[a][b>>6] |= 1 << (b & 63)
   }
   // forward elimination
   for i := 0; i < n; i++ {
       pivot := i
       for pivot < n && ((mat[pivot][i>>6]>>(i&63))&1) == 0 {
           pivot++
       }
       if pivot == n {
           // should not happen, matrix is invertible
           fmt.Fprintln(writer, "NO_SOL")
           return
       }
       if pivot != i {
           mat[i], mat[pivot] = mat[pivot], mat[i]
       }
       rowi := mat[i]
       for j := i + 1; j < n; j++ {
           if ((mat[j][i>>6] >> (i & 63)) & 1) == 1 {
               rowj := mat[j]
               for k := 0; k < blocks; k++ {
                   rowj[k] ^= rowi[k]
               }
           }
       }
   }
   // backward elimination
   for i := n - 1; i >= 0; i-- {
       rowi := mat[i]
       for j := 0; j < i; j++ {
           if ((mat[j][i>>6] >> (i & 63)) & 1) == 1 {
               rowj := mat[j]
               for k := 0; k < blocks; k++ {
                   rowj[k] ^= rowi[k]
               }
           }
       }
   }
   // answer queries: for edge (a,b), minor at (a,b) is inv[b][a] = mat[b][n+a]
   for _, e := range edges {
       a := e[0]
       b := e[1]
       idx := n + a
       bit := (mat[b][idx>>6] >> (idx & 63)) & 1
       if bit == 0 {
           fmt.Fprintln(writer, "YES")
       } else {
           fmt.Fprintln(writer, "NO")
       }
   }
}
