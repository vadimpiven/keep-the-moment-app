// This package implements router group `post`.
package post

import (
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"

	"github.com/FTi130/keep-the-moment-app/back/lib/keyauth"
	"github.com/FTi130/keep-the-moment-app/back/lib/postgres"
)

// ApplyRoutes applies routes for the router group.
func ApplyRoutes(g *echo.Group) {
	post := g.Group("/post")
	{
		post.GET("/visible", getVisiblePosts)
		post.POST("/get-by-id", getPostByID)
		post.POST("/get-by-userid", getPostByUserID)
		post.POST("/get-by-hashtag", getPostByHashtag)
		post.POST("/like-by-id", likePostByID, keyauth.Middleware())
		post.POST("/comment-by-id", commentPostByID, keyauth.Middleware())
		post.POST("/create", createPost, keyauth.Middleware())
		post.POST("/mine", getMinePosts, keyauth.Middleware())
	}
}

type (
	getMinePostsIn struct {
		Page int `json:"page"`
	}
	getMinePostsOut200 struct {
		Page  int                      `json:"page"`
		Posts []postgres.PostAssembled `json:"posts"`
	}
)

// @Summary Returns posts created by the user
// @Accept json
// @Produce json
// @Param page body getMinePostsIn true "wrapped page"
// @Success 200 {object} getMinePostsOut200
// @Failure 400,500 {object} httputil.HTTPError
// @Router /post/mine [post]
func getMinePosts(c echo.Context) error {
	in := new(getMinePostsIn)
	err := c.Bind(in)
	if err != nil || in.Page < 0 {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "input structure not followed",
			Internal: err,
		}
	}

	email, err := keyauth.GetEmail(c)
	if err != nil {
		return echo.ErrInternalServerError
	}

	ids, err := postgres.GetPostsMadeByUser(c, email, in.Page)
	if err != nil {
		return echo.ErrInternalServerError
	}

	posts := make([]postgres.PostAssembled, 0, len(ids))
	for _, id := range ids {
		post, err := postgres.GetPostByID(c, id, email)
		if err != nil {
			return echo.ErrInternalServerError
		}
		posts = append(posts, post)
	}

	return c.JSON(http.StatusOK, getMinePostsOut200{in.Page, posts})
}

type (
	getPostByUserIDIn struct {
		UserID string `json:"user_id"`
	}
	getPostByUserIDOut200 struct {
		Posts []postgres.PostBrief `json:"posts"`
	}
)

// @Summary Returns visible posts made or commented by user with given userID
// @Accept json
// @Produce json
// @Param id body getPostByUserIDIn true "wrapped userId"
// @Success 200 {object} getPostByUserIDOut200
// @Failure 400,500 {object} httputil.HTTPError
// @Router /post/get-by-userid [post]
func getPostByUserID(c echo.Context) error {
	in := new(getPostByUserIDIn)
	err := c.Bind(in)
	if err != nil || in.UserID == "" {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "input structure not followed",
			Internal: err,
		}
	}

	re := regexp.MustCompile("[^a-z0-9_]+")
	if re.Find([]byte(in.UserID)) != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "user_id has the wrong format",
			Internal: err,
		}
	}

	posts, err := postgres.GetPostBriefsByUserID(c, in.UserID)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, getPostByUserIDOut200{Posts: posts})
}

type (
	getPostByHashtagIn struct {
		Hashtag string `json:"hashtag"`
	}
	getPostByHashtagOut200 struct {
		Posts []postgres.PostBrief `json:"posts"`
	}
)

// @Summary Returns visible posts containing hashtag in post or post author account
// @Accept json
// @Produce json
// @Param id body getPostByUserIDIn true "wrapped userId"
// @Success 200 {object} getPostByUserIDOut200
// @Failure 400,500 {object} httputil.HTTPError
// @Router /post/get-by-hashtag [post]
func getPostByHashtag(c echo.Context) error {
	in := new(getPostByHashtagIn)
	err := c.Bind(in)
	if err != nil || in.Hashtag == "" {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "input structure not followed",
			Internal: err,
		}
	}

	re := regexp.MustCompile("[^a-z0-9_]+")
	if re.Find([]byte(in.Hashtag)) != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "one of hashtags has the wrong format",
			Internal: err,
		}
	}

	posts, err := postgres.GetPostBriefsByHashtag(c, in.Hashtag)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, getPostByUserIDOut200{Posts: posts})
}

type (
	getVisiblePostsOut200 struct {
		Posts []postgres.PostBrief `json:"posts"`
	}
)

// @Summary Returns visible posts
// @Produce json
// @Success 200 {object} getVisiblePostsOut200
// @Failure 400,500 {object} httputil.HTTPError
// @Router /post/visible [get]
func getVisiblePosts(c echo.Context) error {
	email, _ := keyauth.GetEmail(c)
	posts, err := postgres.GetVisiblePostBriefs(c, email)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, getVisiblePostsOut200{Posts: posts})
}

type (
	getPostByIDIn struct {
		ID uint64 `json:"id"`
	}
	getPostByIDOut200 postgres.PostAssembled
)

