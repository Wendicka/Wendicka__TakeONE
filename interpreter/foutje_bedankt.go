package wi

import (
  "os"
  "fmt"
)

// Constains the last error as a string.
// This value resets every time a new call to a Wendicka VM is made
// If no errors occurs, this is just an empty string.
var Werror string


// What to do when an error occurs in a Wendicka script.
// panic   = Will use the Go Panic routine to crash out the program. (Most dirty and undesirable method)
// crash   = Error message will show, and the program will end
// warn    = Error message will show, but your program will continue
// nothing = Will show nothing
// both "panic" and "crash" will end your program.
// both "warn" and "nothing" will (mostly) terminate the script execution and store the error in WError
var Werhand = "panic"

type thwerror struct {
    f func(er string)
}

var hwerror = map[string] thwerror{}

func init(){
  hwerror["panic"] = thwerror{
    func(er string) {
      panic ( "ERROR!\n"+er )
      os.Exit(1)
    },
  }
  hwerror["crash"] = thwerror{
    func (er string) {
      fmt.Println("ERROR!\n"+er)
      os.Exit(1)
    },
  }
  hwerror["warn"] = thwerror{
    func (er string) {
      fmt.Println("ERROR!\n"+er)
      Werror = er
    },
  }
  hwerror["nothing"] = thwerror{
    func (er string) {
      Werror = er
    },
  }

}

func wError(e string){
  hwerror[Werhand].f(e)
}
