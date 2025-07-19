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

   // fast integer reader
   readInt := func() int {
       var x int
       var ch byte
       // skip non-digits
       for {
           b, err := reader.ReadByte()
           if err != nil {
               return x
           }
           ch = b
           if ch >= '0' && ch <= '9' {
               break
           }
       }
       // read digits
       for {
           if ch < '0' || ch > '9' {
               break
           }
           x = x*10 + int(ch-'0')
           b, err := reader.ReadByte()
           if err != nil {
               break
           }
           ch = b
       }
       return x
   }

   n := readInt()
   a := make([]int64, n+1)
   b := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       k := readInt()
       a[i] = a[i-1] + int64(k)
   }
   for i := 1; i <= n; i++ {
       k := readInt()
       b[i] = b[i-1] + int64(k)
   }

   // choose branch based on total sums
   if a[n] <= b[n] {
       diffMap := make(map[int64][2]int)
       diffMap[0] = [2]int{0, 0}
       l := 0
       for i := 1; i <= n; i++ {
           for l < n && b[l+1] <= a[i] {
               l++
           }
           diff := a[i] - b[l]
           if prev, ok := diffMap[diff]; ok {
               // output first array segment
               lenA := i - prev[0]
               fmt.Fprintln(writer, lenA)
               for j := prev[0] + 1; j <= i; j++ {
                   if j > prev[0]+1 {
                       writer.WriteByte(' ')
                   }
                   fmt.Fprint(writer, j)
               }
               writer.WriteByte('\n')
               // output second array segment
               lenB := l - prev[1]
               fmt.Fprintln(writer, lenB)
               for j := prev[1] + 1; j <= l; j++ {
                   if j > prev[1]+1 {
                       writer.WriteByte(' ')
                   }
                   fmt.Fprint(writer, j)
               }
               writer.WriteByte('\n')
               return
           }
           diffMap[diff] = [2]int{i, l}
       }
   } else {
       diffMap := make(map[int64][2]int)
       diffMap[0] = [2]int{0, 0}
       l := 0
       for i := 1; i <= n; i++ {
           for l < n && a[l+1] <= b[i] {
               l++
           }
           diff := b[i] - a[l]
           if prev, ok := diffMap[diff]; ok {
               // output first array segment (from A)
               lenA := l - prev[1]
               fmt.Fprintln(writer, lenA)
               for j := prev[1] + 1; j <= l; j++ {
                   if j > prev[1]+1 {
                       writer.WriteByte(' ')
                   }
                   fmt.Fprint(writer, j)
               }
               writer.WriteByte('\n')
               // output second array segment (from B)
               lenB := i - prev[0]
               fmt.Fprintln(writer, lenB)
               for j := prev[0] + 1; j <= i; j++ {
                   if j > prev[0]+1 {
                       writer.WriteByte(' ')
                   }
                   fmt.Fprint(writer, j)
               }
               writer.WriteByte('\n')
               return
           }
           diffMap[diff] = [2]int{i, l}
       }
   }
}
