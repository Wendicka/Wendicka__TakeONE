package wc

import "errors"
import "fmt"
import "trickyunits/qstr"
import "strings"
import "trickyunits/qint"
import "strconv"
//import "bytes"

// When set to 'true' the compiler could send out a few messages on screen
var CHAT = false

var debug = true


func dbgprint(a ...string){
  for i,p:=range a{
    if i!=0 { fmt.Print("\t");}
    fmt.Print(p)
  }
  fmt.Println()
}

func eh(f string,l int, er string) string{
  return fmt.Sprintf("%s, line #%d: %s",f,l+1,er)
}

func qe(f string,l int, er string) error{
  return errors.New(eh(f,l,er))
}

func appstring(ori []byte,s string) []byte {
  ret:=ori
  lb,_:= qint.Int32toBytes(int32(len(s)))
  sb  :=[]byte(s)
  ib  :=[2][]byte { lb,sb }
  for _,qb:= range ib {
    for _,b := range qb {
      ret=append(ret,b)
    }
  }
  return ret
}

func appint32(ori []byte,i int32) []byte {
    ret:=ori
    ai,_:=qint.Int32toBytes(i)
    for _,b := range ai{
      ret=append(ret,b)
    }
    return ret
}

func appint64(ori []byte,i int64) []byte {
    ret:=ori
    ai,_:=qint.Int64toBytes(i)
    for _,b := range ai{
      ret=append(ret,b)
    }
    return ret
}

func appfloat64(ori []byte,i float64) []byte {
    ret:=ori
    ai,_:=qint.FloatToBytes(i)
    for _,b := range ai{
      ret=append(ret,b)
    }
    return ret
}


func appparam(ori []byte,param string) ([]byte,error) {
  ret:=ori
  var err error
  if qstr.Prefixed(param,"\"") && qstr.Suffixed(param,"\""){
    ret=append(ret,0)
    ret=appstring(ret,param[1:len(param)-1])
  } else if qstr.Prefixed(param,"$") {
    ret=append(ret,1)
    ret=appstring(ret,param)
  } else if qstr.Prefixed(param,"0x") {
    var i int64
    i, err = strconv.ParseInt(param[2:], 16, 64)
    ret=append(ret,2)
    ret=appint64(ret,i)
  } else if qstr.Prefixed(param,"0") {
      var i int64
      i, err = strconv.ParseInt(param[1:], 8, 64)
      ret=append(ret,2)
      ret=appint64(ret,i)
  } else if qstr.Prefixed(param,"float:"){
    var i float64
    i, err = strconv.ParseFloat(param[6:] , 64)
    ret=append(ret,3)
    ret=appfloat64(ret,i)
  } else {
    var i int64
    i, err = strconv.ParseInt(param, 10, 64)
    ret=append(ret,2)
    ret=appint64(ret,i)

  }
  return ret,err
}

func chat(a string){
   if CHAT { fmt.Print(a)}
}

