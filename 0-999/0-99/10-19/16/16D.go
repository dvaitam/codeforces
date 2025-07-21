package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
   "strconv"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   // read n
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // consume end of line
   reader.ReadString('\n')
   days := 1
   prevTime := -1
   for i := 0; i < n; i++ {
       line, err := reader.ReadString('\n')
       if err != nil {
           // possibly EOF without newline
       }
       line = strings.TrimRight(line, "\r\n")
       // split time and message by first ": "
       idx := strings.Index(line, ": ")
       if idx < 0 {
           continue
       }
       timeStr := line[:idx]
       // timeStr is "hh:mm x.m."
       parts := strings.Split(timeStr, " ")
       if len(parts) < 2 {
           continue
       }
       hm := parts[0]
       ap := parts[1]
       hmParts := strings.Split(hm, ":")
       if len(hmParts) != 2 {
           continue
       }
       h, _ := strconv.Atoi(hmParts[0])
       m, _ := strconv.Atoi(hmParts[1])
       // convert to minutes since midnight
       var hour24 int
       if len(ap) > 0 && ap[0] == 'a' {
           if h == 12 {
               hour24 = 0
           } else {
               hour24 = h
           }
       } else {
           if h == 12 {
               hour24 = 12
           } else {
               hour24 = h + 12
           }
       }
       t := hour24*60 + m
       if i > 0 && t < prevTime {
           days++
       }
       prevTime = t
   }
   fmt.Println(days)
}
