package main

import (
   "bufio"
   "fmt"
   "io"
   "os"
   "strings"
)

func main() {
   // Read entire input
   data, err := io.ReadAll(bufio.NewReader(os.Stdin))
   if err != nil {
       fmt.Fprintln(os.Stderr, err)
       os.Exit(1)
   }
   s := strings.TrimSpace(string(data))
   indent := 0
   // Iterate through the string, extract tags
   for i := 0; i < len(s); {
       if s[i] == '<' {
           // find end of tag
           j := i + 1
           for j < len(s) && s[j] != '>' {
               j++
           }
           if j >= len(s) {
               break
           }
           token := s[i : j+1]
           // Determine if closing tag
           if len(token) >= 2 && token[1] == '/' {
               indent--
               if indent < 0 {
                   indent = 0
               }
           }
           // Print with current indent
           fmt.Println(strings.Repeat(" ", indent*2) + token)
           // If opening tag, increase indent
           if len(token) >= 2 && token[1] != '/' {
               indent++
           }
           i = j + 1
       } else {
           // Skip any unexpected characters
           i++
       }
   }
}
