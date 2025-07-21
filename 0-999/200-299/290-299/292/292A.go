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
   var curTime, queue, lastTime, maxQueue int64
   for i := 0; i < n; i++ {
       var t, c int64
       fmt.Fscan(reader, &t, &c)
       // process pending messages until time t
       dt := t - curTime
       var processed int64 = queue
       if processed > dt {
           processed = dt
       }
       if processed > 0 {
           lastTime = curTime + processed
           queue -= processed
       }
       curTime = t
       // receive new messages
       queue += c
       if queue > maxQueue {
           maxQueue = queue
       }
   }
   // process remaining messages after last task
   if queue > 0 {
       lastTime = curTime + queue
   }
   fmt.Fprintf(writer, "%d %d", lastTime, maxQueue)
}
