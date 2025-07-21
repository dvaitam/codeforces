package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   var m, k int64
   fmt.Fscan(reader, &n, &m, &k)
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // If n is even, total sum constrained, cannot steal
   if n%2 == 0 {
       fmt.Println(0)
       return
   }
   // Compute auxiliary t such that a[i] = t[i] + sign[i]*x
   t := make([]int64, n)
   // t[0] = 0 implicitly
   for i := 1; i < n; i++ {
       // S[i-1] = a[i-1] + a[i]
       t[i] = a[i-1] + a[i] - t[i-1]
   }
   // Compute minimum allowed x (L) to keep all a[i] >= 0
   // For odd positions (i%2==0), require t[i] + x >= 0 => x >= -t[i]
   var L int64 = -1e18
   for i := 0; i < n; i++ {
       if i%2 == 0 {
           if v := -t[i]; v > L {
               L = v
           }
       }
   }
   x0 := a[0]
   // Maximum total reduction of x
   maxRed := x0 - L
   if maxRed < 0 {
       maxRed = 0
   }
   // Operations cost to change x by 1: floor(n/2)+1
   cost := int64(n/2 + 1)
   // Maximum reduction in x per minute
   per := m / cost
   if per <= 0 {
       fmt.Println(0)
       return
   }
   // Total possible reduction over k minutes
   total := per * k
   if total > maxRed {
       total = maxRed
   }
   // Steal equals total reduction in x when n is odd
   fmt.Println(total)
}
