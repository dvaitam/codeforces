package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
   "strconv"
)

func main() {
   scanner := bufio.NewScanner(os.Stdin)
   scanner.Split(bufio.ScanLines)

   // Read number of procedures
   scanner.Scan()
   n, _ := strconv.Atoi(strings.TrimSpace(scanner.Text()))
   type Proc struct { name string; types []string }
   procs := make([]Proc, 0, n)
   for i := 0; i < n; i++ {
       scanner.Scan()
       line := strings.TrimSpace(scanner.Text())
       lparen := strings.Index(line, "(")
       rparen := strings.LastIndex(line, ")")
       header := line[:lparen]
       headerFields := strings.Fields(header)
       pname := headerFields[len(headerFields)-1]
       typesStr := line[lparen+1:rparen]
       parts := strings.Split(typesStr, ",")
       ts := make([]string, 0, len(parts))
       for _, t := range parts {
           ts = append(ts, strings.TrimSpace(t))
       }
       procs = append(procs, Proc{pname, ts})
   }

   // Read variables
   scanner.Scan()
   m, _ := strconv.Atoi(strings.TrimSpace(scanner.Text()))
   varMap := make(map[string]string, m)
   for i := 0; i < m; i++ {
       scanner.Scan()
       fields := strings.Fields(scanner.Text())
       if len(fields) >= 2 {
           varMap[fields[1]] = fields[0]
       }
   }

   // Read calls
   scanner.Scan()
   k, _ := strconv.Atoi(strings.TrimSpace(scanner.Text()))
   var out strings.Builder
   for i := 0; i < k; i++ {
       scanner.Scan()
       line := strings.TrimSpace(scanner.Text())
       lparen := strings.Index(line, "(")
       rparen := strings.LastIndex(line, ")")
       cname := strings.TrimSpace(line[:lparen])
       argsStr := line[lparen+1:rparen]
       args := []string{}
       if len(argsStr) > 0 {
           parts := strings.Split(argsStr, ",")
           for _, a := range parts {
               a = strings.TrimSpace(a)
               if t, ok := varMap[a]; ok {
                   args = append(args, t)
               } else {
                   args = append(args, "")
               }
           }
       }
       count := 0
       for _, p := range procs {
           if p.name != cname || len(p.types) != len(args) {
               continue
           }
           ok := true
           for j, pt := range p.types {
               if pt != "T" && pt != args[j] {
                   ok = false
                   break
               }
           }
           if ok {
               count++
           }
       }
       out.WriteString(strconv.Itoa(count))
       out.WriteByte('\n')
   }
   fmt.Print(out.String())
}
