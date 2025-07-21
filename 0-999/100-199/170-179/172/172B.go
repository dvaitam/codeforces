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

   var a, b, m, r0 int
   fmt.Fscan(reader, &a, &b, &m, &r0)
   visited := make([]int, m)
   for i := range visited {
       visited[i] = -1
   }
   cur := r0
   for i := 1; ; i++ {
       cur = (a*cur + b) % m
       if visited[cur] != -1 {
           fmt.Fprintln(writer, i-visited[cur])
           return
       }
       visited[cur] = i
   }
}
