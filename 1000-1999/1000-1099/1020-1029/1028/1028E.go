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
   fmt.Fscan(reader, &n)
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   same := true
   for i := 0; i < n; i++ {
       if a[i] != a[0] {
           same = false
           break
       }
   }
   if same {
       if a[0] == 0 {
           writer.WriteString("YES\n")
           for i := 0; i < n; i++ {
               writer.WriteString("1")
               if i+1 < n {
                   writer.WriteByte(' ')
               } else {
                   writer.WriteByte('\n')
               }
           }
       } else {
           writer.WriteString("NO\n")
       }
       return
   }
   it := -1
   for i := 0; i < n; i++ {
       prev := (i + n - 1) % n
       if a[i] > 0 && a[prev] < a[i] {
           it = i
           break
       }
   }
   b := make([]int64, n)
   const mod = int64(1000000007)
   b[it] = a[it]
   for i := 1; i < n; i++ {
       cur := (it - i + n) % n
       nxt := (it - i + 1 + n) % n
       x := b[nxt]
       k := mod / x
       if k < 1 {
           k = 1
       }
       x = k * x
       b[cur] = x + a[cur]
   }
   writer.WriteString("YES\n")
   for i := 0; i < n; i++ {
       fmt.Fprintf(writer, "%d", b[i])
       if i+1 < n {
           writer.WriteByte(' ')
       } else {
           writer.WriteByte('\n')
       }
   }
}
