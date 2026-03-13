package commands

import(
	"errors"
	"fmt"
	"time"
	"context"
	"strings"
	"html"
	"io"
	"strconv"
	"net/http"
	"encoding/xml"
	"database/sql"
	"github.com/Dxnax-RS/gator/internal/config"
	"github.com/Dxnax-RS/gator/internal/database"

	"github.com/google/uuid"
)



func MiddlewareLoggedIn(handler func(s *config.State, cmd command, user database.User) error) func(*config.State, command) error{
	return func(s *config.State, cmd command) error {
		var err error
		user, err := s.Db.GetUser(context.Background(), s.Cfg.Current_user_name)

		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}

func HandlerLogin(s *config.State, cmd command) error{
	if len(cmd.Args) < 1{
		return errors.New("Username missing in arguments")
	}

	var err error

	
	if _, err := s.Db.GetUser(context.Background(), cmd.Args[0]); err != nil{
		fmt.Println("User not found")
		return err
	}

	err = s.Cfg.SetUser(cmd.Args[0]) 
	
	fmt.Println("Login successful")

	return err
}

func RegisterUser(s *config.State, cmd command) error{
	if len(cmd.Args) < 1{
		return errors.New("User missing in arguments")
	}

	var err error

	if _, err := s.Db.GetUser(context.Background(), cmd.Args[0]); err == nil{
		return errors.New("User already exists")
	}

	newUser := database.CreateUserParams{
		ID: 		uuid.New(),
		CreatedAt: 	time.Now(),
		UpdatedAt: 	time.Now(),
		Name: 		cmd.Args[0],
	}

	user, err := s.Db.CreateUser(context.Background(), newUser)

	if err != nil{
		return err
	}

	err = s.Cfg.SetUser(user.Name)

	if err != nil{
		return err
	}

	fmt.Println("User created successfully")
	fmt.Println(user)

	return err
}

func RegisterFeed(s *config.State, cmd command, user database.User) error{
	if len(cmd.Args) < 2{
		return errors.New("Missing arguments")
	}

	var err error

	newFeed := database.CreateFeedParams{
		ID: 		uuid.New(),
		CreatedAt: 	time.Now(),
		UpdatedAt: 	time.Now(),
		Name: 		cmd.Args[0],
		Url:		cmd.Args[1],
		UserID:		user.ID,
	}

	feed, err := s.Db.CreateFeed(context.Background(), newFeed)

	if err != nil{
		return err
	}

	newFollow := database.CreateFeedFollowParams{
		ID: 		uuid.New(),
		CreatedAt: 	time.Now(),
		UpdatedAt: 	time.Now(),
		UserID:		user.ID,
		FeedID:		feed.ID,
	}

	_, err = s.Db.CreateFeedFollow(context.Background(), newFollow)

	if err != nil {
		return err
	}

	fmt.Println("Feed created successfully")
	fmt.Println(feed)

	return err
}

func RegisterFollow(s *config.State, cmd command, user database.User) error{
	if len(cmd.Args) < 1{
		return errors.New("Missing arguments")
	}

	var err error
	
	feedId, err := s.Db.GetFeedIdByUrl(context.Background(), cmd.Args[0])

	if err != nil {
		return err
	}

	newFollow := database.CreateFeedFollowParams{
		ID: 		uuid.New(),
		CreatedAt: 	time.Now(),
		UpdatedAt: 	time.Now(),
		UserID:		user.ID,
		FeedID:		feedId,
	}

	follow, err := s.Db.CreateFeedFollow(context.Background(), newFollow)

	if err != nil {
		return err
	}

	fmt.Printf("The user %s is now following the feed: %s", follow.UserName, follow.FeedName)

	return err

}

func BrowsePosts(s *config.State, cmd command, user database.User) error{
	var err error
	var limit int
	if len(cmd.Args) < 1 {
		limit = 2
	} else {
		limit, err = strconv.Atoi(cmd.Args[0])
	}

	if err != nil {
		return err
	}

	newParams := database.GetPostsForUserParams{
		Limit: 		int32(limit),
		UserID: 	user.ID,
	}

	posts, err := s.Db.GetPostsForUser(context.Background(), newParams)

	if err != nil {
		return err
	}

	for i, v := range posts {
		fmt.Println("\nPost", i+1)
		fmt.Println("Title:", v.Title)
		fmt.Println("Description:", v.Description)
		fmt.Println("Url:", v.Url)
		fmt.Println("Published at:", v.PublishedAt)
		fmt.Printf("\n")
	}

	return err
}

func ResetUserTable(s *config.State, cmd command) error{
	err := s.Db.ResetUsers(context.Background())
	return err
}

func GetAllFeeds(s *config.State, cmd command) error{
	results, err := s.Db.GetFeeds(context.Background())

	if err != nil {
		return err
	}

	for _, value := range results {
		fmt.Printf("\n\nFeed name: %s	Feed url: %s\nFeed created by: %s", value.FeedName, value.Url, value.Name)
	}
	return err
}

func GetAllUsers(s *config.State, cmd command) error{
	users, err := s.Db.GetUsers(context.Background())
	
	if err != nil{
		return err
	}

	for _, v := range users {
		if v.Name == s.Cfg.Current_user_name{
			fmt.Println(v.Name, "(current)")
			continue

		}
		
		fmt.Println(v.Name)

	}

	return err
}

func GetUserFollows(s *config.State, cmd command, user database.User) error{
	follows, err := s.Db.GetFeedFollowsForUser(context.Background(), user.ID)

	if err != nil{
		return err
	}

	fmt.Println("The user", user.Name, "is following the feeds:")

	for _, v := range follows {
		fmt.Println(v.FeedName)
	}

	return err
}

func DeleteFollow(s *config.State, cmd command, user database.User) error{
	if len(cmd.Args) < 1{
		return errors.New("Missing arguments")
	}

	var err error
	
	feedId, err := s.Db.GetFeedIdByUrl(context.Background(), cmd.Args[0])

	if err != nil {
		return err
	}

	newDelete := database.DeleteFeedFollowParams{
		UserID:		user.ID,
		FeedID:		feedId,
	}

	err = s.Db.DeleteFeedFollow(context.Background(), newDelete)

	return err
}

func Aggregator(s *config.State, cmd command) error{
	timeBetweenRequests, err := time.ParseDuration("1m")

	if err != nil {
		return err
	}

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		err = scrapeFeeds(s)

		if err != nil {
			return err
		}
	}

	return nil
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error){
	
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "gator")

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var feed RSSFeed

	if err := xml.Unmarshal(data, &feed); err != nil {
		return nil, err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for _, v := range feed.Channel.Item {
		v.Title = html.UnescapeString(v.Title)
		v.Description = html.UnescapeString(v.Description)
	}

	return &feed, err
}

