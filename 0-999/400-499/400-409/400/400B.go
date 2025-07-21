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

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   distances := make(map[int]struct{})
   for i := 0; i < n; i++ {
       var line string
       fmt.Fscan(reader, &line)
       posG, posS := -1, -1
       for j := 0; j < m; j++ {
           switch line[j] {
           case 'G':
               posG = j
           case 'S':
               posS = j
           }
       }
       // if G is to the right of S, impossible
       if posG > posS {
           fmt.Fprintln(writer, -1)
           return
       }
       d := posS - posG
       if d > 0 {
           distances[d] = struct{}{}
       }
   }
   // result is number of distinct positive distances
   fmt.Fprintln(writer, len(distances))
}
