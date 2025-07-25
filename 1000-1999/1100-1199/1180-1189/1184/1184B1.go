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

   var s, b int
   if _, err := fmt.Fscan(reader, &s, &b); err != nil {
       return
   }
   a := make([]int, s)
   for i := 0; i < s; i++ {
       fmt.Fscan(reader, &a[i])
   }
   bases := make([][2]int, b)
   for i := 0; i < b; i++ {
       fmt.Fscan(reader, &bases[i][0], &bases[i][1])
   }
   // Sort bases by defensive power
   sort.Slice(bases, func(i, j int) bool {
       return bases[i][0] < bases[j][0]
   })
   defs := make([]int, b)
   prefix := make([]int64, b)
   for i := 0; i < b; i++ {
       defs[i] = bases[i][0]
       if i == 0 {
           prefix[i] = int64(bases[i][1])
       } else {
           prefix[i] = prefix[i-1] + int64(bases[i][1])
       }
   }
   // For each spaceship, binary search the total gold
   for i := 0; i < s; i++ {
       attack := a[i]
       idx := sort.Search(b, func(j int) bool {
           return defs[j] > attack
       })
       if idx == 0 {
           fmt.Fprint(writer, 0)
       } else {
           fmt.Fprint(writer, prefix[idx-1])
       }
       if i+1 < s {
           fmt.Fprint(writer, " ")
       }
   }
   fmt.Fprintln(writer)
}
