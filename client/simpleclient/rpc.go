package main

import (
	"fmt"
	_ "errors"
	
	"github.com/spf13/cobra"
	"github.com/hyperledger/fabric/peerex"
)

//the chaincode id of gamecenter.mobi/paicode/chaincode
const defPaicodeName string = "50637ebc88e9c0f2ea9d240784b491c4fde8ebd177a95fbc2f087312111affef1898fea4c267ff1084db244de6c6860f4367b700659d44b7b47fabda27347c23"

var rpcCmd = &cobra.Command{
	Use:   "rpc [command...]",
	Short: fmt.Sprintf("rpc commands."),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error{
		
		if defClient.IsRpcReady(){
			return nil
		}
		
		conn := peerex.ClientConn{}
		err := conn.Dialdefault()
		if err != nil{
			return err
		}		
		
		defClient.PrepareRpc(conn) 
		defClient.Rpc.Rpcbuilder.ChaincodeName = defPaicodeName
		return nil
	},
}

var codenameCmd = &cobra.Command{
	Use:   "chaincode [id]",
	Short: fmt.Sprintf("query or set chaincode id"),
	Run: func(cmd *cobra.Command, args []string) {
		
		if len(args) == 0{
			fmt.Println("Current chaincode id is", defClient.Rpc.Rpcbuilder.ChaincodeName)
			return
		}
		
		if len(args) != 1{
			fmt.Println("Too many arguments")
			return
		}
		
		defClient.Rpc.Rpcbuilder.ChaincodeName = args[0]
		fmt.Println("Set chaincode id as", defClient.Rpc.Rpcbuilder.ChaincodeName)
	},	
}

var userCmd = &cobra.Command{
	Use:   "user [command...]",
	Short: fmt.Sprintf("user commands."),
}

var registerCmd = &cobra.Command{
	Use:       "register",
	Short:     fmt.Sprintf("Register a public key"),
	RunE: func(cmd *cobra.Command, args []string) error{
		
		msg, err := defClient.Rpc.Registry(args...)
		if err != nil{
			return err
		}
		
		fmt.Println("Registry public key ok, TX id is", msg)
		return nil
	},
}

var fundCmd = &cobra.Command{
	Use:       "fund <to:addr> <amount> [message]",
	Short:     fmt.Sprintf("Register a public key"),
	RunE: func(cmd *cobra.Command, args []string) error{
		
		msg, err := defClient.Rpc.Fund(args...)
		if err != nil{
			return err
		}
		
		fmt.Println("Fund ok, TX id is", msg)
		return nil
	},
}

var queryUserCmd = &cobra.Command{
	Use:       "query <userid>",
	Short:     fmt.Sprintf("Query the status of a user"),
	RunE: func(cmd *cobra.Command, args []string) error{
		
		ret, err := defClient.Rpc.QueryUser(args...)
		if err != nil{
			return err
		}
		
		fmt.Println("---------------- Query user", args[0], "----------------")
		fmt.Println(string(ret))
		fmt.Println("-------------------------------------------------------------------------------")
		return nil
	},
}

var queryGlobalCmd = &cobra.Command{
	Use:       "query",
	Short:     fmt.Sprintf("Query the status of chaincode"),
	RunE: func(cmd *cobra.Command, args []string) error{
		
		ret, err := defClient.Rpc.QueryGlobal(args...)
		if err != nil{
			return err
		}
		
		fmt.Println("---------------- Chaincode status ----------------")
		fmt.Println(string(ret))
		fmt.Println("--------------------------------------------------")
		return nil
	},
}


func init(){
	userCmd.AddCommand(registerCmd)
	userCmd.AddCommand(queryUserCmd)
	userCmd.AddCommand(fundCmd)
	
	rpcCmd.AddCommand(userCmd)
	rpcCmd.AddCommand(codenameCmd)
	rpcCmd.AddCommand(queryGlobalCmd)
}