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
   var T int
   fmt.Fscan(reader, &T)
   for T > 0 {
       T--
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       sort.Ints(a)
       // Maximum k is limited by second-largest plank length minus one, and by total planks
       maxByLength := a[n-2] - 1
       maxByCount := n - 2
       if maxByLength < maxByCount {
           fmt.Fprintln(writer, maxByLength)
       } else {
           fmt.Fprintln(writer, maxByCount)
       }
   }
}
