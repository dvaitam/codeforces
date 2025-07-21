package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // counts for values 1, 2, 3
   counts := [4]int{}
   for i := 0; i < n; i++ {
       var x int
       fmt.Fscan(reader, &x)
       if x >= 1 && x <= 3 {
           counts[x]++
       }
   }
   // find maximum frequency
   maxCount := counts[1]
   if counts[2] > maxCount {
       maxCount = counts[2]
   }
   if counts[3] > maxCount {
       maxCount = counts[3]
   }
   // minimal replacements = total - max frequency
   fmt.Fprintln(writer, n-maxCount)
}
