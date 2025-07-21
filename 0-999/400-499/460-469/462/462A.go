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
   board := make([]string, n)
   for i := 0; i < n; i++ {
       var line string
       fmt.Fscan(reader, &line)
       board[i] = line
   }
   // directions: up, down, left, right
   dirs := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           cnt := 0
           for _, d := range dirs {
               ni, nj := i+d[0], j+d[1]
               if ni >= 0 && ni < n && nj >= 0 && nj < n && board[ni][nj] == 'o' {
                   cnt++
               }
           }
           if cnt%2 != 0 {
               fmt.Println("NO")
               return
           }
       }
   }
   fmt.Println("YES")
}
