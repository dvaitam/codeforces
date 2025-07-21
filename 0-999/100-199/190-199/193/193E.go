package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var f int
   if _, err := fmt.Fscan(reader, &f); err != nil {
       return
   }
   const mod = 1013
   // Special cases for first occurrences
   if f == 0 {
       fmt.Println(0)
       return
   }
   if f == 1 {
       fmt.Println(1)
       return
   }
   // Track first occurrence positions, initialize to -1
   pos := make([]int, mod)
   for i := range pos {
       pos[i] = -1
   }
   // Initial Fibonacci values
   pos[0] = 0
   pos[1] = 1
   prev, curr := 0, 1
   // Generate sequence until cycle returns to 0,1
   for i := 2; ; i++ {
       next := (prev + curr) % mod
       if pos[next] == -1 {
           pos[next] = i
       }
       if next == f {
           fmt.Println(pos[next])
           return
       }
       prev, curr = curr, next
       if prev == 0 && curr == 1 {
           break
       }
   }
   // Not found in one full Pisano period
   fmt.Println(-1)
}
