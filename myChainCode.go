package main
import(
     "bufio"
     "bytes"
      "errors"
      "fmt"
     "strconv"
     "github.com/hyperledger/fabric/core/chaincode/shim"
     "runtime"
     "io"
    "os"
    "strings"
    "image"
   "image/gif"
   "image/png"
   "image/jpeg"
   "net/http"
   "time"    
"crypto/aes"
"crypto/cipher"
 "crypto/rand"
"encoding/json"

)
var recType=[]string{"ARTINV","USER","BID","AUCREQ","POSTRAN","OPENAUC","CLAUC","XFER","VERIFY","TRANS","CFER"}
var MyaucTables=[]string{"MyUserTable","MyUserCatTable","MyAssetHistoryTable","MyAssetTable","MyAssetCatTable","MyAssetAuctionTable","MyBidTable","MyCreditHistoryTable"}
type MyCreditLog struct{
     UserID string
     AuctionedBy string
     Amount string
     RecType string
     Desc string
     Date string
}
type MyAssetTransaction struct{
     AuctionID string
     RecType string
     AssetID string
     TransType string
    UserID string
    TransDate string
    HammerTime string
    HammerPrice string
   Details string
}
type MyBid struct{
   AuctionID string
   RecType string
   BideNo string
   AssetID string
   BuyerID string
   BidPrice string
   BideTime string
}
type MyAuctionRequest struct{
    AuctionID string
    RecType string
   AssetID string
    AuctionHouseID string
    SellerID string
     RequestDate string
     ReservePrice string
     BuyItNowPrice string
     Status string
     OpenDate string
    CloseDate string
}  
type MyAssetLog struct{
    AssetID string
    Status string
    RecType string
    AssetName string
    OwnerID string
   Date  string
   AuctionedBy string
}
type MyAssetObject struct{
    AssetID string
    RecType string
    OwnerID string
    AssetImageName string
    ImageType string
   AssetDate string 
  AssetKind string
    AssetPrice string
    AssetName string
    AES_Key []byte 
    AssetImage []byte
}
type MyUserObject struct{
     UserID string
     RecType string
     UserName    string
    UserType string 
   UserPhone string
     UserLevel string
     UserAmount string
     UserPassward string
}
type SimpleChainCode struct{
}
var gopath string
var ccPath string
func (t *SimpleChainCode) Init(stub shim.ChaincodeStubInterface,function string,args []string)([]byte,error){
      fmt.Println("Trade and Auction Application]Init")
      var err error
       for _,val:=range MyaucTables{
            err=stub.DeleteTable(val)
            if err!=nil{
                return nil,fmt.Errorf("init():delte of %s failed",val)
           }   
           err=InitLedger(stub,val)
           if err!=nil{
                return nil,fmt.Errorf("initledger of %s failed",val)
          }
       }
      err=stub.PutState("version",[]byte(strconv.Itoa(23)))
     if err!=nil{
          return nil,err
      }
    fmt.Println("init() initialization complite:",args)
     return []byte("init():initialization compliet"),nil
}
func InitLedger (stub shim.ChaincodeStubInterface,tableName string) error{
      nKeys:=GetNumberOfKeys(tableName)
      if nKeys<1{
           fmt.Println("Atleast 1 key must be provided\n")
           fmt.Println("Aucion_application:fail creating table",tableName)
          return errors.New("Auction_Application:Failed creating Table"+tableName)
          }
        var columnDefsForTbl []*shim.ColumnDefinition
        for i:=0;i<nKeys;i++{
             columnDef:=shim.ColumnDefinition{Name:"keyName"+strconv.Itoa(i),Type:shim.ColumnDefinition_STRING,Key:true}
             columnDefsForTbl=append(columnDefsForTbl,&columnDef)
        }
      columnLastTblDef:=shim.ColumnDefinition{Name:"Details",Type:shim.ColumnDefinition_BYTES,Key:false}
      columnDefsForTbl=append(columnDefsForTbl,&columnLastTblDef)
     err:=stub.CreateTable(tableName,columnDefsForTbl)
     if err!=nil{
         fmt.Println("Auction_application:fail create table",tableName)
        return errors.New("auction_appliction:fail create table"+tableName)
    }
    return err
}
func ChkReqType(args []string)bool{
       for _,rt:=range args{
               for _,val:=range recType{
                      if val==rt{
                            return true
                       }
                 }
           }
           return false
}
func (t *SimpleChainCode) Invoke(stub shim.ChaincodeStubInterface,function string ,args []string,)([]byte,error){
        var err error
        var buff []byte

      if ChkReqType(args)==true{
           InvokeRequest:=InvokeFunction(function)
           if InvokeRequest!=nil{
                  buff,err=InvokeRequest(stub,function,args)
            }  
      }else{
        fmt.Println("Invoke() Invalid recType:",args[1],"\n")
        return nil,errors.New("Invoke():Invalid recType:"+args[0])

     } 
  return buff,err
}
func main(){
 fmt.Println("hello")
   runtime.GOMAXPROCS(runtime.NumCPU())
   gopath=os.Getenv("GOPATH")
  if len(os.Args)==2 && strings.EqualFold(os.Args[1],"DEV"){
      fmt.Println("---------start in dev mode------")
      ccPath=fmt.Sprintf("%s/src/github.com/hyperledger/fabric/auction/art/myChainCode/",gopath)
    }else{
      fmt.Println("--------strat in net mode------")
       ccPath=fmt.Sprintf("%s/src/github.com/julia2804/auction/art/myChainCode/",gopath)
    }  
     err:=shim.Start(new(SimpleChainCode))
     if err!=nil{
          fmt.Printf("eror staring Simple chaincode:%s",err)
      }
}
func (t *SimpleChainCode) delete(stub shim.ChaincodeStubInterface,args []string)([]byte,error){
      if len(args)!=1{
          return nil,errors.New("incorretc numberof argument")
       }
       A:=args[0]
       err:=stub.DelState(A)
       if err!=nil{
           return nil,errors.New("fail to delete state")
        }
     return nil,nil
}
func CreateAssetObject(args []string)(MyAssetObject,error){
     var err error
     var myAsset MyAssetObject
     if len(args)!=8{
           fmt.Println("CreateAssetObject():incorrect number of argument")
           return myAsset,errors.New("CreateAssetobject():incorrect number of argument")
      }
      _,err=strconv.Atoi(args[0])
      if err!=nil{
          fmt.Println("CreateAssetObject():Id should be interger")
          return myAsset,errors.New("CreateOject():id should be integer")
      }

if err!=nil{
          fmt.Println("something wrong")
          panic(err) 
          }
           imagePath:=ccPath+args[2]
    if _,err:=os.Stat(imagePath);err==nil{
            fmt.Println(imagePath," exist")
    }else {
            fmt.Println("createAccetObject():cannot find or load image",imagePath)
            return myAsset,errors.New("createAssetOjet():ART Picture not found")
     }
     imagebytes,fileType:=imageToByteArray(imagePath)
       fmt.Println("image get succes")   
     AES_key,_:=GenAESKey()
   fmt.Println("genaeskey sucess") 
     AES_enc:=Encrypt(AES_key,imagebytes)
   fmt.Println("encrypt success") 
      myAsset=MyAssetObject{args[0],args[1],args[7],args[2],fileType,args[3],args[4],args[5],args[6],AES_key,AES_enc}
    fmt.Println("CreateAssetObject():Asset object created :ID#",myAsset.AssetID,"\n AES key:",myAsset.AES_Key)
    return myAsset,nil 
}
const(
      AESKeyLength=32
      NonceSize=24
)
func Encrypt(key []byte,ba []byte)[]byte{
     block,err:=aes.NewCipher(key)
     if err!=nil{
         panic(err)
     }
     ciphertext:=make([]byte,aes.BlockSize+len(ba))
     iv:=ciphertext[:aes.BlockSize]
     if _,err:=io.ReadFull(rand.Reader,iv);err!=nil{
        panic(err)
     }
     stream:=cipher.NewCFBEncrypter(block,iv)
     stream.XORKeyStream(ciphertext[aes.BlockSize:],ba)
     return ciphertext
} 
func Decrypt(key []byte,ciphertext []byte)[]byte{
      block,err:=aes.NewCipher(key)
      if err!=nil{
          panic(err)
       }
      if len(ciphertext)<aes.BlockSize{
          panic("text is too short")
       }
       iv:=ciphertext[:aes.BlockSize]
      ciphertext=ciphertext[aes.BlockSize:]
       stream:=cipher.NewCFBDecrypter(block,iv)
       stream.XORKeyStream(ciphertext,ciphertext)
      return ciphertext 
}
func GenAESKey()([]byte,error){
      fmt.Println("enter genskey") 
      return GetRandomBytes(AESKeyLength)
}
func GetRandomBytes(len int)([]byte,error){
   fmt.Println("random bytes success") 
     key:=make([]byte,len)
     _,err:=rand.Read(key)
     if err!=nil{
        return nil,err
     }
   fmt.Println("gen get random sucess ending") 
      return key,nil
}
func ValidateMember(stub shim.ChaincodeStubInterface,owner string)([]byte,error){
     args:=[]string{owner,"USER"}
     Avalbytes,err:=QueryLedger(stub,"MyUserTable",args)
     if err!=nil{
        fmt.Println("ValedateMember():fail -cannot find valid owner recodrd for it",owner)
        jsonResp:="{\"error\":\"fail to ger owner information"+owner+"\"}"
        return nil,errors.New(jsonResp)
}
     if Avalbytes==nil{
         fmt.Println("ValidateMember():fail-imcoplite information",owner)
         jsonResp:="{\"error\":\" fail-imcomplete information"+owner+"\"}"
         return nil,errors.New(jsonResp)
        }
       fmt.Println("validateMember():validateMamber success")
   return Avalbytes,nil
}
func UserToCreditLog(io MyUserObject) MyCreditLog{
       iLog:=MyCreditLog{}
       iLog.UserID=io.UserID
       iLog.AuctionedBy="DEFAULT"
       iLog.Amount=io.UserAmount 
        iLog.Desc="created"
       iLog.Date=time.Now().Format("2006-01-02 15:04:05")
       return iLog
}
func CreditLogtoJSON(credit MyCreditLog)([]byte,error){
       ajson,err:=json.Marshal(credit)
       if err!=nil{
             fmt.Println(err)
             return nil,err
        }
        return ajson,nil
}
func JSONtoCreditLog(ithis []byte)(MyCreditLog,error){
      credit:=MyCreditLog{}
     err:=json.Unmarshal(ithis,&credit)
     if err!=nil{
         fmt.Println("JSONtoCreditLog error:",err)
         return credit,err
     }
     return credit,err
}
func GetCreditLog(stub shim.ChaincodeStubInterface,function string,args []string)([]byte,error){
     if len(args)<1{
          fmt.Println("incorrect aument amont")
           return nil,errors.New("incorrect argument")
      }
      rows,err:=GetList(stub,"MyCreditHistoryTable",args)
      if err!=nil{
             return nil,fmt.Errorf("unmarshal son:%s",err)
      }
      nCol:=GetNumberOfKeys("MyCreditHistoryTable")
      tlist:=make([]MyCreditLog,len(rows))
      for i:=0;i<len(rows);i++{
           ts:=rows[i].Columns[nCol].GetBytes()
           il,err:=JSONtoCreditLog(ts)
           if err!=nil{
               fmt.Println("unmarshall error")
               return nil,fmt.Errorf("operation err:%s",err)
            }
            tlist[i]=il
       }
        jsonRows,_:=json.Marshal(tlist)
        return jsonRows,nil
}
func PostCreditLog(stub shim.ChaincodeStubInterface,user MyUserObject,amount string,ah string)([]byte,error){
      iLog:=UserToCreditLog(user)
      iLog.AuctionedBy=ah
    if ((strings.Compare(amount,"0"))!=0){
           iLog.Desc="ammented by "+ah 
           iLog.Amount=amount 
      } else {
           iLog.Desc="updated automatically"
           iLog.Amount="0" 
      } 
      buff,err:=CreditLogtoJSON(iLog)
      if err!=nil{
           fmt.Println("fail to create:",user.UserID)
           return nil,errors.New("failto create "+user.UserID)
       }else{
           keys:=[]string{iLog.UserID,iLog.AuctionedBy,time.Now().Format("2016-01-02 15:04:05")}
      err=UpdateLedger(stub,"MyCreditHistoryTable",keys,buff)
           if err!=nil{
             fmt.Println("write error")
              return buff,err
           }
      }
      return buff,nil
}
func PostAsset(stub shim.ChaincodeStubInterface,function string,args []string)([]byte,error){
    assetObject,err:=CreateAssetObject(args[0:])
   if err!=nil{
      fmt.Println("PostAsset():cannot create item object\n")
      return nil,err
    }
    ownerInfo,err:=ValidateMember(stub,assetObject.OwnerID)
    fmt.Println("owner information",ownerInfo,assetObject.OwnerID)
    if err!=nil{
        fmt.Println("postAsset():failed woner information not found",assetObject.OwnerID)
      }
   buff,err:=ARtoJSON(assetObject)
   if err!=nil{
       fmt.Println("PostAsset():fail cannot create object buff for write:",args[1])
      return nil,errors.New("PostAsset():fail cannot create objet buffer for write:"+ args[1])
     }else {
          keys:=[]string{args[0]}
          err=UpdateLedger(stub,"MyAssetTable",keys,buff)
          if err!=nil{
               fmt.Println("PostAsset():write error while insert\n")
               return buff,err
          }
      _,err=PostAssetLog(stub,assetObject,"INITIAL","DEFAULT")
     if err!=nil{
          fmt.Println("PostAssetLog():write error")
         return nil,err
       }
     fmt.Println("the args[5]:",args[5])
      keys=[]string{"2016",args[4],args[0]}
      err=UpdateLedger(stub,"MyAssetCatTable",keys,buff)
    if err!=nil{
         fmt.Println("PostAsset():write error")
         return buff,err
      }
  }  
     secret_key,_:=json.Marshal(assetObject.AES_Key)
    fmt.Println(string(secret_key))
    return secret_key,nil
}

