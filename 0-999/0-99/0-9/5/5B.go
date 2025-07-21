package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   scanner := bufio.NewScanner(os.Stdin)
   var lines []string
   for scanner.Scan() {
       lines = append(lines, scanner.Text())
   }
   // find maximum line length
   maxLen := 0
   for _, line := range lines {
       if len(line) > maxLen {
           maxLen = len(line)
       }
   }
   // prepare border
   border := strings.Repeat("*", maxLen+2)
   fmt.Println(border)
   flip := true
   // center each line
   for _, line := range lines {
       d := maxLen - len(line)
       left := d / 2
       right := d - left
       if d%2 != 0 {
           if flip {
               left = d / 2
               right = d - left
           } else {
               right = d / 2
               left = d - right
           }
           flip = !flip
       }
       // print framed line
       fmt.Print("*")
       fmt.Print(strings.Repeat(" ", left))
       fmt.Print(line)
       fmt.Print(strings.Repeat(" ", right))
       fmt.Println("*")
   }
   fmt.Println(border)
}
