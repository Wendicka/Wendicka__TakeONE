package wi

import "fmt"

// Args
func (self *VM) Arg_Count() int{
  i:=0
  for _,ok:=self.identifiers.i[fmt.Sprintf("$__ARG%d",i)];ok;ok=ok {
    i++;
  }
  return i
}

func (self *VM) ID_ConvString(id *identifier) string{
  /*
  ak:=v //fmt.Sprintf("$__ARG%d",i)
  id,ok:=self.identifiers.i[ak]
  if !ok { return "nil" }
  */
  switch id.itype {
  case "nil":
    return "nil"
  case "string":
    return id.vstring
  case "int":
    return fmt.Sprint("%d",id.vint)
  case "bool","boolean":
    if id.vint<=0 { return "false"; } else { return "true";}
  case "float":
      return fmt.Sprint("%f",id.vfloat)
  case "chunk","function","procedure","api","table":
    return "<<"+id.itype+":"+id.vchunk+">>"
  default:
    return "Unknown type! ("+id.itype+")"
  }
}

func (self *VM) Arg_ConvString(i int) string{
  //ak:=fmt.Sprintf("$__ARG%d",i)
  ak:=VM.calls[WM.ccall].params.i[i]
  return self.ID_ConvString(ak)
}

func (self *VM) Var_ConvString(v string){
  id,ok:=self.identifiers.i[v]
  if !ok { return "nil" } else { return VM.ID_ConvString(id)}
}


// Minimal APIs
func wi_api_print(w *VM){
  // declare
  r:=[]string{}
  // form this dline
  l:=w.Arg_Count()
  for j:=1;j<=l;j++{
    r=append(r,w.Arg_ConvString(j))
  }
  // execute
  for i,v := range r {
    if i==0 { fmt.Print("\t"); }
    fmt.Print(v)
  }

  // Closure
  fmt.Println()

}

// register
func (w *API) minapi(){
  w.Register("PRINT",wi_api_print)
}
