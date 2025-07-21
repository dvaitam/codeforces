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
   p := make([]int, n)
   pos := make([]int, n+1)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &p[i])
       pos[p[i]] = i
   }
   const M = 300
   for d := 1; d <= M; d++ {
       // for c-d >=1 and c+d <= n
       for c := d + 1; c + d <= n; c++ {
           i := pos[c-d]
           j := pos[c]
           k := pos[c+d]
           if (i < j && j < k) || (i > j && j > k) {
               fmt.Println("YES")
               return
           }
       }
   }
   fmt.Println("NO")
}
