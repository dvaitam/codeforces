package main

import (
   "bufio"
   "fmt"
   "os"
)

func isPowerOfTwo(x int64) bool {
   return x > 0 && (x&(x-1)) == 0
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   var m int64
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   v := make([]int64, n)
   var vmax int64
   for i := 0; i < n; i++ {
       var ai int64
       fmt.Fscan(reader, &ai)
       v[i] = ai - 1
       if v[i] > vmax {
           vmax = v[i]
       }
   }
   // check chain condition: for all v[i], vmax/v[i] must be power of two
   for i := 0; i < n; i++ {
       if vmax%v[i] != 0 || !isPowerOfTwo(vmax/v[i]) {
           fmt.Fprint(writer, 0)
           return
       }
   }
   var ans int64
   // D = vmax * 2^k <= m-1
   var D int64 = vmax
   var twoK int64 = 1
   limitD := m - 1
   for D <= limitD {
       // for this D, x can be 1..min(2^k, m-D)
       maxX := twoK
       rem := m - D
       if rem < maxX {
           maxX = rem
       }
       if maxX <= 0 {
           break
       }
       ans += maxX
       // next k
       // prevent overflow
       if D > limitD/2 {
           break
       }
       D <<= 1
       twoK <<= 1
   }
   fmt.Fprint(writer, ans)
}
