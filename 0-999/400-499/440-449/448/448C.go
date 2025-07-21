package main

import (
   "bufio"
   "fmt"
   "os"
)

var a []int

// solve returns the minimum strokes to paint segment [l, r) of the fence,
// given that it is already painted up to height h.
func solve(l, r, h int) int {
   if l >= r {
       return 0
   }
   // Option 1: paint each plank vertically
   vert := r - l
   // Find minimum height in [l, r)
   minVal := a[l]
   minIdx := l
   for i := l + 1; i < r; i++ {
       if a[i] < minVal {
           minVal = a[i]
           minIdx = i
       }
   }
   // Option 2: paint horizontally up to minVal, then recurse
   // strokes needed to raise from h to minVal
   hori := minVal - h
   // If horizontal strokes already exceed vertical, no need to proceed
   if hori >= vert {
       return vert
   }
   // paint left segment
   hori += solve(l, minIdx, minVal)
   // paint right segment
   hori += solve(minIdx+1, r, minVal)
   if hori < vert {
       return hori
   }
   return vert
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a = make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   result := solve(0, n, 0)
   fmt.Println(result)
}
