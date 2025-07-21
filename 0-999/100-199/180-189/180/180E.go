package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, k int
   if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
       return
   }
   colors := make([][]int, m+1)
   for i := 1; i <= n; i++ {
       var c int
       fmt.Fscan(in, &c)
       colors[c] = append(colors[c], i)
   }
   best := 1
   // for each color, use sliding window on its positions
   for _, pos := range colors {
       sz := len(pos)
       if sz <= best {
           continue
       }
       left := 0
       for right := 0; right < sz; right++ {
           // number of deletions needed to make pos[left..right] contiguous of this color
           // segment length = pos[right] - pos[left] + 1
           // count of this color = right - left + 1
           // deletions = segment length - count
           for left <= right && (pos[right]-pos[left]+1 - (right-left+1) > k) {
               left++
           }
           curr := right - left + 1
           if curr > best {
               best = curr
           }
       }
   }
   // output result
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   fmt.Fprintln(w, best)
}