func AssetToAssetLog(io MyAssetObject) MyAssetLog{
    iLog:=MyAssetLog{}
    iLog.AssetID=io.AssetID
    iLog.Status="INITIAL"
  iLog.AuctionedBy="DEFAULT" 
    iLog.RecType="ALOG"
    iLog.AssetName=io.AssetName
    iLog.OwnerID=io.OwnerID
    iLog.Date=time.Now().Format("2017-03-22 16:33:09")
    return iLog
} 

func PostAssetLog(stub shim.ChaincodeStubInterface,asset MyAssetObject,status string,ah string)([]byte,error){
     iLog:=AssetToAssetLog(asset)
     iLog.Status=status
     iLog.AuctionedBy=ah
     buff,err:=AssetLogtoJSON(iLog)
     if err!=nil{
        fmt.Println("PostAssetLog():failed cannotcreate object buffer "+asset.AssetID)
      return nil,errors.New("PostAssetLog():failed cannot create object"+asset.AssetID)
    }else {
         keys:=[]string{iLog.AssetID,iLog.Status,iLog.AuctionedBy,time.Now().Format("2017-03-22 16:33:09")}
        err=UpdateLedger(stub,"MyAssetHistoryTable",keys,buff)
       if err!=nil{
            fmt.Println("PostAssetLog():write error")
           return buff,err
       }
   }
 return buff,nil
}
func AssetLogtoJSON(asset MyAssetLog)([]byte,error){
     ajson,err:=json.Marshal(asset)
     if err!=nil{
        fmt.Println(err)
        return nil,err
     }
    return ajson,nil
} 
func JSONtoArgs(Avalbytes []byte)(map[string]interface{},error){
     var data map[string]interface{}
     if err:=json.Unmarshal(Avalbytes,&data);err!=nil{
             return nil,err
     }
     return data,nil
} 
func JSONtoAR(data []byte)(MyAssetObject,error){
     ar:=MyAssetObject{}
     err:=json.Unmarshal([]byte(data),&ar)
     if err!=nil{
        fmt.Println("Unmarshal failed:",err)
     }
     return ar,err
}
func ByteArrayToImage(imgByte []byte,imageFile string)error{
     img,_,_:=image.Decode(bytes.NewReader(imgByte))
     fmt.Println("processQueryResult byteArrayToImage:proceeding to create image")
     out,err:=os.Create(imageFile)
     if err!=nil{
        fmt.Println("byteArrayToImage():cannot crate image file ",err)
        return errors.New("byteArrayToImage():proced image file failed")
      }
      fmt.Println("processQueryType byteArrayToImage:proceding to encode image")
      filetype:=http.DetectContentType(imgByte)
      switch filetype{
      case "image/jpeg","image/jpg":
            var opt jpeg.Options
            opt.Quality=100
               err=jpeg.Encode(out,img,&opt)
      case "image/gif":
           var opt gif.Options
             opt.NumColors=256
           err=gif.Encode(out,img,&opt)
     case "image/png":
           err=png.Encode(out,img)
      default:
           err=errors.New("ohly pmng,jpg and gif supported")
      }
      if err!=nil{
         fmt.Println("ByteArrayToImage():cannot encode image file ",err)
         return errors.New("buteArrayToImage():cannot encode image file ")
      }
     fmt.Println("image filegenrated and saved to ",imageFile)
     return nil
}
func ARtoJSON(ar MyAssetObject)([]byte,error){
    fmt.Println("ar to json:",ar.AES_Key) 
      ajson,err:=json.Marshal(ar)
      if err!=nil{
         fmt.Println(err)
         return nil,err
      }
     return ajson,nil
}
func QueryFunction(fname string) func(stub shim.ChaincodeStubInterface,function string,args []string)([]byte,error){
fmt.Println("enter funtion")
fmt.Println("fanme:",fname)
    QueryFunc:=map[string]func(stub shim.ChaincodeStubInterface,function string,args []string)([]byte,error){
       "GetAsset": GetAsset,
       "GetUser": GetUser,
      "GetUserListByCat":GetUserListByCat,
      "GetAssetListByCat":GetAssetListByCat,
      "GetAssetLog":GetAssetLog,
    "ValidateItemOwnership":ValidateItemOwnership,  
   "GetCreditLog":GetCreditLog,  
     "ValidateUser":ValidateUser, 
    }
    return QueryFunc[fname]
}
func GetAssetLog(stub shim.ChaincodeStubInterface,function string,args []string)([]byte,error){
       if len(args)<1{
            fmt.Println("getAssetLog():incorect argument")
            return nil,errors.New("incorrect argumanet")
        }
       rows,err:=GetList(stub,"MyAssetHistoryTable",args)
       if err!=nil{
          return nil,fmt.Errorf("error marshal json:%s",err)
       }
       nCol:=GetNumberOfKeys("MyAssetHistoryTable")
       tlist:=make([]MyAssetLog,len(rows))
        for i:=0;i<len(rows);i++{
            ts:=rows[i].Columns[nCol].GetBytes()
            il,err:=JSONtoAssetLog(ts)
            if err!=nil{
                fmt.Println("unmarshalerror")
                return nil,fmt.Errorf("operation err:%s",err)
             }
            tlist[i]=il
        }
        jsonRows,_:=json.Marshal(tlist)
        return jsonRows,nil
} 
func ValidateItemOwnership(stub shim.ChaincodeStubInterface,function string,args []string)([]byte,error){
       var err error
       if len(args)<3{
            fmt.Println("item,owner,key needer")
            return nil,errors.New("request 3 arguent")
       }
       Avalbytes,err:=QueryLedger(stub,"MyAssetTable",[]string{args[0]})
      if err!=nil{
           fmt.Println("failed to query")
           jsonResp:="{\"error\" get q object data for "+args[0]+"\"}"
           return nil,errors.New(jsonResp)
       }
      if Avalbytes==nil{
           fmt.Println("fail imcoplete query")
            jsonResp:="{error imcomplete informateio for "+args[0]+"\"}"
            return nil,errors.New(jsonResp)
       }
      myItem,err:=JSONtoAR(Avalbytes)
     if err!=nil{
          fmt.Println("faile myitem")
          jsonResp:="{\"error\" get data for(item) "+args[0]+"\"}"
          return nil,errors.New(jsonResp)
        }
       myKey:=GetKeyValue(Avalbytes,"AES_Key")
     myName:=GetKeyValue(Avalbytes,"AssetName")
      myID:=GetKeyValue(Avalbytes,"AssetID") 
        fmt.Println("name string:",myName)
        fmt.Println("id string:",myID) 
        fmt.Println("key string:=",myKey)
       if myKey!=args[2]{
            fmt.Println("key not match",args[2],"-",myKey)
            jsonResp:="{\"error\" jey not match "+args[0]+"\"}"
            return nil,errors.New(jsonResp)
       }
        if myItem.OwnerID!=args[1]{
            fmt.Println("owner not march ",args[1])
            jsonResp:="{\"error\" owner id not marh"+args[0]+"\"}"
           return nil,errors.New(jsonResp) 
       } 
     fmt.Println("successful")
  return Avalbytes,nil
} 
func GetKeyValue(Avalbytes []byte,key string) string{
      var dat map[string]interface{}
      if err:=json.Unmarshal(Avalbytes,&dat);err!=nil{
              panic(err)
      }
     val:=dat[key].(string)
     return val
}
func TransferItem(stub shim.ChaincodeStubInterface,function string,args []string)([]byte,error){
      var err error
      if len(args)<5{
         fmt.Println("Translf():item,ownerid,key,new owner,moner")
          return nil,errors.New("dtransfer new 5 argument")
          }
        err=VerifyIfItemIsOnAuction(stub,args[0])
       if err!=nil{
            fmt.Println("faied",args[0])
            return nil,err
          }
        _,err=ValidateMember(stub,args[3])
        if err!=nil{
              fmt.Println("item not registe yet",args[3])
              return nil,err
        }
      ar,err:=ValidateItemOwnership(stub,"ValidateItemOwnership",args[:3])
      if err!=nil{
            fmt.Println("transfer r fail to authenticate:")
            return nil,err
        }
        myItem,err:=JSONtoAR(ar)
        if err!=nil{
           fmt.Println("faile create item from josn")
           return nil,err
        }
       CurrentAES_Key:=myItem.AES_Key
       image:=Decrypt(CurrentAES_Key,myItem.AssetImage)
      myItem.AES_Key,_=GenAESKey()
    myItem.AssetImage=Encrypt(myItem.AES_Key,image) 
      myItem.OwnerID=args[3]
     ar,err=ARtoJSON(myItem)
      keys:=[]string{myItem.AssetID,myItem.OwnerID}
      err=ReplaceLedgerEntry(stub,"MyAssetTable",keys,ar)
  if err!=nil{
    fmt.Println("transferasset failed to replsde")
    return nil,err
   }
   fmt.Println("transferasset sucess")
  keys=[]string{"2016",myItem.AssetKind,myItem.AssetID}
    err=ReplaceLedgerEntry(stub,"MyAssetCatTable",keys,ar)
    if err!=nil{
       fmt.Println("failed to replace asset at table")
       return nil,err
     }
     _,err=PostAssetLog(stub,myItem,"Transfer",args[1])
     if err!=nil{
         fmt.Println("write error post asset log")
        return nil,err
     }
    fmt.Println("myitem keys:",myItem.AES_Key) 
    fmt.Println("replace cat table success") 
    return myItem.AES_Key,nil
} 
func ReplaceLedgerEntry(stub shim.ChaincodeStubInterface,tableName string,keys []string,args []byte)error{
        nKey:=GetNumberOfKeys(tableName)
        if nKey<1{
             fmt.Println("at lest 1 key")
          }
          var columns []*shim.Column
         for i:=0;i<nKey;i++{
             col:=shim.Column{Value:&shim.Column_String_{String_:keys[i]}}
             columns=append(columns,&col)
         }
         lastCol:=shim.Column{Value:&shim.Column_Bytes{Bytes:[]byte(args)}}
         columns=append(columns,&lastCol)
         row:=shim.Row{columns}
          ok,err:=stub.ReplaceRow(tableName,row)
          if err!=nil{
                return fmt.Errorf("replace row into "+tableName+" able operation failed.%s",err)
          }
          if !ok{
                 return errors.New("replace row into "+tableName+" tabelf failed.Row with given key"+keys[0]+"alreay exist")
         }
        fmt.Println("replace row in "+tableName+" table success")
        return nil
} 
func VerifyIfItemIsOnAuction(stub shim.ChaincodeStubInterface,itemID string)error{
    return nil
}
func InvokeFunction(fname string) func(stub shim.ChaincodeStubInterface,function string,args []string)([]byte,error){
    InvokeFunc:=map[string]func(stub shim.ChaincodeStubInterface,function string,args []string)([]byte,error){
        "PostAsset":    PostAsset,
        "PostUser":     PostUser,
       "TransferCredit":TransferCredit, 
        "TransferItem":TransferItem,
        }
       return InvokeFunc[fname]
}
func ProcessQueryResult(stub shim.ChaincodeStubInterface,Avalbytes []byte,args []string)error{
      var dat map[string]interface{}
      if err:=json.Unmarshal(Avalbytes,&dat);err!=nil{
            panic(err)
        }
      var recType string
      recType=dat["RecType"].(string)
      switch recType{
      case "ARTINV":
          ar,err:=JSONtoAR(Avalbytes)
          if err!=nil{
                   fmt.Println("ProcessRequestType():Cannot creae assetObject \n")
                return err
            }
           image:=Decrypt(ar.AES_Key,ar.AssetImage)
           if err!=nil{
                fmt.Println("processRequestType():image decryption faied")
               return err
           }
          fmt.Println("ProcessRequestType():Image conversion sucessfull")
         err=ByteArrayToImage(image,ccPath+"copy."+ar.AssetImageName)
       if err!=nil{
            fmt.Println("ProcessRequestType():image conversion fail")
            return err
       }
      return err
     case "USER":
          ur,err:=JSONtoUser(Avalbytes)
          if err!=nil{
             return err
          }
         fmt.Println("ProcessRequestType():",ur)
        return err
     case "AUCREQ":
     case "OPENAUC":
     case "CLAUC":
         ar,err:=JSONtoAucReq(Avalbytes)
         if err!=nil{
              return err
         }
        fmt.Println("ProcessRequestType():",ar)
       return err
     case "POSTTRAN":
         atr,err:=JSONtoTran(Avalbytes)
         if err!=nil{
             return err
          }
          fmt.Println("PrcessRequestType():",atr)
     case "BID":
         bid,err:=JSONtoBid(Avalbytes)
         if err!=nil{
             return err
            }
          fmt.Println("processRequestType():",bid)
          return err
     case "DEFAULT":
          return nil
    case "XFER":
          return nil
    case "CFER":
          return nil 
      case "VERIFY":
           return nil
      default:
          return errors.New("unknown")
      }
     return nil
} 
func GetAsset(stub shim.ChaincodeStubInterface,function string,args []string)([]byte ,error){
     Avalbytes,err:=QueryLedger(stub,"MyAssetTable",args)
     if err!=nil{
        fmt.Println("gerAsser():fail to uery object")
        jsonResp:="{\"error\":\"fail to ger data for "+args[0]+"\"}"
        return nil,errors.New(jsonResp)
     }
     if Avalbytes==nil{
        fmt.Println("ger asset():incomplet query")
        jsonResp:="{\"err\":\"incomplete query" +args[0]+"\"}"
        return nil,errors.New(jsonResp)
     }
fmt.Println("get asset:response:success")
     assetObj,_:=JSONtoAR(Avalbytes)
     assetObj.AssetImage=[]byte{}
    fmt.Println("get asset:aeskdy:",assetObj.AES_Key) 
      Avalbytes,_=ARtoJSON(assetObj)
     return Avalbytes,nil
    }

