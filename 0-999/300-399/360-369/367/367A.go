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

   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   n := len(s)
   // prefix sums for x, y, z
   px := make([]int, n+1)
   py := make([]int, n+1)
   pz := make([]int, n+1)
   for i, ch := range s {
       px[i+1] = px[i]
       py[i+1] = py[i]
       pz[i+1] = pz[i]
       switch ch {
       case 'x': px[i+1]++
       case 'y': py[i+1]++
       case 'z': pz[i+1]++
       }
   }
   var m int
   fmt.Fscan(reader, &m)
   for i := 0; i < m; i++ {
       var l, r int
       fmt.Fscan(reader, &l, &r)
       // 1-based to prefix
       length := r - l + 1
       if length < 3 {
           fmt.Fprintln(writer, "YES")
           continue
       }
       cx := px[r] - px[l-1]
       cy := py[r] - py[l-1]
       cz := pz[r] - pz[l-1]
       // check max-min <= 1
       // find max and min
       mx := cx
       if cy > mx {
           mx = cy
       }
       if cz > mx {
           mx = cz
       }
       mn := cx
       if cy < mn {
           mn = cy
       }
       if cz < mn {
           mn = cz
       }
       if mx-mn > 1 {
           fmt.Fprintln(writer, "NO")
       } else {
           fmt.Fprintln(writer, "YES")
       }
   }
}
