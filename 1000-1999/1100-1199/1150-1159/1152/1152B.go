package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
)

// AC returns true if x+1 is a power of two (i.e., x of form 2^k-1)
func AC(x int64) bool {
   return bits.OnesCount64(uint64(x+1)) == 1
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int64
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   var ops []int64
   var t int
   for {
       if AC(n) {
           break
       }
       // operation: flip prefix
       t++
       // find lowest set bit index
       s := bits.TrailingZeros64(uint64(n))
       ops = append(ops, int64(s))
       // flip bits [0..s-1]
       n ^= (1<<uint(s)) - 1
       if AC(n) {
           break
       }
       // increment operation
       t++
       n++
   }
   // output
   fmt.Fprintln(writer, t)
   for i, v := range ops {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   if len(ops) > 0 {
       writer.WriteByte('\n')
   }
}
