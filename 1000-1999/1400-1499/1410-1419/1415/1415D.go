package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // If n is large, answer is 1 by pigeonhole argument
   if n > 130 {
       fmt.Println(1)
       return
   }
   // prefix xor
   px := make([]int, n+1)
   for i := 1; i <= n; i++ {
       px[i] = px[i-1] ^ a[i-1]
   }
   const INF = int(1e9)
   ans := INF
   // brute force small n
   for l := 1; l <= n; l++ {
       for r := l + 1; r <= n; r++ {
           for i := l; i < r; i++ {
               left := px[i] ^ px[l-1]
               right := px[r] ^ px[i]
               if left > right {
                   // operations = (r-l+1) - 2 = r-l-1
                   ops := r - l - 1
                   if ops < ans {
                       ans = ops
                   }
               }
           }
       }
   }
   if ans == INF {
       fmt.Println(-1)
   } else {
       fmt.Println(ans)
   }
}
