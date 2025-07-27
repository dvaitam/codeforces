package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

var (
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
)

func readInt() (int, error) {
   var x int
   _, err := fmt.Fscan(reader, &x)
   return x, err
}

func maxInt64(a, b int64) int64 {
   if a > b {
       return a
   }
   return b
}

func main() {
   defer writer.Flush()
   T, err := readInt()
   if err != nil {
       return
   }
   for tc := 0; tc < T; tc++ {
       n, _ := readInt()
       a := make([]int64, n)
       for i := 0; i < n; i++ {
           v, _ := readInt()
           a[i] = int64(v)
       }
       sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
       // candidate products
       // five largest
       p1 := a[n-1] * a[n-2] * a[n-3] * a[n-4] * a[n-5]
       // two smallest and three largest
       p2 := a[0] * a[1] * a[n-1] * a[n-2] * a[n-3]
       // four smallest and one largest
       p3 := a[0] * a[1] * a[2] * a[3] * a[n-1]
       res := maxInt64(p1, maxInt64(p2, p3))
       fmt.Fprintln(writer, res)
   }
}
