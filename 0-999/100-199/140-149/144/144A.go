package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // find first occurrence of max and last occurrence of min
   maxVal := a[0]
   minVal := a[0]
   maxIdx := 0
   minIdx := 0
   for i, v := range a {
       if v > maxVal {
           maxVal = v
           maxIdx = i
       }
       if v <= minVal {
           minVal = v
           minIdx = i
       }
   }
   // moves: bring max to front, then min to end
   moves := maxIdx + (n - 1 - minIdx)
   // if max is before min in original, moving max shifts min right by 1, no adjustment
   // if max is after min, moving max left shifts min left by 1, reducing one swap
   if maxIdx > minIdx {
       moves--
   }
   fmt.Println(moves)
}
