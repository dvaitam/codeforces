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
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   last := a[n-1]
   // If only one day observed
   if n == 1 {
       switch last {
       case 0:
           fmt.Println("UP")
       case 15:
           fmt.Println("DOWN")
       default:
           fmt.Println(-1)
       }
       return
   }
   // Multiple days observed
   // Next day after full moon or new moon is determinate
   if last == 0 {
       fmt.Println("UP")
       return
   }
   if last == 15 {
       fmt.Println("DOWN")
       return
   }
   prev := a[n-2]
   if last > prev {
       fmt.Println("UP")
   } else {
       fmt.Println("DOWN")
   }
}
