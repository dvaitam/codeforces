package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // Precompute lengths of increasing segments
   dp1 := make([]int, n) // dp1[i]: length of strictly increasing ending at i
   dp2 := make([]int, n) // dp2[i]: length of strictly increasing starting at i
   if n > 0 {
       dp1[0] = 1
       for i := 1; i < n; i++ {
           if a[i] > a[i-1] {
               dp1[i] = dp1[i-1] + 1
           } else {
               dp1[i] = 1
           }
       }
       dp2[n-1] = 1
       for i := n-2; i >= 0; i-- {
           if a[i] < a[i+1] {
               dp2[i] = dp2[i+1] + 1
           } else {
               dp2[i] = 1
           }
       }
   }
   // Compute answer
   maxLen := 1
   // no change
   for i := 0; i < n; i++ {
       if dp1[i] > maxLen {
           maxLen = dp1[i]
       }
   }
   // change first or last element
   if n > 1 {
       if dp2[1]+1 > maxLen {
           maxLen = dp2[1] + 1
       }
       if dp1[n-2]+1 > maxLen {
           maxLen = dp1[n-2] + 1
       }
   }
   // change one middle element to connect segments
   for i := 1; i < n-1; i++ {
       if a[i+1]-a[i-1] >= 2 {
           length := dp1[i-1] + dp2[i+1] + 1
           if length > maxLen {
               maxLen = length
           }
       }
   }
   fmt.Println(maxLen)
}
