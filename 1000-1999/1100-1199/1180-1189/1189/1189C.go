package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   s := make([]uint8, n)
   for i := 0; i < n; i++ {
       var x int
       fmt.Fscan(reader, &x)
       s[i] = uint8(x)
   }
   // compute max power
   maxk := bits.Len(uint(n)) - 1
   // dp_val[k][i]: value at segment of length 2^k starting at i
   // dp_carry[k][i]: total carries in that segment
   dpVal := make([][]uint8, maxk+1)
   dpCarry := make([][]int, maxk+1)
   dpVal[0] = make([]uint8, n)
   dpCarry[0] = make([]int, n)
   copy(dpVal[0], s)
   // build DP
   for k := 1; k <= maxk; k++ {
       length := 1 << k
       half := 1 << (k - 1)
       size := n - length + 1
       if size <= 0 {
           dpVal[k] = nil
           dpCarry[k] = nil
           continue
       }
       dpVal[k] = make([]uint8, size)
       dpCarry[k] = make([]int, size)
       prevVal := dpVal[k-1]
       prevCarry := dpCarry[k-1]
       for i := 0; i < size; i++ {
           l := prevVal[i]
           r := prevVal[i+half]
           sum := int(l) + int(r)
           c := 0
           if sum >= 10 {
               c = 1
               sum -= 10
           }
           dpVal[k][i] = uint8(sum)
           dpCarry[k][i] = prevCarry[i] + prevCarry[i+half] + c
       }
   }
   var q int
   fmt.Fscan(reader, &q)
   for qi := 0; qi < q; qi++ {
       var l, r int
       fmt.Fscan(reader, &l, &r)
       // convert to 0-based
       l--
       length := r - l
       // length+1 is power of two
       k := bits.TrailingZeros(uint(length + 1))
       res := dpCarry[k][l]
       fmt.Fprintln(writer, res)
   }
}
