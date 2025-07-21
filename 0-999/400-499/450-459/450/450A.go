package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   ans := 1
   maxRounds := 0
   for i := 1; i <= n; i++ {
       var a int
       fmt.Fscan(reader, &a)
       rounds := (a + m - 1) / m
       // choose the child with the largest rounds, ties broken by larger index
       if rounds >= maxRounds {
           maxRounds = rounds
           ans = i
       }
   }
   fmt.Println(ans)
}
