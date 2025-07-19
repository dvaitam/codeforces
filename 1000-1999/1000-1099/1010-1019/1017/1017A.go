package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var a int
   if _, err := fmt.Fscan(reader, &a); err != nil {
       return
   }
   sums := make([]int, a)
   for i := 0; i < a; i++ {
       var g, h, k, j int
       fmt.Fscan(reader, &g, &h, &k, &j)
       sums[i] = g + h + k + j
   }
   target := sums[0]
   sort.Ints(sums)
   rank := 0
   for i := len(sums) - 1; i >= 0; i-- {
       rank++
       if sums[i] == target {
           fmt.Println(rank)
           return
       }
   }
}
