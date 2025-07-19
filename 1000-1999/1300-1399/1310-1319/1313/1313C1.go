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
   readInt := func() int64 {
       var sign int64 = 1
       var c byte
       // skip non-numbers
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
       var x int64
       for c >= '0' && c <= '9' {
           x = x*10 + int64(c-'0')
           b, err := reader.ReadByte()
           if err != nil {
               break
           }
           c = b
       }
       return x * sign
   }

   n := int(readInt())
   m := make([]int64, n+2)
   for i := 1; i <= n; i++ {
       m[i] = readInt()
   }
   suml := make([]int64, n+2)
   sumr := make([]int64, n+3)
   // monotonic stacks for indices
   lstack := make([]int, 0, n)
   for i := 1; i <= n; i++ {
       for len(lstack) > 0 && m[lstack[len(lstack)-1]] > m[i] {
           lstack = lstack[:len(lstack)-1]
       }
       if len(lstack) == 0 {
           suml[i] = int64(i) * m[i]
       } else {
           prev := lstack[len(lstack)-1]
           suml[i] = suml[prev] + int64(i-prev) * m[i]
       }
       lstack = append(lstack, i)
   }
   rstack := make([]int, 0, n)
   var maxx int64 = -1
   var pos int = 1
   for i := n; i >= 1; i-- {
       for len(rstack) > 0 && m[rstack[len(rstack)-1]] > m[i] {
           rstack = rstack[:len(rstack)-1]
       }
       if len(rstack) == 0 {
           sumr[i] = int64(n-i+1) * m[i]
       } else {
           nxt := rstack[len(rstack)-1]
           sumr[i] = sumr[nxt] + int64(nxt-i) * m[i]
       }
       rstack = append(rstack, i)
       // compute best
       cur := suml[i] + sumr[i] - m[i]
       if cur > maxx {
           maxx = cur
           pos = i
       }
   }
   // rebuild
   for i := pos - 1; i >= 1; i-- {
       if m[i] > m[i+1] {
           m[i] = m[i+1]
       }
   }
   for i := pos + 1; i <= n; i++ {
       if m[i] > m[i-1] {
           m[i] = m[i-1]
       }
   }
   // output result
   buf := make([]byte, 0, 20)
   for i := 1; i <= n; i++ {
       buf = strconv.AppendInt(buf[:0], m[i], 10)
       writer.Write(buf)
       if i < n {
           writer.WriteByte(' ')
       }
   }
   writer.WriteByte('\n')
}
