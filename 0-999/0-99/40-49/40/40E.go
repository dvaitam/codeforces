package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, k int
   fmt.Fscan(in, &n, &m)
   fmt.Fscan(in, &k)
   // Prepare known cells and parity adjustments
   known := make([][]bool, n)
   for i := 0; i < n; i++ {
       known[i] = make([]bool, m)
   }
   bRow := make([]int, n)
   bCol := make([]int, m)
   // target parity is 1 for each row/col
   for i := 0; i < n; i++ {
       bRow[i] = 1
   }
   for j := 0; j < m; j++ {
       bCol[j] = 1
   }
   // count known per row/col
   knownR := make([]int, n)
   knownC := make([]int, m)
   for t := 0; t < k; t++ {
       var a, b, c int
       fmt.Fscan(in, &a, &b, &c)
       a--
       b--
       // y=0 for +1, 1 for -1
       y := 0
       if c < 0 {
           y = 1
       }
       known[a][b] = true
       knownR[a]++
       knownC[b]++
       // adjust parity
       bRow[a] ^= y
       bCol[b] ^= y
   }
   var p int64
   fmt.Fscan(in, &p)
   // Check isolated rows/columns
   for i := 0; i < n; i++ {
       if knownR[i] == m && bRow[i] != 0 {
           fmt.Println(0)
           return
       }
   }
   for j := 0; j < m; j++ {
       if knownC[j] == n && bCol[j] != 0 {
           fmt.Println(0)
           return
       }
   }
   // Global parity check
   sumR, sumC := 0, 0
   for i := 0; i < n; i++ {
       sumR ^= (bRow[i] & 1)
   }
   for j := 0; j < m; j++ {
       sumC ^= (bCol[j] & 1)
   }
   if sumR != sumC {
       fmt.Println(0)
       return
   }
   // Count connected components in unknown-edge graph
   visitedR := make([]bool, n)
   visitedC := make([]bool, m)
   comps := 0
   // BFS using queue of nodes (0..n-1 rows, n..n+m-1 cols)
   var queue []int
   for i := 0; i < n; i++ {
       if !visitedR[i] && knownR[i] < m {
           comps++
           // start BFS from row i
           visitedR[i] = true
           queue = queue[:0]
           queue = append(queue, i)
           for qi := 0; qi < len(queue); qi++ {
               u := queue[qi]
               if u < n {
                   r := u
                   // neighbors: columns with unknown
                   for cj := 0; cj < m; cj++ {
                       if known[r][cj] {
                           continue
                       }
                       if visitedC[cj] {
                           continue
                       }
                       visitedC[cj] = true
                       queue = append(queue, n+cj)
                   }
               } else {
                   cj := u - n
                   // neighbors: rows with unknown
                   for rr := 0; rr < n; rr++ {
                       if known[rr][cj] {
                           continue
                       }
                       if visitedR[rr] {
                           continue
                       }
                       visitedR[rr] = true
                       queue = append(queue, rr)
                   }
               }
           }
       }
   }
   // count isolated columns (no unknown edges)
   for j := 0; j < m; j++ {
       if knownC[j] == n {
           comps++
       }
   }
   // number of unknown variables
   U := int64(n)*int64(m) - int64(k)
   // free variables = U - ((n+m) - comps)
   F := U - int64(n+m) + int64(comps)
   // compute 2^F mod p
   res := modPow(2, F, p)
   fmt.Println(res)
}

// fast exponentiation
func modPow(a int64, e int64, mod int64) int64 {
   res := int64(1)
   a %= mod
   for e > 0 {
       if e&1 != 0 {
           res = (res * a) % mod
       }
       a = (a * a) % mod
       e >>= 1
   }
   return res
}
