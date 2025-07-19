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
   p := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &p[i])
   }

   for i := 1; i <= n; i++ {
       visited := make([]bool, n+1)
       x := p[i]
       visited[i] = true
       for !visited[x] {
           visited[x] = true
           x = p[x]
       }
       fmt.Fprint(writer, x, " ")
   }
   fmt.Fprintln(writer)
}
