package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, s, t int
   if _, err := fmt.Fscan(reader, &n, &s, &t); err != nil {
       return
   }
   p := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &p[i])
   }
   // If starting and target positions are same
   if s == t {
       fmt.Println(0)
       return
   }
   // Follow the permutation cycle from s
   steps := 0
   curr := s
   for {
       curr = p[curr]
       steps++
       if curr == t {
           fmt.Println(steps)
           return
       }
       if curr == s {
           // Returned to start without reaching t
           fmt.Println(-1)
           return
       }
   }
}
