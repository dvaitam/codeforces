package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   rdr := bufio.NewReader(os.Stdin)
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()

   var n int
   if _, err := fmt.Fscan(rdr, &n); err != nil {
       return
   }
   type person struct {
       a, b, diff int64
   }
   people := make([]person, n)
   for i := 0; i < n; i++ {
       var ai, bi int64
       fmt.Fscan(rdr, &ai, &bi)
       people[i] = person{a: ai, b: bi, diff: ai - bi}
   }
   sort.Slice(people, func(i, j int) bool {
       return people[i].diff > people[j].diff
   })
   var total int64
   nn := int64(n)
   for i, p := range people {
       pos := int64(i)
       total += p.a*pos + p.b*(nn-pos-1)
   }
   fmt.Fprintln(w, total)
}
