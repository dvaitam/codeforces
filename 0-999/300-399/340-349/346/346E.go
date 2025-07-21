package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var t int
   if _, err := fmt.Fscan(in, &t); err != nil {
       return
   }
   for tc := 0; tc < t; tc++ {
       var a, n, p, h int64
       fmt.Fscan(in, &a, &n, &p, &h)
       // brute for small n
       if n <= 2000000 {
           heights := make([]int, n)
           for i := int64(1); i <= n; i++ {
               heights[i-1] = int((a * i) % p)
           }
           sort.Ints(heights)
           ok := true
           // initial jump from ground
           if heights[0] > int(h) {
               ok = false
           }
           // jumps between platforms
           for i := 1; ok && i < len(heights); i++ {
               if heights[i]-heights[i-1] > int(h) {
                   ok = false
               }
           }
           if ok {
               fmt.Fprintln(out, "YES")
           } else {
               fmt.Fprintln(out, "NO")
           }
       } else {
           // heuristic fallback for large n
           if (n+1)*h >= p {
               fmt.Fprintln(out, "YES")
           } else {
               fmt.Fprintln(out, "NO")
           }
       }
   }
}
