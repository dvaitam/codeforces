package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   vis := make([]int, 8)
   for i := 0; i < n; i++ {
       var a int
       if _, err := fmt.Fscan(reader, &a); err != nil {
           return
       }
       if a >= 0 && a < len(vis) {
           vis[a]++
       }
   }
   // invalid numbers
   if vis[5] > 0 || vis[7] > 0 {
       fmt.Println(-1)
       return
   }
   // total triples must be n/3, each uses one '1'
   if vis[1] != n/3 {
       fmt.Println(-1)
       return
   }
   // check feasibility
   if vis[2]-vis[4] < 0 || vis[2]-vis[4] != vis[6]-vis[3] || vis[6] < vis[3] {
       fmt.Println(-1)
       return
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   // output triples: 1 2 4
   for i := 0; i < vis[4]; i++ {
       fmt.Fprintln(writer, "1 2 4")
   }
   // output triples: 1 2 6
   for i := 0; i < vis[6]-vis[3]; i++ {
       fmt.Fprintln(writer, "1 2 6")
   }
   // output triples: 1 3 6
   for i := 0; i < vis[3]; i++ {
       fmt.Fprintln(writer, "1 3 6")
   }
}
