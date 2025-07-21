package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   puzzles := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &puzzles[i])
   }
   sort.Ints(puzzles)
   // select n puzzles with minimal difference between max and min
   ans := puzzles[n-1] - puzzles[0]
   for i := 1; i+n-1 < m; i++ {
       diff := puzzles[i+n-1] - puzzles[i]
       if diff < ans {
           ans = diff
       }
   }
   fmt.Println(ans)
}
