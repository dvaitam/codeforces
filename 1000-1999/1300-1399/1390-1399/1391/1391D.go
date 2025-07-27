package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   // If both dimensions >=4, impossible
   if n >= 4 && m >= 4 {
       fmt.Println(-1)
       return
   }
   // Read matrix and build column masks
   col := make([]int, m)
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       for j := 0; j < m; j++ {
           if s[j] == '1' {
               col[j] |= 1 << i
           }
       }
   }
   // No constraints if only one row
   if n == 1 {
       fmt.Println(0)
       return
   }
   // DP over columns for n=2 or n=3
   states := 1 << n
   // Precompute valid transitions
   ok := make([][]bool, states)
   for u := 0; u < states; u++ {
       ok[u] = make([]bool, states)
       for v := 0; v < states; v++ {
           valid := true
           // for each adjacent row pair
           for i := 0; i+1 < n; i++ {
               sum := ((u>>i)&1) + ((u>>(i+1))&1) + ((v>>i)&1) + ((v>>(i+1))&1)
               if sum&1 == 0 {
                   valid = false
                   break
               }
           }
           ok[u][v] = valid
       }
   }
   const INF = 1000000000
   prev := make([]int, states)
   cur := make([]int, states)
   // Initialize DP for first column
   for v := 0; v < states; v++ {
       prev[v] = bits.OnesCount(uint(v ^ col[0]))
   }
   // Process remaining columns
   for j := 1; j < m; j++ {
       for v := 0; v < states; v++ {
           best := INF
           for u := 0; u < states; u++ {
               if ok[u][v] && prev[u] < best {
                   best = prev[u]
               }
           }
           cur[v] = best + bits.OnesCount(uint(v ^ col[j]))
       }
       // swap prev and cur
       prev, cur = cur, prev
   }
   // Answer is min over last column
   ans := INF
   for v := 0; v < states; v++ {
       if prev[v] < ans {
           ans = prev[v]
       }
   }
   fmt.Println(ans)
}
