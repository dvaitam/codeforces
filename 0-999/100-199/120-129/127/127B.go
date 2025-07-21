package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   cnt := make([]int, 101)
   for i := 0; i < n; i++ {
       var x int
       fmt.Fscan(reader, &x)
       if x >= 1 && x <= 100 {
           cnt[x]++
       }
   }
   totalPairs := 0
   for _, c := range cnt {
       totalPairs += c / 2
   }
   frames := totalPairs / 2
   fmt.Println(frames)
}
