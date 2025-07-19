package main

import (
   "fmt"
   "os"
)

func main() {
   var n int
   if _, err := fmt.Fscan(os.Stdin, &n); err != nil {
       return
   }
   mp := make(map[int]int)
   var s string
   for i := 0; i < n; i++ {
       fmt.Fscan(os.Stdin, &s)
       cnt0, cnt1 := 0, 0
       for j := 0; j < len(s); j++ {
           if s[j] == '(' {
               cnt0++
           } else {
               if cnt0 > 0 {
                   cnt0--
               } else {
                   cnt1++
               }
           }
       }
       if cnt0 == 0 && cnt1 == 0 {
           mp[0]++
       } else if cnt0 == 0 {
           mp[cnt1]++
       } else if cnt1 == 0 {
           mp[-cnt0]++
       }
   }

   res := 0
   for k, v := range mp {
       if k < 0 {
           if v2, ok := mp[-k]; ok {
               if v < v2 {
                   res += v
               } else {
                   res += v2
               }
           }
       }
   }
   // pairs of fully balanced strings
   res += mp[0] / 2
   fmt.Println(res)
}
