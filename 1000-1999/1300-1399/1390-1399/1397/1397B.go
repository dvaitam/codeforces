package main

import (
   "bufio"
   "fmt"
   "math"
   "math/big"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
   // initial best cost for c = 1
   best := int64(0)
   for i := 0; i < n; i++ {
       // target 1^i = 1
       diff := a[i] - 1
       if diff < 0 {
           diff = -diff
       }
       best += diff
   }
   // try bases c >= 2
   if n > 1 {
       // maximum power c^(n-1) that can improve over best
       maxPower := best + a[n-1]
       // approximate max base
       maxC := int(math.Pow(float64(maxPower), 1.0/float64(n-1))) + 1
       for c := 2; c <= maxC; c++ {
           // compute cost for this base
           cost := int64(0)
           cur := big.NewInt(1)
           bigC := big.NewInt(int64(c))
           valid := true
           for i := 0; i < n; i++ {
               if cur.BitLen() > 62 {
                   valid = false
                   break
               }
               tar := cur.Int64()
               diff := a[i] - tar
               if diff < 0 {
                   diff = -diff
               }
               cost += diff
               if cost > best {
                   valid = false
                   break
               }
               if i < n-1 {
                   cur.Mul(cur, bigC)
               }
           }
           if valid && cost < best {
               best = cost
           }
       }
   }
   fmt.Println(best)
}
