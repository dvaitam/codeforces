package main

import "fmt"

func main() {
   var a, b, c int64
   if _, err := fmt.Scan(&a, &b, &c); err != nil {
       return
   }
   // weekly consumption: fish=3, rabbit=2, chicken=2
   wf, wr, wc := int64(3), int64(2), int64(2)
   // maximum full weeks
   kf, kr, kc := a/wf, b/wr, c/wc
   k := kf
   if kr < k { k = kr }
   if kc < k { k = kc }
   // remaining after full weeks
   ra, rb, rc := a-k*wf, b-k*wr, c-k*wc
   // schedule: 0=fish, 1=rabbit, 2=chicken for Mon..Sun
   sched := []int{0, 1, 2, 0, 2, 1, 0}
   var maxExtra int64
   // try each starting day
   for s := 0; s < 7; s++ {
       rem := [3]int64{ra, rb, rc}
       var extra int64
       for i := 0; i < 7; i++ {
           t := sched[(s+i)%7]
           if rem[t] > 0 {
               rem[t]--
               extra++
           } else {
               break
           }
       }
       if extra > maxExtra {
           maxExtra = extra
       }
   }
   // total days = full weeks * 7 + extra days
   fmt.Println(k*7 + maxExtra)
}
