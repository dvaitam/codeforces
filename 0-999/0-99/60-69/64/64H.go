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
   type participant struct {
       name  string
       score int
   }
   parts := make([]participant, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &parts[i].name, &parts[i].score)
   }
   sort.Slice(parts, func(i, j int) bool {
       if parts[i].score != parts[j].score {
           return parts[i].score > parts[j].score
       }
       return parts[i].name < parts[j].name
   })
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   // iterate and group ties
   i := 0
   for i < n {
       j := i + 1
       for j < n && parts[j].score == parts[i].score {
           j++
       }
       // positions i to j-1 inclusive share the same score
       place := ""
       if j-i == 1 {
           place = fmt.Sprintf("%d", i+1)
       } else {
           place = fmt.Sprintf("%d-%d", i+1, j)
       }
       // output each name in this group
       for k := i; k < j; k++ {
           fmt.Fprintf(writer, "%s %s\n", place, parts[k].name)
       }
       i = j
   }
}
