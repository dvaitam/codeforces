package main

import (
   "bufio"
   "fmt"
   "os"
)

func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   rdr := bufio.NewReader(os.Stdin)
   wtr := bufio.NewWriter(os.Stdout)
   defer wtr.Flush()

   readInt := func() int64 {
       var x int64
       var c byte
       var err error
       // skip non-numeric
       for {
           c, err = rdr.ReadByte()
           if err != nil {
               return 0
           }
           if (c >= '0' && c <= '9') || c == '-' {
               break
           }
       }
       sign := int64(1)
       if c == '-' {
           sign = -1
           c, _ = rdr.ReadByte()
       }
       for c >= '0' && c <= '9' {
           x = x*10 + int64(c-'0')
           c, err = rdr.ReadByte()
           if err != nil {
               break
           }
       }
       return x * sign
   }

   n := readInt()
   a := make([]int64, n)
   for i := int64(0); i < n; i++ {
       a[i] = readInt()
   }

   total := make(map[int64]int64)
   prev := make(map[int64]int64)
   for i := int64(0); i < n; i++ {
       curr := make(map[int64]int64)
       // single element subarray
       curr[a[i]]++
       // extend previous subarrays
       for g, cnt := range prev {
           ng := gcd(g, a[i])
           curr[ng] += cnt
       }
       // accumulate totals
       for g, cnt := range curr {
           total[g] += cnt
       }
       prev = curr
   }

   q := readInt()
   for j := int64(0); j < q; j++ {
       x := readInt()
       if cnt, ok := total[x]; ok {
           fmt.Fprintln(wtr, cnt)
       } else {
           fmt.Fprintln(wtr, 0)
       }
   }
}
