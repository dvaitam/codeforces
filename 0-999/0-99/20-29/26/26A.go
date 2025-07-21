package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   var n int
   in := bufio.NewReader(os.Stdin)
   fmt.Fscan(in, &n)
   // Compute smallest prime factor (spf) for each number up to n
   spf := make([]int, n+1)
   for i := 2; i <= n; i++ {
       if spf[i] == 0 {
           for j := i; j <= n; j += i {
               if spf[j] == 0 {
                   spf[j] = i
               }
           }
       }
   }
   // Count numbers with exactly two distinct prime divisors
   count := 0
   for i := 2; i <= n; i++ {
       x := i
       prev := 0
       distinct := 0
       for x > 1 {
           p := spf[x]
           if p != prev {
               distinct++
               prev = p
           }
           x /= p
       }
       if distinct == 2 {
           count++
       }
   }
   fmt.Println(count)
}
