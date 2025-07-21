package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, k int64
   if _, err := fmt.Fscan(in, &n, &k); err != nil {
       return
   }
   // binary search minimal v such that sum_{i=0..} floor(v/k^i) >= n
   var low, high int64 = 1, n
   for low < high {
       mid := (low + high) / 2
       if total(mid, k) >= n {
           high = mid
       } else {
           low = mid + 1
       }
   }
   fmt.Println(low)
}

func total(v, k int64) int64 {
   var sum int64
   for v > 0 {
       sum += v
       v /= k
   }
   return sum
}
