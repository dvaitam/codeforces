package main

import (
   "bufio"
   "fmt"
   "os"
)

func readInt(r *bufio.Reader) int64 {
   var c byte
   var err error
   // skip non-numeric
   for {
       c, err = r.ReadByte()
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
       c, _ = r.ReadByte()
   }
   var val int64
   for err == nil && c >= '0' && c <= '9' {
       val = val*10 + int64(c-'0')
       c, err = r.ReadByte()
   }
   return sign * val
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   n64 := readInt(reader)
   t := readInt(reader)
   n := int(n64)
   a := make([]int64, n+2)
   for i := 1; i <= n; i++ {
       a[i] = readInt(reader)
   }
   x := make([]int, n+2)
   for i := 1; i <= n; i++ {
       x[i] = int(readInt(reader))
   }
   minPossible := make([]int64, n+2)
   maxPossible := make([]int64, n+2)
   for i := 1; i <= n; i++ {
       minPossible[i] = a[i] + t
       maxPossible[i] = 3000000000000000000
   }
   cnt := make([]int, n+3)
   ok := true
   for i := 1; i <= n; i++ {
       if x[i] < i {
           ok = false
           break
       }
       cnt[i]++
       cnt[x[i]]--
       if x[i] != n {
           limit := a[x[i]+1] + t - 1
           if limit < maxPossible[x[i]] {
               maxPossible[x[i]] = limit
           }
       }
   }
   if !ok {
       fmt.Fprintln(writer, "No")
       return
   }
   sum := 0
   for i := 1; i <= n; i++ {
       sum += cnt[i]
       if sum > 0 {
           nextMin := a[i+1] + t
           if nextMin > minPossible[i] {
               minPossible[i] = nextMin
           }
       }
   }
   b := make([]int64, n+2)
   var prev int64
   for i := 1; i <= n; i++ {
       cur := prev + 1
       if minPossible[i] > cur {
           cur = minPossible[i]
       }
       if cur > maxPossible[i] {
           ok = false
           break
       }
       b[i] = cur
       prev = cur
   }
   if !ok {
       fmt.Fprintln(writer, "No")
       return
   }
   fmt.Fprintln(writer, "Yes")
   for i := 1; i <= n; i++ {
       fmt.Fprint(writer, b[i])
       if i < n {
           writer.WriteByte(' ')
       }
   }
   writer.WriteByte('\n')
}
