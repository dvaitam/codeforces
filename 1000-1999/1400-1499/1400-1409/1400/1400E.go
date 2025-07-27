package main

import (
   "bufio"
   "fmt"
   "os"
)

var a []int

// dfs returns minimum operations to delete all in a[l:r] given current height h
func dfs(l, r, h int) int {
   if l >= r {
       return 0
   }
   // cost of removing each individually
   cost1 := r - l
   // find minimum value and its index in a[l:r]
   minVal := a[l]
   minPos := l
   for i := l + 1; i < r; i++ {
       if a[i] < minVal {
           minVal = a[i]
           minPos = i
       }
   }
   // cost of horizontal strokes: raise from h to minVal
   cost2 := minVal - h
   // sum costs for segments split by minVal
   last := l
   for i := l; i < r; i++ {
       if a[i] == minVal {
           cost2 += dfs(last, i, minVal)
           last = i + 1
       }
   }
   // tail segment
   cost2 += dfs(last, r, minVal)
   if cost2 < cost1 {
       return cost2
   }
   return cost1
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a = make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   res := dfs(0, n, 0)
   fmt.Fprint(writer, res)
}
