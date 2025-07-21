package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m, minVal, maxVal int
   if _, err := fmt.Fscan(reader, &n, &m, &minVal, &maxVal); err != nil {
       return
   }
   temps := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &temps[i])
   }
   // Check recorded temperatures in bounds
   hasMin := false
   hasMax := false
   for _, t := range temps {
       if t < minVal || t > maxVal {
           fmt.Println("Incorrect")
           return
       }
       if t == minVal {
           hasMin = true
       }
       if t == maxVal {
           hasMax = true
       }
   }
   // Need to add missing min and max if not present
   needed := 0
   if !hasMin {
       needed++
   }
   if !hasMax {
       needed++
   }
   // Check if we have enough slots to add
   if needed <= (n - m) {
       fmt.Println("Correct")
   } else {
       fmt.Println("Incorrect")
   }
}
