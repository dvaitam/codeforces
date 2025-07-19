package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type Info struct {
   name string
   p, h  int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   v := make([]Info, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &v[i].name, &v[i].p)
   }
   // sort by p ascending
   sort.Slice(v, func(i, j int) bool {
       return v[i].p < v[j].p
   })
   // track occupied heights
   id := make([]bool, n+5)
   for i := n - 1; i >= 0; i-- {
       x := (i - v[i].p) + 1
       p := 1
       if x <= 0 {
           fmt.Fprintln(writer, -1)
           return
       }
       // find x-th free slot
       for x > 0 || id[p] {
           if !id[p] {
               x--
           }
           if x == 0 {
               break
           }
           p++
       }
       id[p] = true
       v[i].h = p
   }
   // output
   for i := 0; i < n; i++ {
       fmt.Fprintln(writer, v[i].name, v[i].h)
   }
}
