package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   fmt.Fscan(reader, &t)
   for t > 0 {
       t--
       var n int
       fmt.Fscan(reader, &n)
       w := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &w[i])
       }
       // frequency of weights (1..n)
       freq := make([]int, n+1)
       for _, v := range w {
           if v >= 1 && v <= n {
               freq[v]++
           }
       }
       best := 0
       // try all possible sums s
       for s := 2; s <= 2*n; s++ {
           cnt := 0
           // pairs with distinct weights i and j where i<j and i+j==s
           for i := 1; i < s-i && i <= n; i++ {
               j := s - i
               if j > n {
                   continue
               }
               cnt += min(freq[i], freq[j])
           }
           // pairs with equal weights i == j when s even
           if s%2 == 0 {
               i := s / 2
               if i >= 1 && i <= n {
                   cnt += freq[i] / 2
               }
           }
           if cnt > best {
               best = cnt
           }
       }
       fmt.Fprintln(writer, best)
   }
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}
