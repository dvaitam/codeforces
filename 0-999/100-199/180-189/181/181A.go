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
   stars := [3][2]int{}
   cnt := 0
   for i := 1; i <= n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       for j := 1; j <= m; j++ {
           if s[j-1] == '*' {
               stars[cnt][0] = i
               stars[cnt][1] = j
               cnt++
           }
       }
   }
   // XOR rows and columns to find the missing vertex
   r := stars[0][0] ^ stars[1][0] ^ stars[2][0]
   c := stars[0][1] ^ stars[1][1] ^ stars[2][1]
   fmt.Println(r, c)
}
