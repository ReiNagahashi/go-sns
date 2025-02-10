package postImpl

import (
	"errors"
	"fmt"
	"go-sns/database"
	"go-sns/database/dataAccess/implementations/userImpl"
	"go-sns/database/dataAccess/interfaces"
	"go-sns/models"
	"time"
)



type PostDAOImpl struct{
	db database.Database
}

func NewPostDAOImpl(db database.Database) *PostDAOImpl{
	return &PostDAOImpl{db: db}
}


func (p PostDAOImpl) AddFavorite(userId, postId int) error{
	posts, err := p.GetAll()
    if err != nil{
		return errors.New("action=PostDAOImpl.AddFavorite msg=Error executing query: " + err.Error())
    }

    // nこのポストはそれぞれposts_id[i~N]のidを持つ
    postMap := make(map[int]bool)
    for _, post := range posts{
        postMap[post.GetId()] = true
    }

    var userDao interfaces.UserDAO = userImpl.NewUserDAOImpl(p.db)
    // 現在のユーザー全てを取ってくる→m個
    users, err := userDao.GetAll()
	if err != nil{
		return errors.New("action=PostDAOImpl.AddFavorite msg=Error executing query: " + err.Error())
    }
    // mこのユーザーはそれぞれusers_id[j~M]のidを持つ
    userMap := make(map[int]bool)
    for _, user := range users{
        userMap[user.GetId()] = true
    }

	_, postIsExisted := postMap[postId]
	_, userIsExisted := userMap[userId]

	if !postIsExisted || !userIsExisted{
		return errors.New("action=PostDAOImpl.AddFavorite msg=post_id or user_id is invalid")
	}

	query := "INSERT INTO users_posts (user_id, post_id) VALUES(?,?)"

	if err := p.db.PrepareAndExecute(query, userId, postId); err != nil{
		return errors.New("action=PostDAOImpl.AddFavorite msg=Error executing query: " + err.Error())
	}

	return nil
}

func (p PostDAOImpl) Create(postData models.Post) error{
	if(postData.GetId() != -1){
		return errors.New("action=PostDAOImpl.Create msg=Cannot create a post data with an existing ID. id: " + string(rune(postData.GetId())))
	}
	
	query := "INSERT INTO posts (title, description, submitted_by, created_at, updated_at) VALUES(?,?,?,?,?)"

	if err := p.db.PrepareAndExecute(query, postData.GetFields()...); err != nil {
		return errors.New("action=PostDAOImpl.Create msg=Error executing query: " + err.Error())
	}

	return nil
}


func (p PostDAOImpl) Delete(id int) error{
	if err := p.db.PrepareAndExecute("DELETE FROM posts WHERE id = ?", id); err != nil{
		return errors.New("action=PostDAOImpl.Delete msg=Error executing query: " + err.Error())	
	}

	return nil
}

func (p PostDAOImpl) DeleteFavorite(userId, postId int) error{
	if err := p.db.PrepareAndExecute("DELETE FROM users_posts WHERE user_id = ? AND post_id = ?", userId, postId); err != nil{
		return errors.New("action=PostDAOImpl.DeleteFavorite msg=Error executing query: " + err.Error())
	}

	return nil
}

func (p PostDAOImpl) GetAll() ([]models.Post, error){
	query := "SELECT * FROM posts"

	postsRecords, err := p.db.PrepareAndFetchAll(query)
	if err != nil {
		return nil, errors.New("action=PostDAOImpl.PrepareAndFetchAll msg=Error executing query: " + err.Error())
	}

	posts := p.resultsToPosts(postsRecords)

	if err := p.updateFavoriteTable(); err != nil{
		return nil, errors.New("action=PostDAOImpl.updateFavoriteTable msg=Error executing query: " + err.Error())
	}

	err = p.initFavorite(posts)
	if err != nil{
		return nil, errors.New("action=PostDAOImpl.initPostFavorite msg=Error executing query: " + err.Error())
	}

	return posts, nil
}

// getAll,getPosts→updatePostFavoritesとやって消えた投稿idを持っているデータを削除するメソッドを実行→initFavorite
// 更新方法は、全ての投稿、全てのユーザーデータをメソッド内で呼び出して、それらがfavoritesRecordsを展開した値を持っていない場合は、そのfavoriteRecordは不適切なので、削除処理をする
func (p PostDAOImpl)initFavorite(posts []models.Post) error{
	postIds := make(map[int]int)
	for i := 0; i < len(posts); i++{
		postIds[posts[i].GetId()] = i
	}

	favoritesRecords, err := p.getAllFavorites()
	if err != nil{
		return err
	}
	var userdao interfaces.UserDAO = userImpl.NewUserDAOImpl(p.db)
	
	for _, favoriteRecord := range favoritesRecords{
		favoriteUserId := int(favoriteRecord["user_id"].(int64))
		favoritePostId := int(favoriteRecord["post_id"].(int64))

		user, err := userdao.GetById(favoriteUserId)
		if err != nil{
			return errors.New("action=PostDAOImpl.GeyById msg=Error executing query: " + err.Error())
		}
		postIndex := postIds[favoritePostId]

		favoriteUsers := posts[postIndex].GetFavoriteUsers()
		favoriteUsers = append(favoriteUsers, *user)
		posts[postIndex].SetFavoriteUsers(favoriteUsers)
	}

	return nil
}

