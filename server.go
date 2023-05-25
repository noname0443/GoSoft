package main

import (
	"GoSoft/graph"
	"context"
	"github.com/gin-gonic/gin"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/lib/pq"
)

const defaultPort = "8080"

func graphqlHandler() gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	return func(ctx *gin.Context) {
		ctx.Cookie("GoSoftToken")

		// Allow unauthenticated users in
		//if err != nil {
		//	h.ServeHTTP(ctx.Writer, ctx.Request)
		//	return
		//}
		cont := context.WithValue(ctx.Request.Context(), "GoSoftToken", "t")
		ctx.Request = ctx.Request.WithContext(cont)
		h.ServeHTTP(ctx.Writer, ctx.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

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
	//router.LoadHTMLGlob("static/*")
	router.Static("/resources", "./resources")
	router.POST("/query", graphqlHandler())
	router.GET("/", playgroundHandler())
	router.Run()
}
