package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   f := make([]int, n+1)
   c := make([]int, n+1)
   du := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fi, ci := 0, 0
       fmt.Fscan(reader, &fi, &ci)
       if fi == -1 {
           fi = 0
       }
       f[i] = fi
       c[i] = ci
       if ci == 0 {
           if fi >= 0 && fi <= n {
               du[fi]++
           }
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   printed := false
   for i := 1; i <= n; i++ {
       if du[i] == 0 && c[i] == 1 && f[i] != 0 {
           fmt.Fprintf(writer, "%d ", i)
           printed = true
       }
   }
   if !printed {
       writer.WriteString("-1")
   }
}
