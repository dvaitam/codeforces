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
   a := make([]int, n)
   // ratings are between 1 and 2000
   freq := make([]int, 2001)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       freq[a[i]]++
   }
   // suffix[r] = number of students with rating > r
   suffix := make([]int, 2002)
   for r := 2000; r >= 1; r-- {
       suffix[r] = suffix[r+1] + freq[r+1]
   }
   for i, rating := range a {
       pos := suffix[rating] + 1
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, pos)
   }
   writer.WriteByte('\n')
}
