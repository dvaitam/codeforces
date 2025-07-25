package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, b, a int
   if _, err := fmt.Fscan(reader, &n, &b, &a); err != nil {
       return
   }
   // maximum capacity of accumulator
   maxA := a
   // read sun exposure array
   s := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &s[i])
   }
   // current charges
   curB, curA := b, a
   passed := 0
   for i := 0; i < n; i++ {
       if s[i] == 1 {
           // sunny segment: prefer battery if can recharge accumulator
           if curB > 0 && curA < maxA {
               curB--
               curA++
           } else if curA > 0 {
               curA--
           } else if curB > 0 {
               curB--
           } else {
               break
           }
       } else {
           // not sunny: prefer accumulator
           if curA > 0 {
               curA--
           } else if curB > 0 {
               curB--
           } else {
               break
           }
       }
       passed++
   }
   fmt.Println(passed)
}
