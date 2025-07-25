package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var a, b, c, d int
   if _, err := fmt.Fscan(reader, &a, &b, &c, &d); err != nil {
       return
   }
   p := []int{a, b, c}
   sort.Ints(p)
   ans := 0
   if p[1]-p[0] < d {
       ans += d - (p[1] - p[0])
   }
   if p[2]-p[1] < d {
       ans += d - (p[2] - p[1])
   }
   writer.WriteString(fmt.Sprintf("%d\n", ans))
}
