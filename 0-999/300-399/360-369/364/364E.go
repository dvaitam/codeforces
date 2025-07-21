package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m, k int
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   grid := make([][]byte, n)
   for i := 0; i < n; i++ {
       line := make([]byte, m)
       if _, err := fmt.Fscan(reader, &line); err != nil {
           return
       }
       grid[i] = line
   }
   var res int64
   // arr holds column sums between top and bottom rows
   arr := make([]int, m)
   for top := 0; top < n; top++ {
       // reset arr
       for j := 0; j < m; j++ {
           arr[j] = 0
       }
       for bottom := top; bottom < n; bottom++ {
           for c := 0; c < m; c++ {
               if grid[bottom][c] == '1' {
                   arr[c]++
               }
           }
           // count subarrays of arr summing to k
           cnt := make(map[int]int)
           cnt[0] = 1
           sum := 0
           for c := 0; c < m; c++ {
               sum += arr[c]
               if sum >= k {
                   if v, ok := cnt[sum-k]; ok {
                       res += int64(v)
                   }
               }
               cnt[sum]++
           }
       }
   }
   fmt.Fprint(writer, res)
}
