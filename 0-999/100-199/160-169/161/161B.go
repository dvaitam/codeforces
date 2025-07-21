package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   votes := make([]int, n)
   sum := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &votes[i])
       sum += votes[i]
   }
   if votes[0] > sum-votes[0] {
       fmt.Println("Yes")
   } else {
       fmt.Println("No")
   }
}
