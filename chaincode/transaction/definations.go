package transaction

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	persistpb "gamecenter.mobi/paicode/protos"
)	

//each function name has a 4-bytes prefix
const (
	
	FuncPrefix int = 4
	
	Admin_funcs string = "ADMN"
	Manage_funcs string = "MANG"
	User_funcs string = "USER"
	Query_funcs string = "QURY"
	
)

var UserFund string = User_funcs + "_FUND"
var UserAuthChaincode string = User_funcs + "_AUTHCHAINCODE"

type UserTx interface{
	HandleUserTx(*persistpb.UserData, shim.ChaincodeStubInterface, []string) (map[string]*persistpb.UserData, error) 
}

var UserTxMap = map[string]UserTx{} 

