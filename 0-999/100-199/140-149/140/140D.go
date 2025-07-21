package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
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
   sort.Ints(a)
   // Total available coding time after initial 10 minutes: 710
   const totalTime = 710
   // Free window before penalty: 350 writing minutes (350+10 = 360 contest minutes)
   const freeWindow = 350

   count := 0
   sumA := 0
   // select max problems
   for _, t := range a {
       if sumA + t > totalTime {
           break
       }
       sumA += t
       count++
   }
   // compute penalty
   penalty := 0
   sumA = 0
   for i := 0; i < count; i++ {
       sumA += a[i]
       if sumA > freeWindow {
           penalty += sumA - freeWindow
       }
   }
   // output number of problems and total penalty
   fmt.Printf("%d %d\n", count, penalty)
}
