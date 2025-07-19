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

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   skills := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &skills[i])
   }
   sort.Ints(skills)
   left := 0
   maxSize := 1
   for i := 1; i < n; i++ {
       for skills[i]-skills[left] > 5 {
           left++
       }
       size := i - left + 1
       if size > maxSize {
           maxSize = size
       }
   }
   fmt.Fprintln(writer, maxSize)
}
