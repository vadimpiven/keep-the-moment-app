// This package implements router group `post`.
package post

import (
	"net/http"

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
		post.POST("/like-by-id", likePostByID, keyauth.Middleware())
		post.POST("/comment-by-id", commentPostByID, keyauth.Middleware())
		post.POST("/create", createPost, keyauth.Middleware())
	}
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
	posts, err := postgres.GetVisiblePostIDs(c, email)
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
		return echo.ErrBadRequest
	}

	if _, visible, err := postgres.CheckPostExists(c, in.ID); err != nil {
		return echo.ErrInternalServerError
	} else if visible != true {
		return echo.ErrBadRequest
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
// @Router /post/get-by-id [post]
func likePostByID(c echo.Context) error {
	in := new(getPostByIDIn)
	err := c.Bind(in)
	if err != nil || in.ID == 0 {
		return echo.ErrBadRequest
	}

	if _, visible, err := postgres.CheckPostExists(c, in.ID); err != nil {
		return echo.ErrInternalServerError
	} else if visible != true {
		return echo.ErrBadRequest
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
		return echo.ErrBadRequest
	}

	if _, visible, err := postgres.CheckPostExists(c, in.ID); err != nil {
		return echo.ErrInternalServerError
	} else if visible != true {
		return echo.ErrBadRequest
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
		return echo.ErrBadRequest
	}

	if exists, err := postgres.CheckImagesExist(c, in.Images); err != nil {
		return echo.ErrInternalServerError
	} else if exists != true {
		return echo.ErrBadRequest
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
