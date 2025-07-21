package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   // Read number of commands
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // current path segments, empty means root
   path := make([]string, 0)
   for i := 0; i < n; i++ {
       var cmd string
       if _, err := fmt.Fscan(reader, &cmd); err != nil {
           return
       }
       switch cmd {
       case "pwd":
           // print current path
           if len(path) == 0 {
               fmt.Println("/")
           } else {
               fmt.Print("/")
               fmt.Print(strings.Join(path, "/"))
               fmt.Println("/")
           }
       case "cd":
           var param string
           if _, err := fmt.Fscan(reader, &param); err != nil {
               return
           }
           // determine starting point
           if strings.HasPrefix(param, "/") {
               // absolute path
               path = path[:0]
           }
           parts := strings.Split(param, "/")
           for _, p := range parts {
               if p == "" {
                   continue
               }
               if p == ".." {
                   // go up one, guaranteed not to go above root
                   if len(path) > 0 {
                       path = path[:len(path)-1]
                   }
               } else {
                   // regular dir
                   path = append(path, p)
               }
           }
       }
   }
}
