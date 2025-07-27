package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var t int
   fmt.Fscan(in, &t)
   for i := 0; i < t; i++ {
       var p, q int64
       fmt.Fscan(in, &p, &q)
       if p%q != 0 {
           fmt.Fprintln(out, p)
           continue
       }
       // factorize q
       q0 := q
       factors := make(map[int64]int)
       for d := int64(2); d*d <= q0; d++ {
           for q0%d == 0 {
               factors[d]++
               q0 /= d
           }
       }
       if q0 > 1 {
           factors[q0]++
       }

       var ans int64 = 1
       // for each prime factor, compute candidate
       for f, e := range factors {
           // count exponent of f in p
           tmp := p
           var k int
           for tmp%f == 0 {
               tmp /= f
               k++
           }
           // remove just enough to make it not divisible by q
           removeExp := int64(k - e + 1)
           // compute divisor = f^removeExp
           div := int64(1)
           for j := int64(0); j < removeExp; j++ {
               div *= f
           }
           cand := p / div
           if cand > ans {
               ans = cand
           }
       }
       fmt.Fprintln(out, ans)
   }
}
