package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var k int
   var q string
   fmt.Fscan(in, &k)
   fmt.Fscan(in, &q)
   seen := make([]bool, 26)
   positions := make([]int, 0, k)
   for i := 0; i < len(q) && len(positions) < k; i++ {
       c := q[i] - 'a'
       if !seen[c] {
           seen[c] = true
           positions = append(positions, i)
       }
   }
   if len(positions) < k {
       fmt.Println("NO")
       return
   }
   fmt.Println("YES")
   for i := 0; i < k; i++ {
       start := positions[i]
       if i+1 < k {
           end := positions[i+1]
           fmt.Println(q[start:end])
       } else {
           fmt.Println(q[start:])
       }
   }
}