func Compile_Lines(source []string,f string) ([]byte,error){
    chunk:=""
    chunkpos:=0
    chunkcount:=0
    var chunklabels map[string]int
    myfile:=f; if myfile=="" {myfile="*SOURCE*"}
    ret := []byte{'W','B','C',8,0,0,0,0,0,0,0} // WBC = Wendicka Byte Code, 8 is a int64 indicating we only use 8bit bytes for instructions. This way, the system has been "prepared for future versions" in case they ever come with more bit.
    for lnum,dline := range (source)    {
       if CHAT {
         fmt.Printf("%3d%%",(lnum/len(source))*100)
         fmt.Print("\b\b\b\b")
       }
       tline:=qstr.MyTrim(dline)
       if tline!="" && !qstr.Prefixed(tline,"#") && !qstr.Prefixed(tline,";") && !qstr.Prefixed(tline,"--") && !qstr.Prefixed(tline,"//"){
         space:=strings.IndexAny(tline," ")
         instruction:=""
         args:=[]string{}
         if space==-1 {
           instruction=tline
         } else {
           instruction=tline[:space]
           inside:=false
           bs:=false
           dl:=0
           ara:=[]byte(tline[space+1:])
           cara:=[]byte{}
           for i:=0;i<len(ara);i++{
               if bs {
                 switch(ara[1]){
                 case 'b','B': cara = append(cara,8)
                 case 'n','N': cara = append(cara,10)
                 case 'r','R': cara = append(cara,13)
                 default:
                    cara = append(cara,ara[i])
                 }
                 bs=false
               } else if dl>0 {
                 a := 0
                 if ara[i]<48 || ara[i]>57 { return nil,errors.New(eh(f,lnum,"Invalid bytecode request!"))}
                 switch(dl){
                 case 3: a+=(int(ara[i])-48)*100
                 case 2: a+=(int(ara[i])-48)*10
                 case 1: a+=(int(ara[i])-48)*1
                 default: return nil,errors.New(eh(f,lnum,"Internal error. dl-messup! Please report!"))
                 }
                 dl--
                 if dl==0 {
                   if a>255 { return nil,errors.New(eh(f,lnum,fmt.Sprintf("Requested bytecode too high! (%d)",a)))}
                 }
                 cara = append(cara,byte(a))
               } else {
                 //fmt.Println(i,ara[i],string(cara)) // debug
                 switch ara[i]{
                    case ',':
                      if !inside{
                        args=append(args,string(cara))
                        cara=[]byte{}
                      }
                    case '"':
                      inside = !inside
                      cara=append(cara,ara[i])
                    case '#':
                      if inside { dl=3 } else {return nil,errors.New(eh(f,lnum,"# outside string"))}
                    case '\\':
                      if inside { bs=true } else {return nil,errors.New(eh(f,lnum,"\\ outside string"))}
                    default:
                      cara=append(cara,ara[i])
                      //args=append(args,string(ara))
                 }
               }
           }
           if inside { return nil,errors.New(eh(f,lnum,"Not properly ended string")) }
           if bs { return nil,errors.New(eh(f,lnum,"Backslash without follow up")) }
           args=append(args,string(cara))
           //dbgprint("Added last:",string(cara))
         }
         instruction=strings.ToUpper(instruction)
         if (chunk=="" && instruction!="CHUNK") { return nil,errors.New(eh(f,lnum,"Before any other action is taken a chunk must be created first"))}
         var err error
         switch(instruction){
            case "CHUNK":
              if len(args)<1 { return nil,qe(f,lnum,"CHUNK needs a name")}
              if chunkcount>0 { ret=append(ret,0xff); } // Make sure older chunks are always ended before starting a new one.
              ret = append(ret,0)
              ret = appstring(ret,args[0])
              chunk=args[0]
              chunkpos = len(ret)
              chunkcount++
              chunklabels=map[string]int{}
            case "LABEL":
              if len(args)<1 { return nil,qe(f,lnum,"LABEL needs a name")}
              //chunklabels[args[0]]=(len(ret)-chunkpos)
              ret = append(ret,254)
              ret = appstring(ret,args[0])
            case "JUMP","JMP":
              if len(args)<1 { return nil,qe(f,lnum,"JUMP needs a label name")}
              //chunklabels[args[0]]=(len(ret)-chunkpos)
              ret = append(ret,253)
              ret = appstring(ret,args[0])              
            case "MOV","MOVE":
              if len(args)<2 { return nil,qe(f,lnum,"MOVE needs 2 parameters")}
              ret = append(ret,1)
              ret = appstring(ret,args[0])
              ret,err = appparam(ret,args[1])
            case "INC","ADD":
              if len(args)<2 { return nil,qe(f,lnum,"INC needs 2 parameters")}
              ret = append(ret,2)
              ret = appstring(ret,args[0])
              ret,err = appparam(ret,args[1])
            case "DEC","SUBTRACT":
              if len(args)<2 { return nil,qe(f,lnum,"DEC needs 2 parameters")}
              ret = append(ret,3)
              ret = appstring(ret,args[0])
              ret,err = appparam(ret,args[1])
            case "SUM":
              if len(args)<3 { return nil,qe(f,lnum,"SUM needs 3 parameters")}
              ret = append(ret,4)
              ret = appstring(ret,args[0])
              ret,err = appparam(ret,args[1])
              if err==nil {
                ret,err = appparam(ret,args[2])
              }
            case "DIF":
              if len(args)<3 { return nil,qe(f,lnum,"DIF needs 3 parameters")}
              ret = append(ret,5)
              ret = appstring(ret,args[0])
              ret,err = appparam(ret,args[1])
              if err==nil {
                ret,err = appparam(ret,args[2])
              }
            case "MUL","MULTIPLY":
              if len(args)<3 { return nil,qe(f,lnum,"MULTIPLY needs 3 parameters")}
              ret = append(ret,6)
              ret = appstring(ret,args[0])
              ret,err = appparam(ret,args[1])
              if err==nil {
                ret,err = appparam(ret,args[2])
              }
            case "DIV","DIVIDE":
              if len(args)<3 { return nil,qe(f,lnum,"DIVIDE needs 3 parameters")}
              ret = append(ret,7)
              ret = appstring(ret,args[0])
              ret,err = appparam(ret,args[1])
              if err==nil {
                ret,err = appparam(ret,args[2])
              }
            case "CALL":
              if len(args)<1 { return nil,qe(f,lnum,"CALL needs 1 parameters")}
              ret = append(ret,10)
              ret = appstring(ret,args[0])
            case "EXIT","END":
              ret = append(ret,255)
            default:
              return nil,errors.New(eh(f,lnum,"Unknown instruction: "+instruction))
         }
       if err!=nil {qe(f,lnum,"Go-Error: "+err.Error())}
       }
    }
    ret = append(ret,255) // Make sure the code is properly ended!
    if chunkcount==1 {
      chat("Complete, 1 chunk\n")
    } else {
      chat(fmt.Sprintf("Complete, %d chunks\n",chunkcount))
    }
    return ret,nil
}

func Compile(source,f string) ([]byte,error){
  r,e:= Compile_Lines(strings.Split(source,"\n"),f)
  return r,e

}