func scrapeFeeds(s *config.State) error{
	feed, err := s.Db.GetNextFeedToFetch(context.Background())

	if err != nil {
		return err
	}

	rss, err := fetchFeed(context.Background(), feed.Url)

	if err != nil {
		return err
	}

	lastFetched := sql.NullTime{
		Time:	time.Now(),
		Valid:	true,
	}

	newMark := database.MarkFeedFetchedParams{
		LastFetchedAt:	lastFetched,
		UpdatedAt: 		lastFetched.Time,
		ID:				feed.ID,
	}

	err = s.Db.MarkFeedFetched(context.Background(), newMark)

	if err != nil {
		return err
	}

	for _, v := range rss.Channel.Item {

		fmt.Println("Title: ", v.Title)
		fmt.Println("Description: ", v.Description)
		fmt.Println("Url: ", v.Link)
		fmt.Println("Date: ", v.PubDate)

		newTitle := sql.NullString{
			String: v.Title,
			Valid: true,
		}

		newDescription := sql.NullString{
			String: v.Description,
			Valid: true,
		}

		if v.Title == ""{
			newTitle.Valid = false
		}
		
		if v.Description == ""{
			newDescription.Valid = false
		}

		v.PubDate = strings.ReplaceAll(v.PubDate, " +0000", "")
		layout := "Mon, 02 Jan 2006 15:04:05"
		newTime, _ := time.Parse(layout, v.PubDate)

		newPost := database.CreatePostParams{
			ID: 			uuid.New(),
			CreatedAt: 		time.Now(),
			UpdatedAt: 		time.Now(),
			Title: 			newTitle,
			Url: 			v.Link,
			Description: 	newDescription,
			PublishedAt: 	newTime,
			FeedID: 		feed.ID,
		}
		
		_, err = s.Db.CreatePost(context.Background(), newPost)

		fmt.Println(err)
		if err != nil{
			if strings.Contains(err.Error(), "posts_url_key"){
				err = nil
			}
		}
	}

	return err
}