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
   fmt.Fscan(reader, &n)
   freq := make(map[int]int)
   for i := 0; i < n; i++ {
       var x int
       fmt.Fscan(reader, &x)
       freq[x]++
   }
   maxCount := 0
   for _, v := range freq {
       if v > maxCount {
           maxCount = v
       }
   }
   threshold := (n + 1) / 2
   if maxCount <= threshold {
       fmt.Fprint(writer, "YES")
   } else {
       fmt.Fprint(writer, "NO")
   }
}
