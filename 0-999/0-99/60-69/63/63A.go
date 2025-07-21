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
   rats := make([]string, 0, n)
   womenChildren := make([]string, 0, n)
   men := make([]string, 0, n)
   var captain string

   for i := 0; i < n; i++ {
       var name, status string
       fmt.Fscan(reader, &name, &status)
       switch status {
       case "rat":
           rats = append(rats, name)
       case "woman", "child":
           womenChildren = append(womenChildren, name)
       case "man":
           men = append(men, name)
       case "captain":
           captain = name
       }
   }

   for _, v := range rats {
       fmt.Fprintln(writer, v)
   }
   for _, v := range womenChildren {
       fmt.Fprintln(writer, v)
   }
   for _, v := range men {
       fmt.Fprintln(writer, v)
   }
   if captain != "" {
       fmt.Fprintln(writer, captain)
   }
}
