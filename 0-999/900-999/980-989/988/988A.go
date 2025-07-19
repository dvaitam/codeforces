package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   ans := make([]int, 0, k)
   seen := make(map[int]bool, n)
   for i := 1; i <= n; i++ {
       var x int
       fmt.Fscan(reader, &x)
       if !seen[x] {
           seen[x] = true
           if len(ans) < k {
               ans = append(ans, i)
           }
       }
   }
   if len(ans) >= k {
       fmt.Println("YES")
       for i := 0; i < k; i++ {
           if i > 0 {
               fmt.Print(" ")
           }
           fmt.Print(ans[i])
       }
       fmt.Println()
   } else {
       fmt.Println("NO")
   }
}
