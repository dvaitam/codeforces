package main

import (
   "bufio"
   "fmt"
   "os"
)

const maxBit = 30

var cnt01 [maxBit + 1]int64
var cnt10 [maxBit + 1]int64

// dfs processes array a for bit down to 0, counting inversion contributions
func dfs(a []int, bit int) {
   if bit < 0 || len(a) < 2 {
       return
   }
   var v0, v1 []int
   v0 = make([]int, 0, len(a))
   v1 = make([]int, 0, len(a))
   var zeros, ones int64
   for _, v := range a {
       if (v>>bit)&1 == 1 {
           cnt01[bit] += zeros
           ones++
           v1 = append(v1, v)
       } else {
           cnt10[bit] += ones
           zeros++
           v0 = append(v0, v)
       }
   }
   dfs(v0, bit-1)
   dfs(v1, bit-1)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   dfs(a, maxBit)
   var x int64
   var inv int64
   for b := 0; b <= maxBit; b++ {
       if cnt10[b] <= cnt01[b] {
           inv += cnt10[b]
       } else {
           inv += cnt01[b]
           x |= 1 << b
       }
   }
   fmt.Fprintln(writer, inv, x)
}
