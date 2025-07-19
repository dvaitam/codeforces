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
   var sum int
   var lim int
   if _, err := fmt.Fscan(in, &sum, &lim); err != nil {
       return
   }
   const maxb = 31
   s := make([]int, maxb)
   for i := 1; i <= lim; i++ {
       // count lowest set bit of i
       for j := 0; j < maxb; j++ {
           if i&(1<<j) != 0 {
               s[j]++
               break
           }
       }
   }
   num := make([]int, maxb)
   f := false
   // for each bit in sum
   for i := 0; i < maxb; i++ {
       if (sum>>i)&1 == 0 {
           continue
       }
       need := 1
       for j := i; j >= 0; j-- {
           if s[j] >= need {
               num[j] += need
               s[j] -= need
               break
           } else {
               need -= s[j]
               num[j] += s[j]
               need <<= 1
               s[j] = 0
               if j == 0 {
                   f = true
                   break
               }
           }
       }
       if f {
           break
       }
   }
   if f {
       fmt.Fprintln(out, -1)
       return
   }
   total := 0
   for _, v := range num {
       total += v
   }
   fmt.Fprintln(out, total)
   // output chosen numbers
   for i := 1; i <= lim; i++ {
       for j := 0; j < maxb; j++ {
           if i&(1<<j) != 0 && num[j] > 0 {
               fmt.Fprint(out, i, " ")
               num[j]--
               break
           }
       }
   }
   fmt.Fprintln(out)
}
