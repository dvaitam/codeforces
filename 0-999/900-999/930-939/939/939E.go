package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   rdr := bufio.NewReader(os.Stdin)
   wtr := bufio.NewWriter(os.Stdout)
   defer wtr.Flush()

   readInt := func() int64 {
       var x int64
       var sign int64 = 1
       b, err := rdr.ReadByte()
       if err != nil {
           return 0
       }
       // skip non-numeric
       for (b < '0' || b > '9') && b != '-' {
           b, _ = rdr.ReadByte()
       }
       if b == '-' {
           sign = -1
           b, _ = rdr.ReadByte()
       }
       for b >= '0' && b <= '9' {
           x = x*10 + int64(b-'0')
           b, _ = rdr.ReadByte()
       }
       return x * sign
   }

   Q := int(readInt())
   // maximum number of insertions is Q
   a := make([]int64, Q+5)
   var idx, cnt int
   var sum int64
   var Max int64
   var aver float64 = 1e18
   for i := 0; i < Q; i++ {
       op := readInt()
       if op == 1 {
           val := readInt()
           idx++
           a[idx] = val
           Max = val
       } else {
           // query
           aver = float64(sum+a[idx]) / float64(cnt+1)
           for cnt < idx-1 {
               tmp := float64(sum + a[idx] + a[cnt+1]) / float64(cnt+2)
               if tmp < aver {
                   aver = tmp
                   sum += a[cnt+1]
                   cnt++
               } else {
                   break
               }
           }
           res := float64(Max) - aver
           fmt.Fprintf(wtr, "%.8f\n", res)
       }
   }
}
