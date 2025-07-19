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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }

   type event struct {
       orig string
       cur  string
   }
   events := make([]event, 0, t)
   curIndex := make(map[string]int)

   for i := 0; i < t; i++ {
       var a, b string
       fmt.Fscan(reader, &a, &b)

       if idx, ok := curIndex[a]; ok {
           events[idx].cur = b
           delete(curIndex, a)
           curIndex[b] = idx
       } else {
           events = append(events, event{orig: a, cur: b})
           curIndex[b] = len(events) - 1
       }
   }

   fmt.Fprintln(writer, len(events))
   for _, e := range events {
       fmt.Fprintln(writer, e.orig, e.cur)
   }
}