func(t *SimpleChainCode) Query(stub shim.ChaincodeStubInterface,function string,args []string)([]byte,error){
    var err error
    var buff []byte
    fmt.Println("ID extracted and type= ",args[0])
    fmt.Println("Args supp;;",args)
    if len(args)<1{
         fmt.Println("at lest 1 arguments key")
         return nil, errors.New("query():expecting transaction type")       }
    QueryRequest:=QueryFunction(function)
    if QueryRequest!=nil{
       buff,err=QueryRequest(stub,function,args)
    }else {
        fmt.Println("query() invalid function call:",function)
        return nil,errors.New("Query():invalid functio acll:" +function)    }
if err!=nil{
    fmt.Println("query() object ot found:",args[0])
    return nil,errors.New("not found:" +args[0])
     }
    return buff,err
   }
func UpdateUserObject(stub shim.ChaincodeStubInterface,ar []byte,hammerUser string,amount string)(string,error){
      var err error
      myUser,err:=JSONtoUser(ar)
      if err!=nil{
           fmt.Println("fail to create")
          return "wrong",err
      }
     number,error:=strconv.Atoi(amount)
    if error!=nil{
          fmt.Println("tarandform failed")
     }
    amo,error:=strconv.Atoi(myUser.UserAmount)
   if error!=nil{
         fmt.Println("trander amount of user fialed")
    }
    amo=amo+number
    myUser.UserAmount=strconv.Itoa(amo)
    ar,err=UsertoJSON(myUser)
    keys:=[]string{myUser.UserID}
    err=ReplaceLedgerEntry(stub,"MyUserTable",keys,ar)
    if err!=nil{
        fmt.Println("fail to replace ledger")
        return "",err
     }
    fmt.Println("repleace user table succesfull")
    keys=[]string{"2016",myUser.UserType,myUser.UserID}
    err=ReplaceLedgerEntry(stub,"MyUserCatTable",keys,ar)
    if err!=nil{
          fmt.Println("replace cat table failed")
          return "",err
    }
     fmt.Println("succesful")
    return myUser.UserID,nil
} 
func TransferCredit(stub shim.ChaincodeStubInterface,function string,args[]string)([]byte,error){
       var err error
      if len(args)<5{
             fmt.Println("transferItem():argument wrong")
             return nil,errors.New("argument wrong")
        }
   ar,err:=ValidateMember(stub,args[0])
   if err!=nil{
      fmt.Println("gail valide", args[0])
      return nil,err
   }
   _,err2:=ValidateUser(stub,"ValidateUser",args)
   if err2!=nil{
       fmt.Println("gail authentic", args[2])
      return nil,err2
   }
  myUser,err:=JSONtoUser(ar)
   if err!=nil{
      fmt.Println("faile tao marshall")
      return nil,err
   }
  str:=myUser.UserAmount
 amo,_:=strconv.Atoi(str)
 count,_:=strconv.Atoi(args[1]) 
  amo=amo+count 
  string_amount:=strconv.Itoa(amo)
   myUser.UserAmount=string_amount 
   ar,err=UsertoJSON(myUser)
   keys:=[]string{myUser.UserID} 
   err=ReplaceLedgerEntry(stub,"MyUserTable",keys,ar)
    if err!=nil{
      fmt.Println("faile to replace user table")
      return nil,err
   }
   fmt.Println("success")
   keys=[]string{"2016",myUser.UserType,myUser.UserID}
   err=ReplaceLedgerEntry(stub,"MyUserCatTable",keys,ar)
  if err!=nil{      fmt.Println("faile to replace user cat table")
      return nil,err
   }
  _,err=PostCreditLog(stub,myUser,args[1],args[2])
   if err!=nil{ 
     fmt.Println("faile to replace user table")
      return nil,err
   }
  _,err=ValidateLevel(stub,"ValidateLevel",keys)
   if err!=nil{
     fmt.Println("faile to replace user table level")
   }
   fmt.Println("repleae success")
   return ar,nil
}
func ValidateLevel(stub shim.ChaincodeStubInterface,function string,args []string)([]byte,error){
{
  var err error
  arg:=[]string{args[2]}
   Avalbytes,err:=QueryLedger(stub,"MyUserTable",arg)
    if err!=nil{
           fmt.Println("fail to quey")
           jsonResp:="{\"error\":get data for"+args[2]+"\"}"
           return nil,errors.New(jsonResp)
    }
   if Avalbytes==nil{
        fmt.Println("incomplete query ojedt")
        jsonResp:="{\"error\":\"get data avalbtes err for "+args[2]+"\"}"
        return nil,errors.New(jsonResp)
   }
   myUser,err:=JSONtoUser(Avalbytes)
 if err!=nil{
           fmt.Println("fail to marshal")
           jsonResp:="{\"error\":\"get marshal for"+args[2]+"\"}"
           return nil,errors.New(jsonResp)
    }
   leve,_:=strconv.Atoi(myUser.UserLevel) 
    fmt.Println("amount:",myUser.UserAmount) 
    amo,_:=strconv.Atoi(myUser.UserAmount)
  fmt.Println("amount number",amo) 
     if (amo>1000 && leve==1 && amo<=3000){
           myUser.UserLevel="2"
            Avalbytes,err=UsertoJSON(myUser) 
            keys:=[]string{myUser.UserID,myUser.UserLevel} 
            err=ReplaceLedgerEntry(stub,"MyUserTable",keys,Avalbytes)  
            if err!=nil{
                fmt.Println("update user level:failed")
             }
            keys=[]string{"2016",myUser.UserType,myUser.UserID}
            err=ReplaceLedgerEntry(stub,"MyUserCatTable",keys,Avalbytes) 
            if err!=nil{
                fmt.Println("update user level:failed")
                return nil,err  
            }
            _,err=PostCreditLog(stub,myUser,"0",args[2])
             if err!=nil{
                fmt.Println("post creditlog error")
                 return nil,err
             }  
            fmt.Println("success update user level")
      } else if(amo>3000 && (leve==1 || leve==2)){
              myUser.UserLevel="3"
             Avalbytes,err=UsertoJSON(myUser)
              keys:=[]string{myUser.UserID,myUser.UserLevel} 
            err=ReplaceLedgerEntry(stub,"MyUserTable",keys,Avalbytes)
            if err!=nil{
                fmt.Println("update user level:failed")
             }
            keys=[]string{"2016",myUser.UserType,myUser.UserID}
            err=ReplaceLedgerEntry(stub,"MyUserCatTable",keys,Avalbytes)
            if err!=nil{
                fmt.Println("update user cat level:failed",err)
             }
            _,err=PostCreditLog(stub,myUser,"0",args[2])
             if err!=nil{
                fmt.Println("post creditlog error")
                 return nil,err
             } 
             fmt.Println("success update user level")
        }
     return Avalbytes,nil
}
}
func ValidateUser(stub shim.ChaincodeStubInterface,function string,args []string)([]byte,error){
     var err error
    arg:=[]string{args[2]} 
     Avalbytes,err:=QueryLedger(stub,"MyUserTable",arg)
    if err!=nil{
           fmt.Println("fail to quey")
           jsonResp:="{\"error\":get data for"+args[2]+"\"}"
           return nil,errors.New(jsonResp)
    }
   if Avalbytes==nil{
        fmt.Println("incomplete query ojedt")
        jsonResp:="{\"error\":\"get data avalbtes err for "+args[2]+"\"}"
        return nil,errors.New(jsonResp)
   }
   myUser,err:=JSONtoUser(Avalbytes)
 if err!=nil{
           fmt.Println("fail to marshal")
           jsonResp:="{\"error\":\"get marshal for"+args[2]+"\"}"
           return nil,errors.New(jsonResp)
    }
   if strings.Compare(myUser.UserType,"商家")!=0{
         fmt.Println("this user has no right tranfer redit")
         jsonResp:="{\"error\":\"not right to transfer"+args[2]+"\"}"
         return nil,errors.New(jsonResp)
   }
   fmt.Println("valide success")
    return Avalbytes,nil
}
func CreateUserObject(args []string)(MyUserObject,error){
      var err error
      var aUser MyUserObject
      if len(args)!=8{
        fmt.Println("CreateUserObject():incorrect argument number,expecting 8")
      }
      _,err=strconv.Atoi(args[0])
     if err!=nil{
         return aUser,errors.New("CreateUserObject():Incorrect number of user id")
    }
    aUser=MyUserObject{args[0],args[1],args[2],args[3],args[4],args[5],args[6],args[7]}
    fmt.Println("CreateUserObject():User Object:",aUser)
  return aUser,nil
} 
func UsertoJSON(user MyUserObject)([]byte,error){
   fmt.Println("user.phone:",user.UserPhone) 
    ajson,err:=json.Marshal(user)
    if err!=nil{
      fmt.Println("UserJSON error:",err)
      return nil,err
      }
     fmt.Println("UsertoJSON created:",ajson)
    record,_:=JSONtoUser(ajson) 
   fmt.Println("record.level,userjson:",record.UserPhone) 
     return ajson,nil
} 
func GetNumberOfKeys(tname string )int{
     TableMap:=map[string]int{
            "MyUserTable":   1,
            "MyUserCatTable": 3,
            "MyAssetCatTable": 3,
            "MyAssetTable":    1,
            "MyAssetHistoryTable":4,
           "MyBidTable":2,
           "MyTransTable":2,
          "MyAssetAuctionTable":1, 
          "MyCreditHistoryTable":3, 
          }
      return TableMap[tname]
}
func QueryLedger2(stub shim.ChaincodeStubInterface,tableName string,args []string)([]byte,error){
      var columns []shim.Column
      nCol:=GetNumberOfKeys(tableName)
      for i:=0;i<nCol;i++{
           colNext:=shim.Column{Value: &shim.Column_String_{String_:args[i]}}
           columns=append(columns,colNext)
}
      row,err:=stub.GetRow(tableName,columns)
     fmt.Println("Length or number of rows retrived ",len(row.Columns))
     if len(row.Columns)==0{
           jsonResp:="{\"error\":\" fail retrieving data"+args[0]+" .\"}"
          fmt.Println("error retriving data record for key="+args[0],"error:",jsonResp)
        return nil,errors.New(jsonResp)
        }
      Avalbytes:=row.Columns[nCol].GetBytes()
      fmt.Println("QueryLedger():successful-proceeding to process quest type")
      err=ProcessQueryResult(stub,Avalbytes,args)
     if err!=nil{
           fmt.Println("QueryLedger():cannot create object:",args[1])
    jsonResp:="{\"QueryLedger()error\":\" cannot create object for key"+args[0]+"\"}"
    return nil,errors.New(jsonResp)
    }
  return Avalbytes,nil
}

