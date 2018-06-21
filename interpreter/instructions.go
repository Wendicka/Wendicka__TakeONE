package wi

import(
  "trickyunits/qint"
  "errors"
  "fmt"
)

type winstruct struct {
  do_it func(w *VM,args [][]byte) bool
  needparam []string
}

var winstructs = map[byte] *winstruct {}

func igidentifier(w *VM,b[]byte) (*tIdentifier,error){
  i:=b[0]
  p:=b[1:]
  switch i{
    case 0:
      r:=tIdentifier{}
      r.itype="string"
      r.vstring=string(p)
      return &r,nil
    case 1:
      r:=w.identifiers.i
      return r[string(p)],nil
    case 2:
      var e error
      r:=tIdentifier{}
      r.itype="integer"
      r.vint,e=qint.BytesToInt64(p)
      //fmt.Println("Form integer",b,i,p,r.vint)
      return &r,e
    case 3:
      var e error
      r:=tIdentifier{}
      r.itype="float"
      r.vfloat,e=qint.BytesToFloat(p)
      return &r,e
    default:
      return nil,errors.New("Unknown Identifier Code. Either the code is compiled incorrectly, or the code is corrupt, or your version of Wendicka is too old for this particular code.")
  }
  return nil,errors.New("I'm completely at a loss. This error could only happen to either a bug or an outdated version of Wendicka.")
}

