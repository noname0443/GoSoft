package main

import (
	"GoSoft/DBMS"
	"GoSoft/graph"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/lib/pq"
)

const defaultPort = "8080"

func graphqlHandler() gin.HandlerFunc {
	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "GoSoftToken", c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/api")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := gin.Default()
	router.LoadHTMLGlob("template/*")
	router.Use(GinContextToContextMiddleware())
	router.Static("/assets", "./assets")
	router.POST("/api", graphqlHandler())
	router.GET("/playground", playgroundHandler())
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	router.GET("/store", func(c *gin.Context) {
		c.HTML(http.StatusOK, "store.html", gin.H{})
	})
	router.GET("/profile", func(c *gin.Context) {
		token, err := c.Cookie("GoSoftToken")
		if err != nil {
			c.AbortWithStatus(400)
			return
		}
		if DBMS.ValidateToken(token) {
			c.HTML(http.StatusOK, "profile.html", gin.H{}) // TODO: Profile rendering
		} else {
			c.HTML(http.StatusOK, "register.html", gin.H{})
		}
	})
	router.GET("/store/:id", func(c *gin.Context) {
		productidString := c.Param("id")
		productid, err := strconv.ParseInt(productidString, 10, 32)
		if err != nil {
			c.AbortWithStatus(404)
		}
		product, err := DBMS.GetProduct(int(productid))
		if err != nil {
			c.AbortWithStatus(404)
			return
		} else {
			comments, err := DBMS.GetComments(int(productid))
			if err != nil {
				c.AbortWithStatus(400)
				return
			}
			token, err := c.Cookie("GoSoftToken")
			if err != nil {
				c.AbortWithStatus(400)
				return
			}
			c.HTML(http.StatusOK, "product.html", gin.H{
				"product":  product,
				"comments": comments,
				"loggedin": DBMS.ValidateToken(token),
			})
		}
	})
	router.Run()
}
