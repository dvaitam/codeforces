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
   if !scanner.Scan() {
       return
   }
   n, err := strconv.Atoi(scanner.Text())
   if err != nil {
       return
   }
   for i := 0; i < n; i++ {
       if !scanner.Scan() {
           break
       }
       s := scanner.Text()
       hasMiao := strings.HasPrefix(s, "miao.")
       hasLala := strings.HasSuffix(s, "lala.")
       switch {
       case hasMiao && hasLala:
           fmt.Println("OMG>.< I don't know!")
       case hasLala:
           fmt.Println("Freda's")
       case hasMiao:
           fmt.Println("Rainbow's")
       default:
           fmt.Println("OMG>.< I don't know!")
       }
   }
