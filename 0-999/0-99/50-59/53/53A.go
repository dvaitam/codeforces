package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
   "strings"
)

func main() {
   scanner := bufio.NewScanner(os.Stdin)
   // Read the input prefix string s
   if !scanner.Scan() {
       return
   }
   s := scanner.Text()
   // Read number of visited pages n
   if !scanner.Scan() {
       return
   }
   n, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
   if err != nil {
       return
   }
   var best string
   // Iterate over n pages
   for i := 0; i < n; i++ {
       if !scanner.Scan() {
           break
       }
       p := scanner.Text()
       if strings.HasPrefix(p, s) {
           if best == "" || p < best {
               best = p
           }
       }
   }
   // Output result: lexicographically smallest completion or original prefix
   if best != "" {
       fmt.Println(best)
   } else {
       fmt.Println(s)
   }
}
