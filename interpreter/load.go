package wi

import(
  //"os"
  "fmt"
  "bytes"
  "errors"
  "trickyunits/qff"
  "trickyunits/qint"
)


func appint32(ori []byte,i int32) []byte {
    ret:=ori
    ai,_:=qint.Int32toBytes(i)
    for _,b := range ai{
      ret=append(ret,b)
    }
    return ret
}

// Convert bytes array into an actual Wendicka VM
func Bytes2VM(l []byte) (*VM,error){
  var e error
  ret:=NewVM(nil)
  bt:=bytes.NewReader(l)
  h:=qff.RawReadString(bt,3)
  readchunk:=""
  if h!="WBC" { return nil,errors.New("Given memory bank is NOT Wendicka Byte Code")}
  bits:=qff.ReadInt64(bt)
  if bits!=8 { return nil,errors.New(fmt.Sprintf("This version of the Wendicka interpreter only supports 8 bit instructions and this bytecode contains %d bit instructions",bits))}
  for bt.Len()>0 { //!qff.EOF(*bt) {
    ins:=qff.ReadByte(bt)
    if ins!=0 && readchunk=="" { return nil,errors.New(fmt.Sprintf("Chunkless code: %X",ins))}
    switch ins {
      case 0:
        namechunk:=qff.ReadString(bt)
        if _,exist:=ret.chunks[namechunk];exist { return nil,errors.New("Duplicate chunk: "+namechunk)}
        ret.chunks[namechunk]=tchunk{}
        readchunk=namechunk
        chat("Created chunk: ",namechunk)
      default:
        q,ok:=winstructs[ins]
        if !ok { return nil,errors.New(fmt.Sprintf("Unknown instruction %X",ins))}
        a:=[][]byte {}//*tIdentifier{}
        for i,partype:=range q.needparam {
          switch partype {
            case "string":
              s:= qff.ReadString(bt)
              a = append(a,[]byte(s))
            case "identifier":
              b:=[]byte{}
              t:=qff.ReadByte(bt)
              b=append(b,t)
              switch(t){
                case 0,1: // string + identifiername
                  l:=qff.ReadInt32(bt)
                  b=appint32(b,l)
                  for i:=int32(0);i<l;i++{ b=append(b,qff.ReadByte(bt));}
                case 2,3: // int64 + float64
                  for i:=0;i<8;i++{ b=append(b,qff.ReadByte(bt));}
                default:
                  return nil,errors.New("Invalid identifier type")
              }
            default:
              return nil,errors.New(fmt.Sprintf("Unknown paramter type for instruction %X, position %d: %s",ins,i,partype))
          }
        }
        cc:=tchunkins{}
        cc.param=a
        cc.ins=ins
        chat("Handled instruction:",fmt.Sprintf("%X",ins))
    }
  }
  return ret,e
}

// Load a file and convert that into a Wendicka VM
func File2VM(filename string) (*VM,error){
  b,e:=qff.EGetFile(filename)
  if e!=nil { return nil,e }
  var r *VM
  r,e=Bytes2VM(b)
  return r,e
}