func init(){
  winstructs[1] = &winstruct{
    // MOV
    func (w *VM,args [][]byte) bool{
      v:=string(args[0])
      if args[0][0]!='$' {
        wError("Invalid MOV! "+v)
        return false
      }
      if len(args)<=1 { wError("No values received to appoint to: "+string(args[0])); return false}
      id,e:=igidentifier(w,args[1])
      if e!=nil { wError(e.Error()); return false;}
      switch v {
        case "$__ARG":
          w.nextcallparam.i = append(w.nextcallparam.i,id)
          chat("Argument ",id.itype,string(args[1])," added for next call")
          return true
        default:
          kut:=&w.identifiers.i
          (*kut)[v]=id
          chat("Value ",string(args[1]),"("+id.itype+")"," assigned to variable ",v,fmt.Sprintf("%d",id.vint))
          return true
        }
    },
    []string{"string","identifier"},
  }

  winstructs[2] = &winstruct{
    // INC
    func (w *VM,args [][]byte) bool{
      v:=string(args[0])
      if args[0][0]!='$' {
        wError("Invalid MOV! "+v)
        return false
      }
      if len(args)<=1 { wError("No values received to alter: "+string(args[0])); return false}
      vr,ok:=w.identifiers.i[string(args[0])]
      if !ok { wError("Non-existent variable "+string(args[0])+" cannot be altered"); return false}
      ch,e:=igidentifier(w,args[1])
      if e!=nil { wError(e.Error()); return false;}
      switch(vr.itype){
		case "integer":
			switch(ch.itype){
				case "integer":
					vr.vint += ch.vint
				case "float":
					vr.vint += int64(ch.vfloat)
				default:
					wError("I cannot alter a "+vr.itype+" with a "+ch.itype)
					return false
			}
		case "float":
			switch(ch.itype){
				case "integer":
					vr.vfloat += float64(ch.vint)
				case "float":
					vr.vfloat += ch.vfloat
				default:
					wError("I cannot alter a "+vr.itype+" with a "+ch.itype)
					return false
			}
		default:
			wError("Illegal identifier for INC "+vr.itype)
			return false
		}
      return true
	  },
    []string{"string","identifier"},
  }
  winstructs[3] = &winstruct{
    // DEC
    func (w *VM,args [][]byte) bool{
      v:=string(args[0])
      if args[0][0]!='$' {
        wError("Invalid MOV! "+v)
        return false
      }
      if len(args)<=1 { wError("No values received to alter: "+string(args[0])); return false}
      vr,ok:=w.identifiers.i[string(args[0])]
      if !ok { wError("Non-existent variable "+string(args[0])+" cannot be altered"); return false}
      ch,e:=igidentifier(w,args[1])
      if e!=nil { wError(e.Error()); return false;}
      switch(vr.itype){
		case "integer":
			switch(ch.itype){
				case "integer":
					vr.vint -= ch.vint
				case "float":
					vr.vint -= int64(ch.vfloat)
				default:
					wError("I cannot alter a "+vr.itype+" with a "+ch.itype)
					return false
			}
		case "float":
			switch(ch.itype){
				case "integer":
					vr.vfloat -= float64(ch.vint)
				case "float":
					vr.vfloat -= ch.vfloat
				default:
					wError("I cannot alter a "+vr.itype+" with a "+ch.itype)
					return false
			}
		default:
			wError("Illegal identifier for INC "+vr.itype)
			return false
		}
      return true
    },
    []string{"string","identifier"},
  }


  // CALL
  winstructs[10] = &winstruct{
    func (w *VM,args [][]byte) bool{
          ch:=string(args[0])
          if args[0][0]=='$' {
            id:=w.identifiers.i[ch]
            if id.itype=="chunk" || id.itype=="api" { ch=id.vstring } else {wError("Call identifiers must refer to apis or chunks"); return false}
          }
          afs:=w.vapi.fs
          if _,ok:=w.chunks[ch];ok {
            w.callChunk(ch)
          } else if _,aok:=(*afs)[ch];aok {
            s,e:=w.callAPI(ch)
            if !s {wError("API Call error: "+e)}
          }
      return true
    },
    []string{"string"},
  }
  // Check
  winstructs[250] = &winstruct{
	  func(w *VM, args[][]byte) bool{
		  id1,e1:=igidentifier(w,args[1])
		  id2,e2:=igidentifier(w,args[2])
		  d:=args[0][0]
		  if e1!=nil { wError("CHECK#1"+e1.Error()); return false; }
		  if e2!=nil { wError("CHECK#2"+e2.Error()); return false; }
		  switch d{
			  case 0,1:
				if id1.itype==id2.itype && id1.itype=="string" { 
					w.lastcompare=id1.vstring == id2.vstring
				} else if id1.itype==id2.itype && id1.itype=="integer" {
					w.lastcompare=id1.vint == id2.vint
				} else if id1.itype==id2.itype && id1.itype=="float" {
					w.lastcompare=id1.vfloat == id2.vfloat
				} else {
					fmt.Println("WARNING! Unsupported compare ("+id1.itype+"<=>"+id2.itype+")")
					w.lastcompare = false
				}
				if d==1 { w.lastcompare=!w.lastcompare }
			  case 3:
				if id1.itype==id2.itype && id1.itype=="integer" {
					w.lastcompare = id1.vint < id2.vint
					chat(fmt.Sprintf("%s %d < %s %d ",id1.itype,id1.vint,id2.itype,id2.vint))
				} else if id1.itype==id2.itype && id1.itype=="float" {
					w.lastcompare = id1.vfloat < id2.vfloat
				} else if id1.itype=="float" && id2.itype=="integer" {
					w.lastcompare = id1.vfloat < float64(id2.vint)
				} else if id1.itype=="float" && id2.itype=="integer" {
					w.lastcompare = float64(id1.vint) < id2.vfloat
				} else {
					fmt.Println("WARNING! Unsupported smaller compare ("+id1.itype+"<=>"+id2.itype+")")
					w.lastcompare = false
				}
			  case 4:
				if id1.itype==id2.itype && id1.itype=="integer" {
					w.lastcompare = id1.vint > id2.vint
				} else if id1.itype==id2.itype && id1.itype=="float" {
					w.lastcompare = id1.vfloat > id2.vfloat
				} else if id1.itype=="float" && id2.itype=="integer" {
					w.lastcompare = id1.vfloat > float64(id2.vint)
				} else if id1.itype=="float" && id2.itype=="integer" {
					w.lastcompare = float64(id1.vint) > id2.vfloat
				} else {
					fmt.Println("WARNING! Unsupported greater compare ("+id1.itype+"<=>"+id2.itype+")")
					w.lastcompare = false
				}
			  case 5:
				if id1.itype==id2.itype && id1.itype=="integer" {
					w.lastcompare = id1.vint >= id2.vint
				} else if id1.itype==id2.itype && id1.itype=="float" {
					w.lastcompare = id1.vfloat >= id2.vfloat
				} else if id1.itype=="float" && id2.itype=="integer" {
					w.lastcompare = id1.vfloat >= float64(id2.vint)
				} else if id1.itype=="float" && id2.itype=="integer" {
					w.lastcompare = float64(id1.vint) >= id2.vfloat
				} else {
					fmt.Println("WARNING! Unsupported greater equal compare ("+id1.itype+"<=>"+id2.itype+")")
					w.lastcompare = false
				}
			  case 6:
				if id1.itype==id2.itype && id1.itype=="integer" {
					w.lastcompare = id1.vint <= id2.vint
				} else if id1.itype==id2.itype && id1.itype=="float" {
					w.lastcompare = id1.vfloat <= id2.vfloat
				} else if id1.itype=="float" && id2.itype=="integer" {
					w.lastcompare = id1.vfloat <= float64(id2.vint)
				} else if id1.itype=="float" && id2.itype=="integer" {
					w.lastcompare = float64(id1.vint) <= id2.vfloat
				} else {
					fmt.Println("WARNING! Unsupported smaller equal compare ("+id1.itype+"<=>"+id2.itype+")")
					w.lastcompare = false
				}
			  default:
				wError(fmt.Sprintf("Unknown comparing code: %d!",d))
				return false
		  }
		  chat(fmt.Sprintf("d = %d #%d",d,len(args[0])))
		  if w.lastcompare { chat("CHECK = TRUE"); } else { chat("CHECK = FALSE"); }
		  return true
	  },
	  []string{"byte","identifier","identifier"},
  }
  
  // Jump
  winstructs[0xfd] = &winstruct{
	  func(w *VM,args [][]byte) bool{
		  l:=string(args[0])
		  cll:=len(w.calls)-1
		  cl:=w.calls[cll]
		  ch:=cl.achunk
		  pos,ok:=ch.labels[l]
		  if !ok { wError("Jump to non-existent label: "+l); return false }
		  cl.pos=pos-1 // -1 has to be taken in order for the pos++ instruction that always comes next!
		  return true
	  },
	  []string{"string"},  
  }
  
  // Conditional jumping
  winstructs[252] = &winstruct{ // false
	  func(w *VM,args [][]byte) bool{
		  if !w.lastcompare {
			  d:=winstructs[0xfd]
			  return d.do_it(w,args)
		  } else { return true }
	  },
	  []string{"string"},  
  }
  winstructs[251] = &winstruct{ // true
	  func(w *VM,args [][]byte) bool{
		  if w.lastcompare {
			  d:=winstructs[0xfd]
			  return d.do_it(w,args)
		  } else { return true }
	  },
	  []string{"string"},  
  }

  // END CALL
  winstructs[0xff] = &winstruct{
    func (w *VM,args [][]byte) bool{
		cll:=len(w.calls)-1
		cl:=w.calls[cll]
		cl.ended=true
		w.calls=w.calls[:cll]
      return true
    },
    []string{},

  }
}
