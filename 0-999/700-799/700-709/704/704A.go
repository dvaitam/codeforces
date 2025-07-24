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

   var n, q int
   fmt.Fscan(reader, &n, &q)
   // notifications: store app id for each notification
   notifications := make([]int, 0, q)
   // read status of notifications
   read := make([]bool, q)
   // per-app list of notification ids
   appQueues := make([][]int, n+1)
   // per-app pointer to processed in type2
   appPtr := make([]int, n+1)
   // global pointer for type3
   globalPtr := 0
   unread := 0

   for i := 0; i < q; i++ {
       var typ, x int
       fmt.Fscan(reader, &typ, &x)
       switch typ {
       case 1:
           // new notification for app x
           id := len(notifications)
           notifications = append(notifications, x)
           appQueues[x] = append(appQueues[x], id)
           // initially unread
           // read[id] is false by default
           unread++
       case 2:
           // read all notifications of app x
           queue := appQueues[x]
           ptr := appPtr[x]
           for ptr < len(queue) {
               id := queue[ptr]
               if !read[id] {
                   read[id] = true
                   unread--
               }
               ptr++
           }
           appPtr[x] = ptr
       case 3:
           // read first x notifications globally
           t := x
           for globalPtr < t {
               if !read[globalPtr] {
                   read[globalPtr] = true
                   unread--
               }
               globalPtr++
           }
       }
       // output unread count
       fmt.Fprintln(writer, unread)
   }
}
