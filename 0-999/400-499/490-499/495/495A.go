package main

import (
   "fmt"
   "os"
)

func main() {
   var s string
   if _, err := fmt.Fscan(os.Stdin, &s); err != nil {
       return
   }
   if len(s) != 2 {
       return
   }
   mask := [10]int{0x3F, 0x06, 0x5B, 0x4F, 0x66, 0x6D, 0x7D, 0x07, 0x7F, 0x6F}
   obs1 := mask[s[0]-'0']
   obs2 := mask[s[1]-'0']
   cnt := 0
   for x := 0; x < 100; x++ {
       d1 := x / 10
       d2 := x % 10
       if (obs1 & mask[d1]) == obs1 && (obs2 & mask[d2]) == obs2 {
           cnt++
       }
   }
   fmt.Println(cnt)
}
