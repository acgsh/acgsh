package main

import (
	"log"
	"strconv"
	"time"

	"strings"

	"github.com/acgshare/acgsh/db"
	"github.com/acgshare/acgsh/rpc"
)

const sync_posts_number = 2000

type updatePostsInfo struct {
	maxK  int64
	minK  int64
	lastK int64
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func syncTimeLine() {
	// Update publishers
	twisterUsers, err := rpc.GetFollowing(adminTwisterUsername)
	if err != nil {
		log.Println("Error: can not fetch following users for", adminTwisterUsername, "from Twister RPC server.", err)
	}
	log.Println(twisterUsers)
	log.Println("Get", len(*twisterUsers), "following users for", adminTwisterUsername)
	db.AddPublishersIfNotExist(twisterUsers)

	// Get all publishers' username, max id, since id and latest id (k), from DB
	publishers, err := db.GetPublishers()
	if err != nil {
		log.Println(err)
		log.Printf("Error: syncTimeLine GetPublishers: %+v\n", publishers)
	}
	//log.Printf("%+v\n", publishers)

	// Get posts from Twister for all publishers
	var gup getUserPosts
	for name, data := range publishers {
		gup.addUser(name, data.Max, data.Since)
	}

	newPosts, err := rpc.GetPosts(sync_posts_number, gup.data)
	//log.Printf("%+v\n", presult)
	log.Println(len(*newPosts), "new posts")

	// Update publishers sync data with info from retrieved posts.
	var newShPosts db.ShPosts
	var newShPubReplyPosts db.ShPubReplyPosts
	hasNewPostsPublishers := make(map[string]updatePostsInfo)
	for _, tp := range *newPosts {

		// Make sure current post's username is in the publishers list.(May contain promoted post)
		_, ok := publishers[tp.Userpost.N]
		if !ok {
			continue
		}

		// Get lastK, may be nil pointer.
		var lastK int64
		if tp.Userpost.Lastk == nil {
			lastK = -1
		} else {
			lastK = *(tp.Userpost.Lastk)
		}

		upi, ok := hasNewPostsPublishers[tp.Userpost.N]
		if ok {
			upi.lastK = lastK
			if upi.maxK < tp.Userpost.K {
				upi.maxK = tp.Userpost.K
			}
			if upi.minK > tp.Userpost.K {
				upi.minK = tp.Userpost.K
			}
			hasNewPostsPublishers[tp.Userpost.N] = upi
		} else {
			hasNewPostsPublishers[tp.Userpost.N] = updatePostsInfo{
				lastK: lastK,
				maxK:  tp.Userpost.K,
				minK:  tp.Userpost.K,
			}

		}

		// Add New ShPost
		if tp.Userpost.Reply == nil && tp.Userpost.Rt == nil {
			btih, category, fileSize, title, ok := retrieveMagnetInfo(tp.Userpost.Msg)
			if ok {
				shPost := db.ShPost{
					Msg:      tp.Userpost.Msg,
					N:        tp.Userpost.N,
					K:        tp.Userpost.K,
					Lastk:    lastK,
					Time:     tp.Userpost.Time,
					Category: category,
					Title:    title,
					Magnet:   btih,
					Size:     fileSize,
				}

				newShPosts = append(newShPosts, shPost)
				//fmt.Println(index, tp.Userpost.N, tp.Userpost.K, lastK)
			}
		}

		// Add new ShPubReplyPost
		if tp.Userpost.Reply != nil {
			if tp.Userpost.Reply.N == tp.Userpost.N {
				shPubReplyPosts := db.ShPubReplyPost{
					Msg:   tp.Userpost.Msg,
					Time:  tp.Userpost.Time,
					N:     tp.Userpost.N,
					K:     tp.Userpost.Reply.K,
					Lastk: lastK,
				}
				newShPubReplyPosts = append(newShPubReplyPosts, shPubReplyPosts)
				//fmt.Println(index, tp.Userpost.N, tp.Userpost.K, lastK)
			}
		}
	}

	// Save posts to DB
	err = db.AddPosts(&newShPosts)
	if err != nil {
		log.Println(err)
		log.Printf("Error: syncTimeLine db.AddPosts: %+v\n", len(newShPosts))
	}
	err = db.AddPublishersReplyPosts(&newShPubReplyPosts)
	if err != nil {
		log.Println(err)
		log.Printf("Error: syncTimeLine db.AddPublishersReplyPosts: %+v\n", len(newShPubReplyPosts))
	}

	//log.Printf("%+v\n", hasNewPostsPublishers)

	// Update publishers in DB
	updatedPublishers := make(map[string]db.SyncData)
	for name, postsInfo := range hasNewPostsPublishers {
		var newSd db.SyncData
		oldSd := publishers[name]

		newSd.Latest = oldSd.Latest
		if postsInfo.maxK > oldSd.Latest {
			newSd.Latest = postsInfo.maxK
		}

		newSd.Since = oldSd.Since
		lastK := min(postsInfo.lastK, postsInfo.minK-1)
		if lastK <= oldSd.Since {
			newSd.Since = newSd.Latest
			newSd.Max = max_post_id
		} else {
			newSd.Max = lastK
		}

		updatedPublishers[name] = newSd
	}

	err = db.UpdatePublishers(&updatedPublishers)
	if err != nil {
		log.Println(err)
		log.Printf("Error: syncTimeLine: %+v\n", updatedPublishers)
	}

}

func retrieveMagnetInfo(ss string) (string, string, uint64, string, bool) {

	lowerS := strings.ToLower(ss)
	// If not magnet return
	if !(strings.Contains(lowerS, "magnet")) {
		return "", "", 0, "", false
	}

	idx := strings.Index(lowerS, "magnet")

	magnetString := ss[idx:]

	idx = strings.Index(magnetString, "?")
	if idx == -1 {
		return "", "", 0, "", false
	}
	if len(magnetString) <= idx+1 {
		return "", "", 0, "", false
	}
	paramString := magnetString[idx+1:]
	//log.Println(paramString)

	params := strings.Split(paramString, "&")

	var category, title string
	var fileSize uint64

	for _, param := range params {
		//log.Println(param)
		k := strings.Split(param, "=")
		//log.Println(k, len(k))
		if len(k) >= 2 {
			n := strings.ToLower(k[0])
			if strings.Contains(n, "xl") {
				var err error
				fileSize, err = strconv.ParseUint(k[1], 10, 64)
				if err != nil {
					fileSize = 0
				}
			}
			if strings.Contains(n, "x.c") {
				category = k[1]
			}
			if strings.Contains(n, "dn") {
				title = k[1]
			}
		}
	}

	magnetString = strings.TrimSpace(magnetString)
	category = strings.TrimSpace(category)
	title = strings.TrimSpace(title)

	return magnetString, category, fileSize, title, true
}

func runSyncTimeLine() {
	for {
		syncTimeLine()
		time.Sleep(60 * time.Second)
	}
}

// todo: add x.l
