package main

import (
	"errors"
	"fmt"
	"sync"
	_ "strconv"
	_ "encoding/hex"
	
	"github.com/op/go-logging"
	"github.com/hyperledger/fabric/core/chaincode/shim"	
	proto "github.com/golang/protobuf/proto"
	
	persistpb "gamecenter.mobi/paicode/protos"
	sec "gamecenter.mobi/paicode/chaincode/security" 
	tx "gamecenter.mobi/paicode/chaincode/transaction"
)

type paiStatus struct{
	totalPai 	int64
	frozenPai 	int64
}

type PaiChaincode struct {
	globalLock 	sync.RWMutex
	cacheOK    	bool
	paistat		paiStatus
}

const (

	global_setting_entry string = "global_setting"
	
)

var privilege_Def = map[string]string{
	tx.Admin_funcs: sec.AdminPrivilege,
	tx.Manage_funcs: sec.ManagerPrivilege,
	tx.User_funcs: sec.DelegatePrivilege}

var logger = logging.MustGetLogger("chaincode")

func (s *paiStatus) init(set *persistpb.DeploySetting){
	s.totalPai = set.TotalPais
	s.frozenPai = set.UnassignedPais
}

func (s *paiStatus) set(set *persistpb.DeploySetting){
	 set.TotalPais = s.totalPai
	 set.UnassignedPais = s.frozenPai
}

func (t *PaiChaincode) updateCache(stub shim.ChaincodeStubInterface) error{
	t.globalLock.RLock()
	defer t.globalLock.RUnlock()
	
	if !t.cacheOK{
		t.globalLock.Lock()
		defer t.globalLock.RUnlock()
		
		rawset, err := stub.GetState(global_setting_entry)
		if err != nil{
			return err
		}
		
		if rawset == nil{
			return errors.New("FATAL: No global setting found")
		}
		
		setting := &persistpb.DeploySetting{}
		err = proto.Unmarshal(rawset, setting)
		
		if err != nil{
			return err
		}
		
		sec.InitSecHelper(setting)
		t.paistat.init(setting)
		logger.Info("Update global setting:", setting)
		
		t.cacheOK = true	
	}
	
	return nil
}

func (t *PaiChaincode) saveGlobalStatus(stub shim.ChaincodeStubInterface) error{
	t.globalLock.RLock()
	defer t.globalLock.RUnlock()
	
	if !t.cacheOK {
		return errors.New("FATAL: Invalid cache")
	}	
	
	//a Write After Read process
	rawset, err := stub.GetState(global_setting_entry)
	if err != nil{
		return err
	}
	
	if rawset == nil{
		t.cacheOK = false
		return errors.New("FATAL: No global setting found")
	}
	
	setting := &persistpb.DeploySetting{}
	err = proto.Unmarshal(rawset, setting)
	
	if err != nil{
		return err
	}
	
	t.paistat.set(setting)
	
	logger.Info("Save current global setting:", setting)	
	
	rawset, err = proto.Marshal(setting)
	if err != nil{
		return err
	}
	
	return stub.PutState(global_setting_entry, rawset)	
}

func (t *PaiChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	return nil, nil
}

func (t *PaiChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	
	if err := t.updateCache(stub); err != nil{
		return nil, err
	}
	
	rolePriv, region := sec.Helper.GetPrivilege(stub)
	
	funcGrp := function[:tx.FuncPrefix]	
	expectPriv := privilege_Def[funcGrp]
	
	//check priviledge
	if !sec.Helper.VerifyPrivilege(rolePriv, expectPriv){ 
		sec.Helper.ActiveAudit(stub, fmt.Sprintf("Call function <%s> without require priviledge", function))
		return nil, errors.New("No priviledge")
	}	
	
	var err error
	switch funcGrp{
		case tx.Admin_funcs:
		case tx.Manage_funcs:
		case tx.User_funcs:
			err = t.handleUserFuncs(stub, function, region, args)
		default:
			return nil, errors.New("Function group not exist or invokable")
	}

	return nil, err
}


func (t *PaiChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if err := t.updateCache(stub); err != nil{
		return nil, err
	}

	return nil, nil
}

func main() {
	err := shim.Start(new(PaiChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
