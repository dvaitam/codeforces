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
   var T int
   if _, err := fmt.Fscan(reader, &T); err != nil {
       return
   }
   for t := 0; t < T; t++ {
       solve(reader, writer)
   }
}

func solve(reader *bufio.Reader, writer *bufio.Writer) {
   var n int
   fmt.Fscan(reader, &n)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   var ans [][2]int
   // align endpoints
   if a[0] != a[n-1] {
       ans = append(ans, [2]int{1, n})
       if (a[0]+a[n-1])%2 == 1 {
           a[n-1] = a[0]
       } else {
           a[0] = a[n-1]
       }
   }
   // adjust middle elements
   for i := 1; i < n-1; i++ {
       if a[i] == a[0] {
           continue
       }
       if (a[i]+a[0])%2 == 1 {
           ans = append(ans, [2]int{1, i + 1})
       } else {
           ans = append(ans, [2]int{i + 1, n})
       }
   }
   // output
   fmt.Fprintln(writer, len(ans))
   for _, p := range ans {
       fmt.Fprintln(writer, p[0], p[1])
   }
}
