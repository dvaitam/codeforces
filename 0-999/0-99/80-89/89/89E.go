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

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n+1)
   l, r := n, 0
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
       if a[i] > 0 {
           if i < l {
               l = i
           }
           if i > r {
               r = i
           }
       }
   }
   p := -1
   for r > 0 {
       i := p + 2
       for i <= r && a[i] == 0 {
           writer.WriteString("AR")
           p++
           i++
       }
       value := 1
       d := 0
       first := true
       for i = p + 2; i <= r; i++ {
           if a[i] > 0 {
               value += 4
               if first {
                   d++
               }
           } else {
               value--
               first = false
           }
           if value <= 0 {
               tar := d
               for j := 0; j < tar; j++ {
                   writer.WriteString("AR")
               }
               writer.WriteByte('A')
               for j := 0; j < tar; j++ {
                   writer.WriteByte('L')
               }
               writer.WriteByte('A')
               for j := p + 2; j < p + d + 2; j++ {
                   if a[j] > 0 {
                       a[j]--
                   }
               }
               break
           }
       }
       if value > 0 {
           tar := r - p - 1
           for i = 0; i < tar; i++ {
               writer.WriteString("AR")
           }
           writer.WriteString("AL")
           p = r - 2
           if a[r] == 1 {
               for p+1 >= l && a[p+1] <= 1 {
                   p--
                   writer.WriteByte('L')
               }
           }
           for p+1 >= l && a[p+1] > 1 {
               p--
               writer.WriteByte('L')
           }
           writer.WriteByte('A')
           for i = p + 2; i <= r; i++ {
               if a[i] > 0 {
                   a[i]--
               }
           }
       }
       l = n
       r = 0
       for i := 1; i <= n; i++ {
           if a[i] > 0 {
               if i < l {
                   l = i
               }
               if i > r {
                   r = i
               }
           }
       }
   }
}
