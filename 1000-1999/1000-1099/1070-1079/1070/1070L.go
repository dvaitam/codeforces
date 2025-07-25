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
   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for tc := 0; tc < t; tc++ {
       var n, m int
       fmt.Fscan(reader, &n, &m)
       adj := make([][]int, n)
       deg := make([]int, n)
       for i := 0; i < m; i++ {
           var u, v int
           fmt.Fscan(reader, &u, &v)
           u--;
           v--;
           adj[u] = append(adj[u], v)
           adj[v] = append(adj[v], u)
           deg[u]++
           deg[v]++
       }
       // if all degrees even, one state suffices
       allEven := true
       for i := 0; i < n; i++ {
           if deg[i]&1 != 0 {
               allEven = false
               break
           }
       }
       if allEven {
           fmt.Fprintln(writer, 1)
           for i := 0; i < n; i++ {
               if i > 0 {
                   writer.WriteByte(' ')
               }
               writer.WriteByte('1')
           }
           writer.WriteByte('\n')
           continue
       }
       // build linear system mod2: A * x = b
       // n equations, n variables
       cols := n + 1
       w := (cols + 63) >> 6
       A := make([][]uint64, n)
       for i := 0; i < n; i++ {
           A[i] = make([]uint64, w)
       }
       // fill adjacency and diag
       for i := 0; i < n; i++ {
           for _, j := range adj[i] {
               A[i][j>>6] |= 1 << (uint(j) & 63)
           }
           if deg[i]&1 != 0 {
               // diag
               A[i][i>>6] |= 1 << (uint(i) & 63)
               // rhs
               idx := n >> 6
               A[i][idx] |= 1 << (uint(n) & 63)
           }
       }
       pivot := make([]int, n)
       for i := range pivot {
           pivot[i] = -1
       }
       row := 0
       // gauss elimination
       for col := 0; col < n && row < n; col++ {
           sel := -1
           for i := row; i < n; i++ {
               if (A[i][col>>6]>>(uint(col)&63))&1 == 1 {
                   sel = i
                   break
               }
           }
           if sel < 0 {
               continue
           }
           A[row], A[sel] = A[sel], A[row]
           pivot[col] = row
           // eliminate
           for i := 0; i < n; i++ {
               if i != row && ((A[i][col>>6]>>(uint(col)&63))&1 == 1) {
                   // A[i] ^= A[row]
                   for k := 0; k < w; k++ {
                       A[i][k] ^= A[row][k]
                   }
               }
           }
           row++
       }
       // recover solution
       x := make([]int, n)
       for col := 0; col < n; col++ {
           if pivot[col] != -1 {
               bit := (A[pivot[col]][n>>6] >> (uint(n) & 63)) & 1
               x[col] = int(bit)
           } else {
               x[col] = 0
           }
       }
       // output r=2
       fmt.Fprintln(writer, 2)
       for i := 0; i < n; i++ {
           if i > 0 {
               writer.WriteByte(' ')
           }
           // state = x[i]+1
           writer.WriteByte(byte('0' + x[i] + 1))
       }
       writer.WriteByte('\n')
   }
}
