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

   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   Base := [64]uint64{}
   use := [64]uint64{}
   id := [64]int{}
   ans := make([]bool, n)
   var sum uint64
   var tot int

   for i := 0; i < n; i++ {
       var x uint64
       fmt.Fscan(in, &x)
       sum ^= x
       var p uint64
       for j := 63; j >= 0; j-- {
           if (x>>uint(j))&1 == 0 {
               continue
           }
           if Base[j] != 0 {
               x ^= Base[j]
               p ^= use[j]
           } else {
               Base[j] = x
               use[j] = p | (1 << uint(tot))
               id[tot] = i
               tot++
               break
           }
       }
   }
   // check function closure
   check := func(x uint64) bool {
       var b [64]uint64
       for i := 0; i < 64; i++ {
           y := Base[i] & x
           for j := 63; j >= 0; j-- {
               if (y>>uint(j))&1 == 0 {
                   continue
               }
               if b[j] != 0 {
                   y ^= b[j]
               } else {
                   b[j] = y
                   break
               }
           }
       }
       for i := 63; i >= 0; i-- {
           if (x>>uint(i))&1 == 0 {
               continue
           }
           if b[i] != 0 {
               x ^= b[i]
           } else {
               return false
           }
       }
       return true
   }

   var x uint64
   // first bits where sum bit is 0
   for i := 63; i >= 0; i-- {
       if (sum>>uint(i))&1 == 0 {
           if check(x | (1 << uint(i))) {
               x |= 1 << uint(i)
           }
       }
   }
   // then bits where sum bit is 1
   for i := 63; i >= 0; i-- {
       if (sum>>uint(i))&1 == 1 {
           if check(x | (1 << uint(i))) {
               x |= 1 << uint(i)
           }
       }
   }
   var p uint64
   for i := 63; i >= 0; i-- {
       if (x>>uint(i))&1 == 1 {
           x ^= Base[i]
           p ^= use[i]
       }
   }
   for i := 0; i < 64; i++ {
       if (p>>uint(i))&1 == 1 {
           idx := id[i]
           ans[idx] = true
       }
   }
   // output
   for i := 0; i < n; i++ {
       if ans[i] {
           out.WriteString("2 ")
       } else {
           out.WriteString("1 ")
       }
   }
   out.WriteByte('\n')
}
