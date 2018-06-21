package wi

import "fmt"

// When set "true" all kinds of needless debugshit will be displayed
var DebugCHAT = true


func chat(crap ...string){
  if !DebugCHAT{ return }
  for i,a := range crap {
    if i!=0 { fmt.Print("\t");}
    fmt.Print(a)
  }
  fmt.Println()
}
