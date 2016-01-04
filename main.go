package main

import (
	"log"

	"github.com/BurntSushi/toml"
	"github.com/acgshare/acgsh/db"
	"github.com/acgshare/acgsh/rpc"
)

var adminTwisterUsername string
var config acgshConfig

const (
	max_post_id = 99999999
)

type acgshConfig struct {
	TwisterUsername string
	TwisterServer   string
	HttpServerPort  string
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Load config

	if _, err := toml.DecodeFile("acgsh.conf", &config); err != nil {
		log.Fatalln("Error: can not load acgsh.conf", err)
		return
	}
	adminTwisterUsername = config.TwisterUsername

	//Init DB
	db.Init()
	defer db.Close()

	rpc.SetAddress(config.TwisterServer)

	go runSyncTimeLine()

	//btih, category, fileSize, title, ok := retrieveMagnetInfo("#acgsh maGnet:? dn = =& xt=urn:btih:A3TU7P63QSNXXSYN2PDQYDZV4IYRU2CG& x.C =       動畫 &xl=123124&dn=[诸神字幕组][高校星歌剧][High School Star Musical][12][繁日双语字幕][720P][CHT MP4]")
	//println(btih, category, fileSize, title, ok)

	startHttpServer()
}

// todo: httpjsonrpcClient unmarshal json error
// todo: httpjsonrpcClient no connection error
// todo: httpjsonrpcClient log.fatal modify
