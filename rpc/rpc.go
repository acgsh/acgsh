package rpc

import (
	"encoding/json"
	"log"

	"fmt"
	"github.com/acgshare/Go-HTTP-JSON-RPC/httpjsonrpc"
)

var Address string
var id int64

func SetAddress(newAddress string) {
	Address = newAddress
}

func Follow(name string, followUserNames []string) (interface{}, error) {
	//
	id++
	resp, err := httpjsonrpc.Call(Address, "follow", id, []interface{}{name, followUserNames})
	if err != nil {
		log.Println(err)
		return resp, err
	}
	result := resp.Result
	//log.Println(resp)

	return result, err
}

func UnFollow(name string, followUserNames []string) (interface{}, error) {
	//
	id++
	resp, err := httpjsonrpc.Call(Address, "unfollow", id, []interface{}{name, followUserNames})
	if err != nil {
		log.Println(err)
		return resp, err
	}
	result := resp.Result
	//log.Println(resp)

	return result, err
}

func GetFollowing(name string) (*[]string, error) {
	//
	id++
	resp, err := httpjsonrpc.Call(Address, "getfollowing", id, []interface{}{name})
	if err != nil {
		log.Println(err)
		log.Printf("Error: RPC GetFollowing: %+v", resp)
		return nil, err
	}
	result := resp.Result
	if resp.Error != nil {
		log.Printf("Error: RPC GetFollowing: %+v", resp)
		return nil, fmt.Errorf("RPC error")
	}

	var following []string
	err = json.Unmarshal(result, &following)
	if err != nil {
		log.Println(err)
		log.Printf("Error: RPC GetFollowing Unmarshal: %+v", result)
		return nil, err
	}

	return &following, nil
}
func ListWalletUsers() (interface{}, error) {
	//
	id++
	resp, err := httpjsonrpc.Call(Address, "listwalletusers", id, nil)
	if err != nil {
		log.Println(err)
		return resp, err
	}
	result := resp.Result
	//log.Println(resp)

	return result, err
}
func GetPosts(count int64, params []interface{}) (*TwisterPosts, error) {
	//
	id++
	resp, err := httpjsonrpc.Call(Address, "getposts", id, []interface{}{count, params})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	result := resp.Result
	if resp.Error != nil {
		log.Printf("Error: RPC GetPosts: %+v", resp)
		return nil, fmt.Errorf("RPC error")
	}
	//fmt.Println(string(result))
	var tp TwisterPosts
	err = json.Unmarshal(result, &tp)
	if err != nil {
		log.Println(err)
		log.Printf("Error: RPC GetPosts Unmarshal: %+v", result)
		return nil, err
	}

	return &tp, nil
}