// @Summary Returns existing posts
// @Accept json
// @Produce json
// @Param id body getPostByIDIn true "wrapped id"
// @Success 200 {object} getPostByIDOut200
// @Failure 400,500 {object} httputil.HTTPError
// @Router /post/get-by-id [post]
func getPostByID(c echo.Context) error {
	in := new(getPostByIDIn)
	err := c.Bind(in)
	if err != nil || in.ID == 0 {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "input structure not followed",
			Internal: err,
		}
	}

	if _, visible, err := postgres.CheckPostExists(c, in.ID); err != nil {
		return echo.ErrInternalServerError
	} else if visible != true {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "post with given id is hidden",
			Internal: err,
		}
	}

	email, _ := keyauth.GetEmail(c)
	post, err := postgres.GetPostByID(c, in.ID, email)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, post)
}

type (
	likePostByIDIn struct {
		ID uint64 `json:"id"`
	}
	likePostByIDOut200 postgres.PostAssembled
)

// @Summary Toggle like on post
// @Accept json
// @Produce json
// @Security Bearer
// @Param id body likePostByIDIn true "wrapped id"
// @Success 200 {object} likePostByIDOut200
// @Failure 400,401,500 {object} httputil.HTTPError
// @Router /post/like-by-id [post]
func likePostByID(c echo.Context) error {
	in := new(getPostByIDIn)
	err := c.Bind(in)
	if err != nil || in.ID == 0 {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "input structure not followed",
			Internal: err,
		}
	}

	if _, visible, err := postgres.CheckPostExists(c, in.ID); err != nil {
		return echo.ErrInternalServerError
	} else if visible != true {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "post with given id is hidden",
			Internal: err,
		}
	}

	email, err := keyauth.GetEmail(c)
	if err != nil {
		return echo.ErrInternalServerError
	}

	err = postgres.LikePostByID(c, in.ID, email)
	if err != nil {
		return echo.ErrInternalServerError
	}

	post, err := postgres.GetPostByID(c, in.ID, email)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, post)
}

type (
	commentPostByIDIn struct {
		ID      uint64 `json:"id"`
		Comment string `json:"comment"`
	}
	commentPostByIDOut200 postgres.PostAssembled
)

// @Summary Add comment to post
// @Accept json
// @Produce json
// @Security Bearer
// @Param id body commentPostByIDIn true "post id and comment"
// @Success 200 {object} commentPostByIDOut200
// @Failure 400,401,500 {object} httputil.HTTPError
// @Router /post/comment-by-id [post]
func commentPostByID(c echo.Context) error {
	in := new(commentPostByIDIn)
	err := c.Bind(in)
	if err != nil || in.ID == 0 || in.Comment == "" {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "input structure not followed",
			Internal: err,
		}
	}

	if _, visible, err := postgres.CheckPostExists(c, in.ID); err != nil {
		return echo.ErrInternalServerError
	} else if visible != true {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "post with given id is hidden",
			Internal: err,
		}
	}

	email, err := keyauth.GetEmail(c)
	if err != nil {
		return echo.ErrInternalServerError
	}

	err = postgres.CommentPostByID(c, in.ID, email, in.Comment)
	if err != nil {
		return echo.ErrInternalServerError
	}

	post, err := postgres.GetPostByID(c, in.ID, email)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, post)
}

type (
	createPostIn struct {
		Background []int32  `json:"background"`
		Content    string   `json:"content"`
		Hashtags   []string `json:"hashtags"`
		Images     []string `json:"images"`
	}
	createPostOut200 postgres.PostAssembled
)

// @Summary Creates new post.
// @Security Bearer
// @Accept json
// @Produce json
// @Security Bearer
// @Param post_data body createPostIn true "post content"
// @Success 200 {object} createPostOut200
// @Failure 400,401,500 {object} httputil.HTTPError
// @Router /post/create [post]
func createPost(c echo.Context) error {
	in := new(createPostIn)
	err := c.Bind(in)
	if in.Background == nil {
		in.Background = []int32{}
	}
	if in.Hashtags == nil {
		in.Hashtags = []string{}
	}
	if err != nil || len(in.Images) > 5 ||
		(len(in.Images) == 0 && len(in.Hashtags) == 0 && in.Content == "") {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "input structure not followed",
			Internal: err,
		}
	}

	re := regexp.MustCompile("[^a-z0-9_]+")
	for _, hashtag := range in.Hashtags {
		if re.Find([]byte(hashtag)) != nil {
			return &echo.HTTPError{
				Code:     http.StatusBadRequest,
				Message:  "one of hashtags has the wrong format",
				Internal: err,
			}
		}
	}

	if exists, err := postgres.CheckImagesExist(c, in.Images); err != nil {
		return echo.ErrInternalServerError
	} else if exists != true {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "given image not exists",
			Internal: err,
		}
	}

	email, err := keyauth.GetEmail(c)
	if err != nil {
		return echo.ErrInternalServerError
	}

	data := postgres.Post{
		Email:      email,
		Background: in.Background,
		Content:    in.Content,
		Hashtags:   in.Hashtags,
	}
	if len(in.Images) > 0 {
		data.Image1 = in.Images[0]
	}
	if len(in.Images) > 1 {
		data.Image2 = in.Images[1]
	}
	if len(in.Images) > 2 {
		data.Image3 = in.Images[2]
	}
	if len(in.Images) > 3 {
		data.Image4 = in.Images[3]
	}
	if len(in.Images) > 4 {
		data.Image5 = in.Images[4]
	}

	if err = postgres.CreateNewPost(c, &data); err != nil {
		return echo.ErrInternalServerError
	}
	return c.JSON(http.StatusOK, postgres.PostAssembled{
		Post:     data,
		Comments: []postgres.PostComment{},
		IsLiked:  false,
	})
}
