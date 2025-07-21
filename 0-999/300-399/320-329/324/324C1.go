package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   line, err := reader.ReadString('\n')
   if err != nil && err.Error() != "EOF" {
       fmt.Fprintln(os.Stderr, "Error reading input:", err)
       return
   }
   line = line
   // trim whitespace
   s := line
   // remove possible \r and \n
   if len(s) > 0 && (s[len(s)-1] == '\n' || s[len(s)-1] == '\r') {
       s = s[:len(s)-1]
   }
   if len(s) > 0 && (s[len(s)-1] == '\r') {
       s = s[:len(s)-1]
   }
   n, err := strconv.Atoi(s)
   if err != nil {
       fmt.Fprintln(os.Stderr, "Invalid number:", err)
       return
   }
   if n <= 0 {
       fmt.Println(0)
       return
   }
   dp := make([]int, n+1)
   dp[0] = 0
   const INF = 1<<30
   for i := 1; i <= n; i++ {
       x := i
       best := INF
       for x > 0 {
           d := x % 10
           x /= 10
           if d > 0 {
               prev := dp[i-d]
               if prev < best {
                   best = prev
               }
           }
       }
       // best must have been set since i>=1 and has at least one non-zero digit
       dp[i] = best + 1
   }
   fmt.Println(dp[n])
}
