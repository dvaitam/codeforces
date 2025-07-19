package main

import (
   "bufio"
   "fmt"
   "os"
)

func gcd(a, b int) int {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   dat := make([]int, n)
   odd := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &dat[i])
       if dat[i]&1 == 1 {
           odd++
       }
   }
   if n == 1 {
       fmt.Fprintln(writer, dat[0])
       ch := byte('a')
       for i := 0; i < dat[0]; i++ {
           writer.WriteByte(ch)
       }
       return
   }
   if odd >= 2 {
       fmt.Fprintln(writer, 0)
       for i := 0; i < n; i++ {
           ch := byte('a' + i)
           for j := 0; j < dat[i]; j++ {
               writer.WriteByte(ch)
           }
       }
       return
   }
   g := dat[0]
   for i := 1; i < n; i++ {
       g = gcd(g, dat[i])
   }
   fmt.Fprintln(writer, g)
   if odd == 0 {
       for tt := 0; tt < g/2; tt++ {
           // left half
           for i := 0; i < n; i++ {
               ch := byte('a' + i)
               times := dat[i] / g
               for j := 0; j < times; j++ {
                   writer.WriteByte(ch)
               }
           }
           // right half
           for i := n - 1; i >= 0; i-- {
               ch := byte('a' + i)
               times := dat[i] / g
               for j := 0; j < times; j++ {
                   writer.WriteByte(ch)
               }
           }
       }
   } else {
       for tt := 0; tt < g; tt++ {
           // left even parts
           for i := 0; i < n; i++ {
               if dat[i]&1 == 1 {
                   continue
               }
               ch := byte('a' + i)
               times := dat[i] / g / 2
               for j := 0; j < times; j++ {
                   writer.WriteByte(ch)
               }
           }
           // middle odd part
           for i := 0; i < n; i++ {
               if dat[i]&1 == 1 {
                   ch := byte('a' + i)
                   times := dat[i] / g
                   for j := 0; j < times; j++ {
                       writer.WriteByte(ch)
                   }
                   break
               }
           }
           // right even parts
           for i := n - 1; i >= 0; i-- {
               if dat[i]&1 == 1 {
                   continue
               }
               ch := byte('a' + i)
               times := dat[i] / g / 2
               for j := 0; j < times; j++ {
                   writer.WriteByte(ch)
               }
           }
       }
   }
}
