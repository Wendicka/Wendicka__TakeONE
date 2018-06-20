package main
import "wendicka/compiler"
import "flag"
import "trickyunits/ansistring"
import "trickyunits/qstr"
import "trickyunits/qff"
import "os"
import "fmt"

const Red = ansistring.A_Red
const Yellow = ansistring.A_Yellow
const Cyan = ansistring.A_Cyan
const Magenta = ansistring.A_Magenta
const Green = ansistring.A_Green
const Blue = ansistring.A_Blue
const White = ansistring.A_White

const Blink  = ansistring.A_Blink
const Bright = ansistring.A_Bright

func crashout(e string){
   fmt.Println(ansistring.SCol("ERROR! ",Red,0)+ansistring.SCol(e,White,Bright))
   os.Exit(1)
 }

 func ecrash(e error){
     crashout(e.Error())
 }


func main(){
  var b []byte
  wc.CHAT=true
  fl_out:=flag.String("o","","Custom output file")
  flag.Parse()
  nonflags:=flag.Args()
	if len(nonflags)<1 {
  		fmt.Print(ansistring.SCol("Usage: ",Red,0),ansistring.SCol("wendicka build ",Yellow,0),ansistring.SCol("[ flags ] ",Magenta,ansistring.A_Dark),ansistring.SCol("<source> ",Cyan,0),"\n\n")
      flag.PrintDefaults()
      fmt.Println("\n\n(c) Jeroen P. Broks")
      os.Exit(0)
  }
  outfile:=*fl_out
  infile:=nonflags[0]
  if outfile=="" {
     outfile=qstr.StripExt(infile)+".wcc"
  }
  fmt.Print(ansistring.SCol("Processing: ",Yellow,0))
  fmt.Print(ansistring.SCol(infile,Cyan,0))
  fmt.Print(ansistring.SCol(" ... ",Magenta,0))
  s,e:=qff.EGetString(infile)
  if e!=nil { ecrash(e) }
  b,e = wc.Compile(s,infile)
  if e!=nil { ecrash(e) }
  fmt.Print(ansistring.SCol("Writing:    ",Yellow,0))
  fmt.Print(ansistring.SCol(outfile,Cyan,0))
  fmt.Print(ansistring.SCol(" ... ",Magenta,0))
  e = qff.WriteStringToFile(outfile,string(b))
  if e!=nil { ecrash(e) }
  fmt.Println(ansistring.SCol("Done ",Green,0))
}
