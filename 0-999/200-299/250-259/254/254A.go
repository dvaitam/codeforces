package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// pair holds value and its original index
type pair struct {
   val int
   idx int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var m int
   if _, err := fmt.Fscan(reader, &m); err != nil {
       return
   }
   n := 2 * m
   data := make([]pair, n)
   count := make(map[int]int, n)
   for i := 0; i < n; i++ {
       var v int
       fmt.Fscan(reader, &v)
       data[i] = pair{val: v, idx: i + 1}
       count[v]++
   }
   for _, c := range count {
       if c%2 != 0 {
           fmt.Fprintln(writer, -1)
           return
       }
   }
   sort.Slice(data, func(i, j int) bool {
       return data[i].val < data[j].val
   })
   for i := 0; i < n; i += 2 {
       fmt.Fprintln(writer, data[i].idx, data[i+1].idx)
   }
}
