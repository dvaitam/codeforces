package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
   "strings"
)

// clean removes signs '-', ';', '_' and converts letters to lowercase
func clean(s string) string {
   var b strings.Builder
   for _, r := range s {
       switch r {
       case '-', ';', '_':
           // skip sign
       default:
           if 'A' <= r && r <= 'Z' {
               b.WriteRune(r + ('a' - 'A'))
           } else if 'a' <= r && r <= 'z' {
               b.WriteRune(r)
           }
           // ignore other chars (none expected)
       }
   }
   return b.String()
}

func main() {
   scanner := bufio.NewScanner(os.Stdin)
   // read initial three strings
   orig := make([]string, 3)
   for i := 0; i < 3; i++ {
       if !scanner.Scan() {
           return
       }
       orig[i] = scanner.Text()
   }
   // clean initial strings
   t := make([]string, 3)
   for i, s := range orig {
       t[i] = clean(s)
   }
   // prepare all concatenation permutations
   permList := [][3]int{{0, 1, 2}, {0, 2, 1}, {1, 0, 2}, {1, 2, 0}, {2, 0, 1}, {2, 1, 0}}
   valid := make(map[string]struct{}, len(permList))
   for _, p := range permList {
       concat := t[p[0]] + t[p[1]] + t[p[2]]
       valid[concat] = struct{}{}
   }
   // read number of students
   if !scanner.Scan() {
       return
   }
   n, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
   if err != nil {
       return
   }
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   // process each answer
   for i := 0; i < n; i++ {
       if !scanner.Scan() {
           break
       }
       answer := clean(scanner.Text())
       if _, ok := valid[answer]; ok {
           fmt.Fprintln(out, "ACC")
       } else {
           fmt.Fprintln(out, "WA")
       }
   }
}
