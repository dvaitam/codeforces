package main

import (
   "bufio"
   "fmt"
   "os"
)

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}

func reverseBytes(s []byte) {
   for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
       s[i], s[j] = s[j], s[i]
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var T int
   fmt.Fscan(reader, &T)
   for tc := 0; tc < T; tc++ {
       var n, a, b int
       fmt.Fscan(reader, &n, &a, &b)
       var s string
       fmt.Fscan(reader, &s)
       // dp[i][j][k][l] = reachable
       dp := make([][][][]bool, n+1)
       for i := 0; i <= n; i++ {
           dp[i] = make([][][]bool, a)
           for j := 0; j < a; j++ {
               dp[i][j] = make([][]bool, b)
               for k := 0; k < b; k++ {
                   dp[i][j][k] = make([]bool, n+1)
               }
           }
       }
       dp[0][0][0][0] = true
       // DP transition
       for i := 0; i < n; i++ {
           x := int(s[i] - '0')
           for j := 0; j < a; j++ {
               for k := 0; k < b; k++ {
                   for l := 0; l <= i; l++ {
                       if !dp[i][j][k][l] {
                           continue
                       }
                       // assign R
                       nj := (j*10 + x) % a
                       dp[i+1][nj][k][l+1] = true
                       // assign B
                       nk := (k*10 + x) % b
                       dp[i+1][j][nk][l] = true
                   }
               }
           }
       }
       // check possible
       possible := false
       for l := 1; l < n; l++ {
           if dp[n][0][0][l] {
               possible = true
               break
           }
       }
       if !possible {
           fmt.Fprintln(writer, -1)
           continue
       }
       // find best split
       // choose red count to minimize difference
       bestL := -1
       bestDiff := n
       for l1 := 1; l1 < n; l1++ {
           if dp[n][0][0][l1] {
               d := abs(2*l1 - n)
               if bestL < 0 || d < bestDiff {
                   bestDiff = d
                   bestL = l1
               }
           }
       }
       // backtrack
       remA, remB := 0, 0
       l := bestL
       res := make([]byte, 0, n)
       for i := n - 1; i >= 0; i-- {
           x := int(s[i] - '0')
           // try R
           found := false
           if l > 0 {
               for j := 0; j < a; j++ {
                   if (j*10+x)%a == remA && dp[i][j][remB][l-1] {
                       remA = j
                       l--
                       res = append(res, 'R')
                       found = true
                       break
                   }
               }
           }
           if !found {
               // B
               for k := 0; k < b; k++ {
                   if (k*10+x)%b == remB && dp[i][remA][k][l] {
                       remB = k
                       res = append(res, 'B')
                       break
                   }
               }
           }
       }
       reverseBytes(res)
       writer.Write(res)
       writer.WriteByte('\n')
   }
}
