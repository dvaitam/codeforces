package main

import (
   "bufio"
   "os"
   "strconv"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   // fast integer reading
   readInt := func() int {
       var sign, x int = 1, 0
       var c byte
       // skip non-numeric characters
       for {
           b, err := reader.ReadByte()
           if err != nil {
               return 0
           }
           c = b
           if (c >= '0' && c <= '9') || c == '-' {
               break
           }
       }
       if c == '-' {
           sign = -1
           b, _ := reader.ReadByte()
           c = b
       }
       for c >= '0' && c <= '9' {
           x = x*10 + int(c-'0')
           b, err := reader.ReadByte()
           if err != nil {
               break
           }
           c = b
       }
       return x * sign
   }

   n := readInt()
   m := make([]int, n+2)
   for i := 1; i <= n; i++ {
       m[i] = readInt()
   }
   sumL := make([]int64, n+2)
   sumR := make([]int64, n+2)
   // compute prefix contributions: sum of minimums of subarrays ending at i
   lstack := make([]int, 0, n)
   for i := 1; i <= n; i++ {
       for len(lstack) > 0 && m[lstack[len(lstack)-1]] > m[i] {
           lstack = lstack[:len(lstack)-1]
       }
       if len(lstack) == 0 {
           sumL[i] = int64(i) * int64(m[i])
       } else {
           prev := lstack[len(lstack)-1]
           sumL[i] = sumL[prev] + int64(i-prev)*int64(m[i])
       }
       lstack = append(lstack, i)
   }
   // compute suffix contributions: sum of minimums of subarrays starting at i
   rstack := make([]int, 0, n)
   for i := n; i >= 1; i-- {
       for len(rstack) > 0 && m[rstack[len(rstack)-1]] > m[i] {
           rstack = rstack[:len(rstack)-1]
       }
       if len(rstack) == 0 {
           sumR[i] = int64(n-i+1) * int64(m[i])
       } else {
           nxt := rstack[len(rstack)-1]
           sumR[i] = sumR[nxt] + int64(nxt-i)*int64(m[i])
       }
       rstack = append(rstack, i)
   }
   // find best peak position
   var best int64 = -1
   peak := 1
   for i := 1; i <= n; i++ {
       total := sumL[i] + sumR[i] - int64(m[i])
       if total > best {
           best = total
           peak = i
       }
   }
   // rebuild solution
   a := make([]int, n+2)
   a[peak] = m[peak]
   for i := peak - 1; i >= 1; i-- {
       if m[i] > a[i+1] {
           a[i] = a[i+1]
       } else {
           a[i] = m[i]
       }
   }
   for i := peak + 1; i <= n; i++ {
       if m[i] > a[i-1] {
           a[i] = a[i-1]
       } else {
           a[i] = m[i]
       }
   }
   // output
   buf := make([]byte, 0, 20)
   for i := 1; i <= n; i++ {
       buf = strconv.AppendInt(buf[:0], int64(a[i]), 10)
       writer.Write(buf)
       if i < n {
           writer.WriteByte(' ')
       }
   }
   writer.WriteByte('\n')
}
