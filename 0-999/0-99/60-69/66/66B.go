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
   heights := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &heights[i])
   }
   maxWatered := 0
   for i := 0; i < n; i++ {
       count := 1
       // flow to the left
       j := i
       for j > 0 && heights[j-1] <= heights[j] {
           count++
           j--
       }
       // flow to the right
       j = i
       for j < n-1 && heights[j+1] <= heights[j] {
           count++
           j++
       }
       if count > maxWatered {
           maxWatered = count
       }
   }
   fmt.Println(maxWatered)
}
