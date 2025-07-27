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
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   arr := make([]int, 0, n)
   last := 0
   for i := 0; i < n; i++ {
       var a1, l1, r1, k1 int
       fmt.Fscan(reader, &a1, &l1, &r1, &k1)
       a := (a1 + last) % (n + 1)
       l := (l1 + last) % (i + 1)
       r := (r1 + last) % (i + 1)
       // convert to 1-based
       l++
       r++
       if l > r {
           l, r = r, l
       }
       k := (k1 + last) % (n + 1)
       // append
       arr = append(arr, a)
       // brute count
       ans := 0
       // for each subarray [x,y]
       for x := l - 1; x < r; x++ {
           for y := x; y < r; y++ {
               // compute mex of arr[x:y]
               // use map or slice
               seen := make(map[int]bool)
               for t := x; t <= y; t++ {
                   seen[arr[t]] = true
               }
               mex := 0
               for {
                   if !seen[mex] {
                       break
                   }
                   mex++
               }
               if mex == k {
                   ans++
               }
           }
       }
       fmt.Fprintln(writer, ans)
       last = ans
   }
}