func (p PostDAOImpl) getAllFavorites()([]map[string]interface{}, error){
	// users_postsからレコードを全て取得
	query := "SELECT * FROM users_posts"
	favofavoritesRecords, err := p.db.PrepareAndFetchAll(query)
	if err != nil{
		return nil, errors.New("action=PostDAOImpl.getAllFavorites msg=Error executing query: " + err.Error())
	}
	return favofavoritesRecords, nil
}


func (p PostDAOImpl) GetById(id int) (*models.Post, error){
	results,err := p.db.PrepareAndFetchAll("SELECT * FROM posts WHERE id = ?", id)
	if err != nil {
		return nil, errors.New("action=PostDAOImpl.GetById msg=Error executing query: " + err.Error())
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no post found with id %d", id)
	}

	post := p.resultToPost(results[0])

	return &post, nil
}


func (p PostDAOImpl) GetPosts(limit int)([]models.Post, error){
	recordNums,err := p.db.GetTableLength("posts")
	if err != nil{
		return nil, errors.New("action=PostDAOImpl.GetTableLength msg=Error executing query: " + err.Error())
	}

	if limit < 0 || limit > recordNums{
		limit = recordNums
	}

	query := "SELECT * FROM posts LIMIT ?"

	postsRecords, err := p.db.PrepareAndFetchAll(query, []interface{}{limit}...)
	if err != nil {
		return nil, errors.New("action=PostDAOImpl.PrepareAndFetchAll msg=Error executing query: " + err.Error())
	}

	posts := p.resultsToPosts(postsRecords)

	if err := p.updateFavoriteTable(); err != nil{
		return nil, errors.New("action=PostDAOImpl.updateFavoriteTable msg=Error executing query: " + err.Error())
	}
	
	err = p.initFavorite(posts)
	if err != nil{
		return nil, errors.New("action=PostDAOImpl.initPostFavorite msg=Error executing query: " + err.Error())
	}

	return posts, nil
}


func (p PostDAOImpl) resultToPost(post map[string]interface{}) models.Post{
	return *models.NewPost(
		int(post["id"].(int64)),
		int(post["submitted_by"].(int64)),
		post["title"].(string),
		post["description"].(string),
		*models.NewDateTimeStamp(post["created_at"].(time.Time), post["updated_at"].(time.Time)),
		[]models.User{})
}


func (p PostDAOImpl) resultsToPosts(results []map[string]interface{}) []models.Post{
	posts := make([]models.Post, 0)
	
	for _, result := range results{
		posts = append(posts, p.resultToPost(result))
	}

	return posts
}


func (p PostDAOImpl) updateFavoriteTable()error{
	favoritesRecords, err := p.getAllFavorites()
	if err != nil{
		return err
	}
	var userdao interfaces.UserDAO = userImpl.NewUserDAOImpl(p.db)
	
	for _, favoriteRecord := range favoritesRecords{
		favoriteUserId := int(favoriteRecord["user_id"].(int64))
		favoritePostId := int(favoriteRecord["post_id"].(int64))

		_, err := userdao.GetById(favoriteUserId)
		if err != nil{
			fmt.Printf("Invalid user id:%v This record was deleted.", favoriteUserId)
			p.DeleteFavorite(favoriteUserId, favoritePostId)
			continue
		}

		_, err = p.GetById(favoritePostId)
		if err != nil{
			fmt.Printf("Invalid post id:%v This record was deleted.", favoritePostId)
			p.DeleteFavorite(favoriteUserId, favoritePostId)
			continue
		}
	}

	return nil
}


func (p PostDAOImpl) ValidatePostField(post models.Post) error {
	title := post.GetTitle()
	description := post.GetDescription()

	if len(title) == 0 {
		return errors.New("title is required")
	}
	if len(title) > 100 {
		return errors.New("title must be less than 100 characters")
	}
	if description == "" {
		return errors.New("description is required")
	}
	if len(description) > 1000 {
		return errors.New("description must be less than 1000 characters")
	}
	return nil
}

