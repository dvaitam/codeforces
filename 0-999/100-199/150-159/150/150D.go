package main

import (
   "bufio"
   "fmt"
   "os"
)

func max(a, b int64) int64 {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var l int
   if _, err := fmt.Fscan(reader, &l); err != nil {
       return
   }
   a := make([]int64, l+1)
   for i := 1; i <= l; i++ {
       var ai int64
       fmt.Fscan(reader, &ai)
       a[i] = ai
   }
   var s string
   fmt.Fscan(reader, &s)
   // dp1[i][j]: max score to fully delete s[i..j], or -inf if impossible
   const negInf = -1 << 60
   dp1 := make([][]int64, l)
   for i := range dp1 {
       dp1[i] = make([]int64, l)
       for j := range dp1[i] {
           dp1[i][j] = negInf
       }
   }
   // base length 1
   for i := 0; i < l; i++ {
       if a[1] >= 0 {
           dp1[i][i] = a[1]
       }
   }
   // lengths 2..l
   for length := 2; length <= l; length++ {
       for i := 0; i+length-1 < l; i++ {
           j := i + length - 1
           // delete whole [i..j] if palindrome and interior deletable
           if s[i] == s[j] && a[length] >= 0 {
               if length == 2 {
                   dp1[i][j] = max(dp1[i][j], a[length])
               } else if dp1[i+1][j-1] >= 0 {
                   dp1[i][j] = max(dp1[i][j], dp1[i+1][j-1] + a[length])
               }
           }
           // split
           for k := i; k < j; k++ {
               if dp1[i][k] >= 0 && dp1[k+1][j] >= 0 {
                   dp1[i][j] = max(dp1[i][j], dp1[i][k] + dp1[k+1][j])
               }
           }
       }
   }
   // dp2[i]: max score for prefix s[0..i-1]
   dp2 := make([]int64, l+1)
   // dp2[0] = 0
   for i := 1; i <= l; i++ {
       // leave char i-1
       dp2[i] = dp2[i-1]
       // try deleting a suffix j..i-1 fully
       for j := 0; j < i; j++ {
           if dp1[j][i-1] >= 0 {
               dp2[i] = max(dp2[i], dp2[j] + dp1[j][i-1])
           }
       }
   }
   fmt.Println(dp2[l])
