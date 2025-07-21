package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000003

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   n := len(s)
   // prefix sum of non-digit count
   nonD := make([]int, n+1)
   for i := 0; i < n; i++ {
       nonD[i+1] = nonD[i]
       if s[i] < '0' || s[i] > '9' {
           nonD[i+1]++
       }
   }
   // dp[i][j] number of parses of s[i..j]
   dp := make([][]int, n)
   for i := range dp {
       dp[i] = make([]int, n)
   }
   // process lengths
   for length := 1; length <= n; length++ {
       for i := 0; i+length-1 < n; i++ {
           j := i + length - 1
           // if substring all digits
           if nonD[j+1]-nonD[i] == 0 {
               dp[i][j] = 1
               continue
           }
           ways := 0
           // unary
           if (s[i] == '+' || s[i] == '-') && i+1 <= j {
               ways = dp[i+1][j]
           }
           // binary splits
           for k := i + 1; k < j; k++ {
               c := s[k]
               if c == '+' || c == '-' || c == '*' || c == '/' {
                   left := dp[i][k-1]
                   if left != 0 {
                       right := dp[k+1][j]
                       if right != 0 {
                           ways = (ways + left*right) % MOD
                       }
                   }
               }
           }
           dp[i][j] = ways
       }
   }
   fmt.Println(dp[0][n-1] % MOD)
}
