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

   var n, m int64
   fmt.Fscan(reader, &n, &m)
   // maximum teams limited by total participants and smallest group
   total := (n + m) / 3
   ans := total
   if n < ans {
       ans = n
   }
   if m < ans {
       ans = m
   }
   fmt.Fprintln(writer, ans)
}
