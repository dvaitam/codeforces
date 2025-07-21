package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   cnt0, cnt5 := 0, 0
   for i := 0; i < n; i++ {
       var x int
       if _, err := fmt.Fscan(reader, &x); err != nil {
           return
       }
       if x == 0 {
           cnt0++
       } else if x == 5 {
           cnt5++
       }
   }
   // Need at least one zero for divisibility by 10
   if cnt0 == 0 {
       fmt.Println(-1)
       return
   }
   // Need sum of digits divisible by 9 => at least 9 fives
   if cnt5 < 9 {
       // can only make 0
       fmt.Println(0)
       return
   }
   // Use the maximum multiple of 9 fives
   k := (cnt5 / 9) * 9
   // Build result: k times '5', then cnt0 times '0'
   res := make([]byte, k+cnt0)
   for i := 0; i < k; i++ {
       res[i] = '5'
   }
   for i := k; i < k+cnt0; i++ {
       res[i] = '0'
   }
   fmt.Println(string(res))
}
