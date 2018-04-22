package main

import (
	//"golang.org/x/oauth2"
	//"golang.org/x/oauth2/google"
	//"google.golang.org/api/firestore/v1beta1"
	//firebase "firebase.google.com/go"
	//"firebase.google.com/go/auth"
	"cloud.google.com/go/firestore"
	"golang.org/x/net/context"
	//"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	//"strings"
	_ "strconv"
	"time"
)

var del_re *regexp.Regexp = regexp.MustCompile("<a href=\"[^?]+\\?bo_table=[^&]+&wr_id=([0-9]+)[^\"]*\">(<[^<]+){7,20}.+<img src='\\.\\./skin/board/[^/]+/img/icon_secret\\.gif'")
var title_re *regexp.Regexp = regexp.MustCompile("/img/icon_subject.gif\"[^<]+<a href=\"[^?]+\\?bo_table=[^&]+&wr_id=([0-9]+)[^\"]*\">")

func main() {
	ctx := context.Background()
	opt := option.WithCredentialsFile("etorrent-cb3fd-firebase-adminsdk-jdmom-bc2c5fa9fb.json")
	client, err := firestore.NewClient(ctx, "etorrent-cb3fd", opt)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()
	da := client.Collection("etorrent").Doc("deleted_articles")

	/*doc, err := da.Collection("etohumor").Doc("1532842").Get(ctx)
	if err != nil {
		log.Println(err)
	}
	d := doc.Data()
	log.Println(d["prev"], d["next"])*/

	for {
		for _, bo_table := range []string{"etohumor", "etoboard"} { //, "star", "movie", "any",
			//for page := 1; page <= 2; page++ {
			r, err := http.Get("https://etorrent.co.kr/bbs/board.php?bo_table=" + bo_table) // + "&page=" + strconv.Itoa(2))
			if err != nil {
				return
			}
			defer r.Body.Close()
			body, _ := ioutil.ReadAll(r.Body)
			bodystr := string(body)
			find := del_re.FindAllString(bodystr, -1)
			if len(find) > 0 {
				for _, f := range find {
					wr_id := del_re.FindStringSubmatch(f)[1]
					doc := da.Collection(bo_table).Doc(wr_id)
					_, err := doc.Get(ctx)
					if err != nil {
						log.Println(f)
						log.Println(wr_id)
						wrIdSearch := title_re.FindAllString(bodystr, -1)
						prevWrId := title_re.FindStringSubmatch(wrIdSearch[0])[1]
						// ignore first ( doesn't have prev article to link )
						for i, ft := range wrIdSearch[1:] {
							wrId := title_re.FindStringSubmatch(ft)[1]
							if wrId == wr_id {
								nextWrId := title_re.FindStringSubmatch(wrIdSearch[i+1])[1]
								_, err = doc.Set(ctx, map[string]interface{}{
									"prev": prevWrId,
									"next": nextWrId,
								})
								if err != nil {
									log.Println(err)
								}
								_, err = da.Collection(bo_table).Doc(prevWrId).Set(ctx, map[string]interface{}{
									"next": nextWrId,
								})
								if err != nil {
									log.Println(err)
								}
								_, err = da.Collection(bo_table).Doc(nextWrId).Set(ctx, map[string]interface{}{
									"prev": prevWrId,
								})
								if err != nil {
									log.Println(err)
								}
							}
							prevWrId = wrId
						}
					}
				}
				//}
			}
		}
		log.Println("rechecked...")
		time.Sleep(10 * time.Minute)
	}
}

// https://firestore.googleapis.com/v1beta1/projects/etorrent-cb3fd/databases/(default)/documents/etorrent/deleted_articles/etohumor/1532842?key={YOUR_API_KEY}
