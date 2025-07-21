package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var t int
   if _, err := fmt.Fscan(in, &t); err != nil {
       return
   }
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for i := 0; i < t; i++ {
       var n, k, d1, d2 int64
       fmt.Fscan(in, &n, &k, &d1, &d2)
       // total wins must be divisible by 3
       if n%3 != 0 {
           fmt.Fprintln(out, "no")
           continue
       }
       target := n / 3
       ok := false
       // try all sign combinations for differences
       for _, s1 := range []int64{1, -1} {
           for _, s2 := range []int64{1, -1} {
               // w1 - w2 = s1*d1, w2 - w3 = s2*d2
               // sum w1+w2+w3 = k => 3*w2 + s1*d1 - s2*d2 = k
               num := k - s1*d1 + s2*d2
               if num%3 != 0 {
                   continue
               }
               w2 := num / 3
               w1 := w2 + s1*d1
               w3 := w2 - s2*d2
               if w1 < 0 || w2 < 0 || w3 < 0 {
                   continue
               }
               if w1 > target || w2 > target || w3 > target {
                   continue
               }
               ok = true
               break
           }
           if ok {
               break
           }
       }
       if ok {
           fmt.Fprintln(out, "yes")
       } else {
           fmt.Fprintln(out, "no")
       }
   }
}
