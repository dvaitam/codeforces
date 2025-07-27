package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var T int64
   fmt.Fscan(in, &n, &T)
   var s string
   fmt.Fscan(in, &s)
   // v_i = 2^pos(s[i])
   // last character forced +, second last forced -
   // rest can choose signs
   // compute T' = T - v_n + v_{n-1}
   pos := func(c byte) int {
       return int(c - 'a')
   }
   vLast := int64(1) << pos(s[n-1])
   vSecond := int64(1) << pos(s[n-2])
   x := T - vLast + vSecond
   // collect counts of rest
   var cnt [26]int64
   for i := 0; i < n-2; i++ {
       cnt[pos(s[i])]++
   }
   // sum of rest
   var sumRest int64
   for i := 0; i < 26; i++ {
       if cnt[i] > 0 {
           sumRest += cnt[i] * (int64(1) << i)
       }
   }
   // need (x + sumRest) / 2
   // must be non-negative integer
   if (x+sumRest)%2 != 0 {
       fmt.Println("No")
       return
   }
   need := (x + sumRest) / 2
   if need < 0 {
       fmt.Println("No")
       return
   }
   // greedy take from largest
   for i := 25; i >= 0; i-- {
       if need <= 0 {
           break
       }
       if cnt[i] == 0 {
           continue
       }
       val := int64(1) << i
       // max usable
       maxUse := need / val
       if maxUse > cnt[i] {
           maxUse = cnt[i]
       }
       need -= maxUse * val
   }
   if need == 0 {
       fmt.Println("Yes")
   } else {
       fmt.Println("No")
   }
}
