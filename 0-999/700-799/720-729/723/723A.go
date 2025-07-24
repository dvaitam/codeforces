package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var x [3]int
   fmt.Fscan(reader, &x[0], &x[1], &x[2])
   // Sort to find median
   sort.Ints(x[:])
   m := x[1]
   // Compute total distance to median
   ans := abs(x[0]-m) + abs(x[1]-m) + abs(x[2]-m)
   fmt.Fprintln(writer, ans)
}

// abs returns absolute value of integer
func abs(a int) int {
   if a < 0 {
       return -a
   }
   return a
}
