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
   events := make(map[int]int)
   for i := 0; i < n; i++ {
       var b, d int
       fmt.Fscan(reader, &b, &d)
       events[b]++
       events[d]--
   }
   years := make([]int, 0, len(events))
   for y := range events {
       years = append(years, y)
   }
   sort.Ints(years)
   cur := 0
   maxCount := 0
   maxYear := 0
   for _, y := range years {
       cur += events[y]
       if cur > maxCount {
           maxCount = cur
           maxYear = y
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintf(writer, "%d %d\n", maxYear, maxCount)
}
