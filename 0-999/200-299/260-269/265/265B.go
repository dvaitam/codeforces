package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   // number of trees
   fmt.Fscan(reader, &n)
   heights := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &heights[i])
   }
   var total int64
   curH := 0
   for i := 0; i < n; i++ {
       if i > 0 {
           // jump to next tree
           total++
       }
       // climb up to the nut
       if heights[i] > curH {
           total += int64(heights[i] - curH)
           curH = heights[i]
       }
       // eat the nut
       total++
       // prepare height for next jump by descending if needed
       if i < n-1 && curH > heights[i+1] {
           total += int64(curH - heights[i+1])
           curH = heights[i+1]
       }
   }
   fmt.Println(total)
}
