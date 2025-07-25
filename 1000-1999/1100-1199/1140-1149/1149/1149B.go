package main

import (
   "bufio"
   "fmt"
   "os"
)

const MAX = 251

var (
   n, q    int
   s       string
   nxt     [][]int
   dp      [MAX][MAX][MAX]int
   A, B, C [MAX]byte
   la, lb, lc int
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fscan(reader, &n, &q)
   fmt.Fscan(reader, &s)
   // build nxt array, size n+1 by 26
   nxt = make([][]int, n+2)
   for i := range nxt {
       nxt[i] = make([]int, 26)
   }
   // initialize next positions to n+1
   for c := 0; c < 26; c++ {
       nxt[n][c] = n + 1
       nxt[n+1][c] = n + 1
   }
   // s is 0-indexed, positions 1..n
   for i := n - 1; i >= 0; i-- {
       for c := 0; c < 26; c++ {
           nxt[i][c] = nxt[i+1][c]
       }
       nxt[i][s[i]-'a'] = i + 1
   }
   // initialize dp[0][0][0] = 0
   for i := 0; i < MAX; i++ {
       for j := 0; j < MAX; j++ {
           for k := 0; k < MAX; k++ {
               dp[i][j][k] = n + 1
           }
       }
   }
   dp[0][0][0] = 0

   for qi := 0; qi < q; qi++ {
       var op string
       var idx int
       fmt.Fscan(reader, &op, &idx)
       if op == "+" {
           var ch byte
           fmt.Fscan(reader, &ch)
           if idx == 1 {
               la++
               A[la] = ch
               i := la
               for j := 0; j <= lb; j++ {
                   for k := 0; k <= lc; k++ {
                       dp[i][j][k] = n + 1
                       // from A
                       prev := dp[i-1][j][k]
                       if prev <= n {
                           pos := nxt[prev][A[i]-'a']
                           if pos <= n {
                               dp[i][j][k] = min(dp[i][j][k], pos)
                           }
                       }
                       // from B
                       if j > 0 {
                           prev = dp[i][j-1][k]
                           if prev <= n {
                               pos := nxt[prev][B[j]-'a']
                               if pos <= n {
                                   dp[i][j][k] = min(dp[i][j][k], pos)
                               }
                           }
                       }
                       // from C
                       if k > 0 {
                           prev = dp[i][j][k-1]
                           if prev <= n {
                               pos := nxt[prev][C[k]-'a']
                               if pos <= n {
                                   dp[i][j][k] = min(dp[i][j][k], pos)
                               }
                           }
                       }
                   }
               }
           } else if idx == 2 {
               lb++
               B[lb] = ch
               j := lb
               for i := 0; i <= la; i++ {
                   for k := 0; k <= lc; k++ {
                       dp[i][j][k] = n + 1
                       // from A
                       if i > 0 {
                           prev := dp[i-1][j][k]
                           if prev <= n {
                               pos := nxt[prev][A[i]-'a']
                               if pos <= n {
                                   dp[i][j][k] = min(dp[i][j][k], pos)
                               }
                           }
                       }
                       // from B
                       prev := dp[i][j-1][k]
                       if prev <= n {
                           pos := nxt[prev][B[j]-'a']
                           if pos <= n {
                               dp[i][j][k] = min(dp[i][j][k], pos)
                           }
                       }
                       // from C
                       if k > 0 {
                           prev = dp[i][j][k-1]
                           if prev <= n {
                               pos := nxt[prev][C[k]-'a']
                               if pos <= n {
                                   dp[i][j][k] = min(dp[i][j][k], pos)
                               }
                           }
                       }
                   }
               }
           } else {
               lc++
               C[lc] = ch
               k := lc
               for i := 0; i <= la; i++ {
                   for j := 0; j <= lb; j++ {
                       dp[i][j][k] = n + 1
                       // from A
                       if i > 0 {
                           prev := dp[i-1][j][k]
                           if prev <= n {
                               pos := nxt[prev][A[i]-'a']
                               if pos <= n {
                                   dp[i][j][k] = min(dp[i][j][k], pos)
                               }
                           }
                       }
                       // from B
                       if j > 0 {
                           prev := dp[i][j-1][k]
                           if prev <= n {
                               pos := nxt[prev][B[j]-'a']
                               if pos <= n {
                                   dp[i][j][k] = min(dp[i][j][k], pos)
                               }
                           }
                       }
                       // from C
                       prev := dp[i][j][k-1]
                       if prev <= n {
                           pos := nxt[prev][C[k]-'a']
                           if pos <= n {
                               dp[i][j][k] = min(dp[i][j][k], pos)
                           }
                       }
                   }
               }
           }
       } else {
           // remove last
           if idx == 1 {
               la--
           } else if idx == 2 {
               lb--
           } else {
               lc--
           }
       }
       if dp[la][lb][lc] <= n {
           fmt.Fprintln(writer, "YES")
       } else {
           fmt.Fprintln(writer, "NO")
       }
   }
}
