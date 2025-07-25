package main

import (
   "bufio"
   "fmt"
   "math/rand"
   "os"
   "time"
)

func gcd(a, b int) int {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   rw := bufio.NewReadWriter(bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout))
   defer rw.Flush()
   var n int
   // read size of the array
   fmt.Fscan(rw, &n)
   // find maximum value by binary searching using query type 2
   lo, hi := 0, 1000000000
   for lo < hi {
       mid := (lo + hi) / 2
       // query if any element > mid
       fmt.Fprintf(rw, "? 2 %d\n", mid)
       rw.Flush()
       var res int
       fmt.Fscan(rw, &res)
       if res == 1 {
           lo = mid + 1
       } else {
           hi = mid
       }
   }
   maxVal := lo
   // sample values to determine d
   rand.Seed(time.Now().UnixNano())
   samples := make(map[int]struct{})
   values := make([]int, 0, 30)
   // always include one sample of max if possible by searching equal
   // sample up to min(n, 30) distinct positions
   k := 30
   if n < k {
       k = n
   }
   for len(values) < k {
       idx := rand.Intn(n) + 1
       if _, ok := samples[idx]; ok {
           continue
       }
       samples[idx] = struct{}{}
       fmt.Fprintf(rw, "? 1 %d\n", idx)
       rw.Flush()
       var v int
       fmt.Fscan(rw, &v)
       values = append(values, v)
   }
   // compute gcd of differences
   d := 0
   for _, v := range values {
       diff := maxVal - v
       if diff < 0 {
           diff = -diff
       }
       if diff == 0 {
           continue
       }
       if d == 0 {
           d = diff
       } else {
           d = gcd(d, diff)
       }
   }
   // if all sampled equal or n==1, set d=0
   if d == 0 && n > 1 {
       // fallback: arithmetic progression with positive d, but we got zero diff
       // can assume d = maxVal - values[0]
       d = maxVal - values[0]
   }
   // smallest element
   x1 := maxVal - (n-1)*d
   // answer
   fmt.Fprintf(rw, "! %d %d\n", x1, d)
   rw.Flush()
}
