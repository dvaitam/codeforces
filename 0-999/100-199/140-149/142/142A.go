package main

import (
   "fmt"
   "math"
)

// divisors returns all positive divisors of n
func divisors(n int64) []int64 {
   var ds []int64
   // iterate up to sqrt(n)
   lim := int64(math.Sqrt(float64(n)))
   for i := int64(1); i <= lim; i++ {
       if n%i == 0 {
           ds = append(ds, i)
           if j := n / i; j != i {
               ds = append(ds, j)
           }
       }
   }
   return ds
}

func main() {
   var n int64
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   // initialize min and max stolen blocks
   const inf64 = 1<<63 - 1
   minStolen := inf64
   maxStolen := int64(0)
   // iterate possible x = A-1
   for _, x := range divisors(n) {
       m := n / x
       // iterate possible y = B-2
       for _, y := range divisors(m) {
           z := m / y // z = C-2
           // stolen S = (x+1)*(y+2)*(z+2) - x*y*z
           // expand to avoid overflow: all within int64
           // S = 2*x*y + 2*x*z + 4*x + y*z + 2*y + 2*z + 4
           s := 2*x*y + 2*x*z + 4*x + y*z + 2*y + 2*z + 4
           if s < minStolen {
               minStolen = s
           }
           if s > maxStolen {
               maxStolen = s
           }
       }
   }
   fmt.Printf("%d %d\n", minStolen, maxStolen)
}
