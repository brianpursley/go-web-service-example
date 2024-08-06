package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	_ "web-service-gin/docs"
)

// @title           Example API
// @version         1.0
// @description     This is an API Example.
// @termsOfService  https://example.com

// @contact.name   API Support
// @contact.url    https://example.com
// @contact.email  support@example.com

// @license.name  Example License
// @license.url   https://example.com

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey  ApiKeyAuth
// @in header
// @name Authorization
// @type apiKey
// @description Use "key1" for read-only access or "key2" for full access.

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://example.com
func main() {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		albums := v1.Group("/albums")
		{
			albums.GET("", authenticate, getAlbums)
			albums.GET("/:id", authenticate, getAlbumByID)
			albums.POST("", authenticate, postAlbums)
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	_ = r.Run(":8080")
}

func authenticate(ctx *gin.Context) {
	if ctx.GetHeader("Authorization") == "key1" {
		ctx.Next()
		return
	}

	if ctx.GetHeader("Authorization") == "key2" {
		ctx.Set("Role", "Admin")
		ctx.Next()
		return
	}

	ctx.AbortWithStatusJSON(http.StatusUnauthorized, apiError{Error: "Unauthorized", Message: "Invalid API key"})
}

type apiError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

// getAlbums godoc
// @Summary      Gets all albums
// @Description  Gets all albums
// @Tags         albums
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {object}  album
// @Router       /albums [get]
func getAlbums(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, albums)
}

// getAlbumByID godoc
// @Summary      Gets an album by ID
// @Description  Gets an album by ID
// @Tags         albums
// @Security     ApiKeyAuth
// @Produce      json
// @Param        id  path  string  true  "Album ID"
// @Success      200  {object}  album
// @Failure      404  {object}  apiError
// @Router       /albums/{id} [get]
func getAlbumByID(ctx *gin.Context) {
	id := ctx.Param("id")

	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range albums {
		if a.ID == id {
			ctx.IndentedJSON(http.StatusOK, a)
			return
		}
	}

	ctx.AbortWithStatusJSON(http.StatusNotFound, apiError{Error: "Not Found", Message: "Album not found"})
}

// postAlbums godoc
// @Summary      Adds a new album
// @Description  Adds a new album
// @Tags         albums
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        album  body  album  true  "Album to add"
// @Success      201  {object}  album
// @Failure      400  {object}  apiError
// @Router       /albums [post]
func postAlbums(ctx *gin.Context) {
	var newAlbum album

	if role, _ := ctx.Get("Role"); role != "Admin" {
		ctx.AbortWithStatusJSON(http.StatusForbidden, apiError{Error: "Forbidden", Message: "Admin role required"})
		return
	}

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := ctx.BindJSON(&newAlbum); err != nil {
		return
	}

	// Check if an album with the same ID already exists.
	for _, a := range albums {
		if a.ID == newAlbum.ID {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, apiError{Error: "Bad Request", Message: "Album ID already exists"})
			return
		}
	}

	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	ctx.IndentedJSON(http.StatusCreated, newAlbum)
}
