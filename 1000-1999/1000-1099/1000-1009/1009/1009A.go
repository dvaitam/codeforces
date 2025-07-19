package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   fmt.Fscan(reader, &n, &m)
   games := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &games[i])
   }
   bills := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &bills[i])
   }
   count := 0
   billIndex := 0
   for i := 0; i < n; i++ {
       if billIndex < m {
           if games[i] <= bills[billIndex] {
               count++
               billIndex++
           }
       } else {
           break
       }
   }
   fmt.Println(count)
}
