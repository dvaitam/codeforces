package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

type Type int

const (
   Err Type = -1
)

func (t Type) pointer() Type {
   if t == Err {
       return Err
   }
   return t + 1
}

func (t Type) deref() Type {
   if t == Err {
       return Err
   }
   if t > 0 {
       return t - 1
   }
   // &void => err
   return Err
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   types := make(map[string]Type)
   // initial types
   types["void"] = 0
   types["errtype"] = Err

   for i := 0; i < n; i++ {
       var op string
       if _, err := fmt.Fscan(reader, &op); err != nil {
           break
       }
       switch op {
       case "typedef":
           var a, b string
           fmt.Fscan(reader, &a, &b)
           t := evalType(a, types)
           types[b] = t
       case "typeof":
           var a string
           fmt.Fscan(reader, &a)
           t := evalType(a, types)
           if t == Err {
               fmt.Println("errtype")
           } else {
               // void with t pointers
               sb := strings.Builder{}
               sb.WriteString("void")
               for k := 0; k < int(t); k++ {
                   sb.WriteByte('*')
               }
               fmt.Println(sb.String())
           }
       }
   }
}

// evalType parses a type expression like &&name*** and returns its Type
func evalType(s string, types map[string]Type) Type {
   // count prefix '&'
   na := 0
   for na < len(s) && s[na] == '&' {
       na++
   }
   // count suffix '*'
   ns := 0
   for ns < len(s)-na && s[len(s)-1-ns] == '*' {
       ns++
   }
   name := s[na : len(s)-ns]
   t, ok := types[name]
   if !ok {
       t = Err
   }
   // apply pointer operations first (asterisks)
   for i := 0; i < ns; i++ {
       t = t.pointer()
   }
   // then dereference operations
   for i := 0; i < na; i++ {
       t = t.deref()
   }
   return t
}
