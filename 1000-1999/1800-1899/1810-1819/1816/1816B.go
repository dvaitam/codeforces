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

   var tc int
   if _, err := fmt.Fscan(reader, &tc); err != nil {
       return
   }
   for tc > 0 {
       tc--
       var n int
       fmt.Fscan(reader, &n)
       arr := make([][]int64, 2)
       arr[0] = make([]int64, n)
       arr[1] = make([]int64, n)
       // place largest numbers
       mx := int64(2 * n)
       tgl := 0
       arr[0][0] = mx
       mx--
       arr[1][n-1] = mx
       mx--
       // fill descending large numbers alternately
       for i := n - 2; i > 0; i-- {
           arr[tgl][i] = mx
           mx--
           tgl ^= 1
       }
       // fill ascending small numbers alternately
       mx = 1
       tgl = 1
       for i := 0; i < n; i++ {
           arr[tgl][i] = mx
           mx++
           tgl ^= 1
       }
       // output
       for i := 0; i < 2; i++ {
           for j := 0; j < n; j++ {
               if j > 0 {
                   writer.WriteByte(' ')
               }
               fmt.Fprint(writer, arr[i][j])
           }
           writer.WriteByte('\n')
       }
   }
}
