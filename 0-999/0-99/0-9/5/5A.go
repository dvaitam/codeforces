package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   scanner := bufio.NewScanner(os.Stdin)
   participants := make(map[string]struct{})
   var total int64

   for scanner.Scan() {
       line := scanner.Text()
       if line == "" {
           continue
       }
       switch line[0] {
       case '+':
           // Add command
           name := line[1:]
           participants[name] = struct{}{}
       case '-':
           // Remove command
           name := line[1:]
           delete(participants, name)
       default:
           // Send command: format sender:message
           parts := strings.SplitN(line, ":", 2)
           if len(parts) != 2 {
               // invalid, skip
               continue
           }
           msg := parts[1]
           // length in bytes equals len(msg) for ASCII
           l := int64(len(msg))
           // number of participants
           n := int64(len(participants))
           total += n * l
       }
   }
   // Output total traffic
   fmt.Println(total)
}
