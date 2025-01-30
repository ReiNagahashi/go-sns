package seeds

import (
	"fmt"
	"go-sns/database"
	"go-sns/database/dataAccess/implementations/postImpl"
	"go-sns/database/dataAccess/implementations/userImpl"
	"go-sns/database/dataAccess/interfaces"
	"log"
	"sync"
	"time"

	"github.com/go-faker/faker/v4"
)

func PostSeed(sqlite *database.SqliteBase, wg *sync.WaitGroup) {
	defer wg.Done()
	stmt, err := sqlite.DbConnection.Prepare("INSERT INTO posts (title, description, submitted_by, created_at, updated_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var userDao interfaces.UserDAO = userImpl.NewUserDAOImpl(sqlite)
	fmt.Printf("%v", userDao)

	users, err := userDao.GetAll()
	if err != nil {
		log.Fatalln(err)
	}

	if len(users) == 0 {
		return
	}

	user_ids := make([]int, len(users))
	for i := range user_ids {
		user_ids[i] = users[i].GetId()
	}

	numRecords := 5

	for i := 0; i < numRecords; i++ {
		title := faker.Paragraph()
		description := faker.Paragraph()

		randomIntSubmitted_by, err := faker.RandomInt(0, len(users)-1)
		if err != nil {
			log.Fatalln(err)
		}
		submitted_by := user_ids[randomIntSubmitted_by[0]]		

		now := time.Now()
		_, err = stmt.Exec(title, description, submitted_by, now, now)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func PostFavoriteSeed(sqlite database.Database, wg *sync.WaitGroup) {
	defer wg.Done()
	// シードするいいねの数
	favoriteCnt := 10

	var postDao interfaces.PostDAO = postImpl.NewPostDAOImpl(sqlite)
	// 現在のポスト全てを取ってくる→n個
	posts, err := postDao.GetAll()
	if err != nil {
		log.Println(err)
		return
	}
	var userDao interfaces.UserDAO = userImpl.NewUserDAOImpl(sqlite)
	// nこのポストはそれぞれposts_id[i~N]のidを持つ
	postsIds := make([]int, len(posts))
	for i, post := range posts {
		postsIds[i] = post.GetId()
	}
	// 現在のユーザー全てを取ってくる→m個
	users, err := userDao.GetAll()
	if err != nil {
		log.Println(err)
		return
	}
	// mこのユーザーはそれぞれusers_id[j~M]のidを持つ
	usersIds := make([]int, len(users))
	for i, user := range users {
		usersIds[i] = user.GetId()
	}

	postsIndexAddedFavorite, err := faker.RandomInt(0, len(posts)-1, favoriteCnt)
	if err != nil {
		log.Fatalln(err)
	}

	for _, postIndex := range postsIndexAddedFavorite {
		randomUserIndex, err := faker.RandomInt(0, len(users)-1)
		if err != nil {
			log.Fatalln(err)
		}
		// 抽出されたi個目において、更にusers_idから1~mの間からランダムなuser_idを1つ取得し、そのuser_idとポストidをそれぞれカラム値にして挿入
		if err := postDao.AddFavorite(usersIds[randomUserIndex[0]], postsIds[postIndex]); err != nil{
            log.Fatalln(err)
        }
	}

}
