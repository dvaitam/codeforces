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

   var ns, nm, nl, nxl, nxxl int
   fmt.Fscan(reader, &ns, &nm, &nl, &nxl, &nxxl)
   counts := []int{ns, nm, nl, nxl, nxxl}

   var k int
   fmt.Fscan(reader, &k)
   sizes := []string{"S", "M", "L", "XL", "XXL"}
   results := make([]string, k)

   for i := 0; i < k; i++ {
       var desired string
       fmt.Fscan(reader, &desired)
       // map desired size to index
       var d int
       switch desired {
       case "S":
           d = 0
       case "M":
           d = 1
       case "L":
           d = 2
       case "XL":
           d = 3
       case "XXL":
           d = 4
       }
       // find closest available size
       assigned := -1
       for dist := 0; dist < len(counts); dist++ {
           if dist == 0 {
               if counts[d] > 0 {
                   assigned = d
                   break
               }
           } else {
               // prefer bigger size on ties
               if d+dist < len(counts) && counts[d+dist] > 0 {
                   assigned = d + dist
                   break
               }
               if d-dist >= 0 && counts[d-dist] > 0 {
                   assigned = d - dist
                   break
               }
           }
       }
       counts[assigned]--
       results[i] = sizes[assigned]
   }

   for _, sz := range results {
       fmt.Fprintln(writer, sz)
   }
}