func UpdateLedger(stub shim.ChaincodeStubInterface,tableName string,keys []string,args []byte) error{
     nKeys:=GetNumberOfKeys(tableName)
     var record MyUserObject
     var err error 
     fmt.Println("the compare result:",strings.Compare(tableName,"MyUserTable")) 
     if strings.Compare(tableName,"MyUserTable")==0{
       record,err=JSONtoUser(args)
            fmt.Println("record.phone:",record.UserPhone)
        } 
     fmt.Println("iam entering updateledger")
     if nKeys<1{
        fmt.Println("At least 1 key must be porovide\n")
     }
     var columns []*shim.Column
     for i:=0;i<nKeys;i++{
        col:=shim.Column{Value:&shim.Column_String_{String_:keys[i]}}
        columns=append(columns,&col)
       }
    lastCol:=shim.Column{Value:&shim.Column_Bytes{Bytes:[]byte(args)}}
    columns=append(columns,&lastCol)
    row:=shim.Row{columns}
fmt.Println("i am inserting") 
   ok,err:=stub.InsertRow(tableName,row)
     if err!=nil{
fmt.Println(" inser row,err")
        return fmt.Errorf("UpdateLedger:InsertRow into "+tableName+" Table operateon failed,%s",err)
      }
     if !ok{
      fmt.Println("insert ok,but,existed")   
           return errors.New("UpdateLedger:Insert Row into"+tableName+" Table failed.Given keys"+keys[0]+"already existed")
        }
  fmt.Println("UpdateLedger:InsertRoew into "+tableName +" table successufull")
    return nil
}
func imageToByteArray(imageFile string)([]byte,string){
       file,err:=os.Open(imageFile)
       if err!=nil{
           fmt.Println("imageToByteAraay():cannot open image file",err)
         return nil,string("imageToByteAray():cannot open image file")
}
     defer file.Close()
     fileInfo,_:=file.Stat()
     var size int64=fileInfo.Size()
     bytes:=make([]byte,size)
     buff:=bufio.NewReader(file)
     _,err=buff.Read(bytes)
      if err!=nil{     
      return nil,string("imageToByteArray():cannot read image")
        }
     filetype:=http.DetectContentType(bytes)
     fmt.Println("imageToByteArray():",filetype)
     return bytes,filetype
}
func JSONtoAssetLog(ithis []byte)(MyAssetLog,error){
     item:=MyAssetLog{}
     err:=json.Unmarshal(ithis,&item)
     if err!=nil{
         fmt.Println("log error:",err)
         return item,err
     }
     return item,err
}
func JSONtoUser(user []byte)(MyUserObject,error){
     ur:=MyUserObject{}
     err:=json.Unmarshal(user,&ur)
     if err!=nil{
        fmt.Println("JSONtoUsr error:",err)
        return ur,err
     }
     fmt.Println("JSONtoUser created:",ur)
     return ur,err
}
func JSONtoAucReq(areq []byte)(MyAuctionRequest,error){
     ar:=MyAuctionRequest{}
     err:=json.Unmarshal(areq,&ar)
     if err!=nil{
        fmt.Println("JSONtoAucReq error:",err)
        return ar,err
      }
     return ar,err
}
func JSONtoBid(areq []byte)(MyBid,error){
      myHand:=MyBid{}
       err:=json.Unmarshal(areq,&myHand)
      if err!=nil{
          fmt.Println("JSONtoBid error:",err)
          return myHand,err
        }
       return myHand,err
}
func JSONtoTran(areq []byte)(MyAssetTransaction,error){
     at:=MyAssetTransaction{}
     err:=json.Unmarshal(areq,&at)
     if err!=nil{
         fmt.Println("JSONtoTran error:",err)
         return at,err
     }
     return at,err
}
func GetUser(stub shim.ChaincodeStubInterface,function string,args []string)([]byte,error){
var err error
Avalbytes,err:=QueryLedger(stub,"MyUserTable",args)
if err!=nil{
    fmt.Println("GetUser():Failed to Qeuery objet")
     jsonResp:="{\"Error\":\"failed to ger oject data faor "+args[0]+"\"}"
     return nil,errors.New(jsonResp)
}
if Avalbytes==nil{
    fmt.Println("Get user():inomplet query")
    jsonResp:="{\"error\":\"imcomplete for "+args[0]+"\"}"
    return nil,errors.New(jsonResp)
}
fmt.Println("GetUSr():successful")
return Avalbytes,nil
}
func QueryLedger(stub shim.ChaincodeStubInterface,tableName string,args []string)([]byte,error){
fmt.Println("enter query elgder")
    var columns []shim.Column
    nCol:=GetNumberOfKeys(tableName)
fmt.Println("ncol:",nCol,".tabelName:",tableName)
    for i:=0;i<nCol;i++{
        colNext:=shim.Column{Value: &shim.Column_String_{String_:args[i]}}
        columns=append(columns,colNext)
       }
     fmt.Println("append successful") 
     row,err:=stub.GetRow(tableName,columns)
    fmt.Println("lenth or number of rows retrieved",len(row.Columns))
if len(row.Columns)==0{
     jsonResp:="{\"error\":\"failed retrieving data "+args[0] +".\"}"
     fmt.Println("fail retrieving for key = ",args[0],"error",jsonResp)
     return nil,errors.New(jsonResp)
}
Avalbytes:=row.Columns[nCol].GetBytes()
fmt.Println("successful")
err=ProcessQueryResult(stub,Avalbytes,args)
if err!=nil{
      fmt.Println("queryLedger():cannot create object:",args[1])
      jsonResp:="{\"error\":\"cannot create "+args[0] +"\"}"
     return nil,errors.New(jsonResp)
}
  return Avalbytes,nil
}
func GetList(stub shim.ChaincodeStubInterface,tableName string,args []string)([]shim.Row,error){
       var columns []shim.Column
        nKeys:=GetNumberOfKeys(tableName)
        nCol:=len(args)
         if nCol<1{
               fmt.Println("at least one key\n")
               return nil,errors.New("getlist failed")
         }
               for i:=0;i<nCol;i++{
                  colNext:=shim.Column{Value:&shim.Column_String_{String_:args[i]}}
                  columns=append(columns,colNext)
                }
         rowChannel,err:=stub.GetRows(tableName,columns)
         if err!=nil{
               return nil,fmt.Errorf(" operation fail.%s",err)
         }
         var rows []shim.Row
         for{
              select{
                case row,ok:=<-rowChannel:
                      if !ok{
                            rowChannel=nil
                       }else{
                             rows=append(rows,row)
                       }
                 }
                 if rowChannel==nil{
                       break
                 }
            }
            fmt.Println("keys retrieved:",nKeys)
            fmt.Println("rows retrieved:",len(rows))
            return rows,nil
}
func GetUserListByCat(stub shim.ChaincodeStubInterface,function string,args []string)([]byte,error){
      if len(args)<1{
           fmt.Println("getUSerListBycat():incorect number of arguents")
           return nil,errors.New("ccreateUserObjetct():incorrect")
       }
       rows,err:=GetList(stub,"MyUserCatTable",args)
       if err!=nil{
              return nil,fmt.Errorf("ger failed.error marshaling json:%s",err)
      }
      nCol:=GetNumberOfKeys("MyUserCatTable")
      tlist:=make([]MyUserObject,len(rows))
      for i:=0;i<len(rows);i++{
          ts:=rows[i].Columns[nCol].GetBytes()
          uo,err:=JSONtoUser(ts)
          if err!=nil{
              fmt.Println("GerUserListByCat()failed:ummarshasll error")
              return nil,fmt.Errorf("operaion faile.%s",err)
           }
           tlist[i]=uo
       }
       jsonRows,_:=json.Marshal(tlist)
       return jsonRows,nil
}
func GetAssetListByCat(stub shim.ChaincodeStubInterface,function string,args []string)([]byte,error){
    if len(args)<1{
           fmt.Println("getAssetListByCat():ncorrect arguent")
           return nil,errors.New("incorect argument")
      }
      rows,err:=GetList(stub,"MyAssetCatTable",args)
      if err!=nil{
           return nil,fmt.Errorf("getItem List cat .errot gerlist:%s",err)
       }
       nCol:=GetNumberOfKeys("MyAssetCatTable")
       tlist:=make([]MyAssetObject,len(rows))
       for i:=0;i<len(rows);i++{
                ts:=rows[i].Columns[nCol].GetBytes()
                io,err:=JSONtoAR(ts)
                if err!=nil{
                     fmt.Println("unmarshall error")
                     return nil,fmt.Errorf("operation fial.%s",err)
                 }
                 io.AssetImage=[]byte{}
                 fmt.Println("list asset,aes-key:",io.AES_Key)  
                tlist[i]=io
        }
       jsonRows,_:=json.Marshal(tlist)
       return jsonRows,nil
}
func PostUser(stub shim.ChaincodeStubInterface,function string,args []string)([]byte,error){
       record,err:=CreateUserObject(args[0:])
        fmt.Println("args[5]",args[5])
      if err!=nil{
         return nil,err
     }
     buff,err:=UsertoJSON(record)
     if err!=nil{
        fmt.Println("PostUserObject():failed cannot create object:",args[1])
        return nil,errors.New("PostUserObject():faile cannot create object:"+args[1])
     }else {
          keys:=[]string{args[0]}
             err=UpdateLedger(stub,"MyUserTable",keys,buff)
           if err!=nil{
              fmt.Println("PostUser():write error while inserting recode")
              return nil,err
            }
          _,err=PostCreditLog(stub,record,record.UserAmount,"DEFAULT")
          if err!=nil{
               fmt.Println("Postcredit og werite error")
               return nil,err
         } 
          keys=[]string{"2016",args[3],args[0]}
     err=UpdateLedger(stub,"MyUserCatTable",keys,buff)
     if err!=nil{
         fmt.Println("PostUser():write error wihle inserting recode into usercatTable")
      }
   }
   return buff,err
}
