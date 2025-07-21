package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, k int
   if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
       return
   }
   const mod = 1000000007
   ans := 0
   // Good arrays exist only for even length r = 2*h (h>=1)
   // For each h such that 2*h <= n, number of arrays = m - h
   maxH := n / 2
   for h := 1; h <= maxH; h++ {
       val := m - h
       if val <= 0 {
           break
       }
       ans = (ans + val) % mod
   }
   fmt.Println(ans)
}
