package main

import (
   "bufio"
   "fmt"
   "math/rand"
   "os"
   "time"
)

var (
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
)

func flush() {
   writer.Flush()
}

// ask queries whether the train is in [l, r]
func ask(l, r int64) bool {
   fmt.Fprintf(writer, "? %d %d\n", l, r)
   flush()
   var resp string
   fmt.Fscan(reader, &resp)
   return resp == "Yes" || resp == "YES" || resp == "Y" || resp == "yes"
}

func answer(x int64) {
   fmt.Fprintf(writer, "! %d\n", x)
   flush()
   os.Exit(0)
}

func main() {
   defer flush()
   rand.Seed(time.Now().UnixNano())
   var n, k int64
   fmt.Fscan(reader, &n, &k)
   l, r := int64(1), n
   const threshold = 50
   for i := 0; i < 4500; i++ {
       if r-l+1 <= threshold {
           // brute force with random probes
           pos := l + rand.Int63n(r-l+1)
           if ask(pos, pos) {
               answer(pos)
           }
       } else {
           mid := (l + r) / 2
           if ask(l, mid) {
               r = mid
           } else {
               l = mid + 1
           }
       }
       // expand possible range by k
       if l-k > 1 {
           l -= k
       } else {
           l = 1
       }
       if r+k < n {
           r += k
       } else {
           r = n
       }
   }
   // Fallback: linear scan
   for x := l; x <= r; x++ {
       if ask(x, x) {
           answer(x)
       }
   }
}
