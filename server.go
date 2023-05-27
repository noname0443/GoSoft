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
			c.HTML(http.StatusOK, "register.html", gin.H{})
			return
		}
		if DBMS.ValidateToken(token) {
			profile, err := DBMS.GetProfile(token)
			if err != nil {
				c.AbortWithStatus(404)
				return
			}
			history, err := DBMS.CartHistory(token)
			if err != nil {
				c.AbortWithStatus(404)
				return
			}
			for _, v := range history {
				v.Date = v.Date[0:10]
			}
			c.HTML(http.StatusOK, "profile.html", gin.H{
				"profile":   profile,
				"purchases": history,
			})
			return
		} else {
			c.HTML(http.StatusOK, "register.html", gin.H{})
			return
		}
	})
	router.GET("/files/:file", func(c *gin.Context) {
		// TODO: make file download if product is bought
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
	router.GET("/cart", func(c *gin.Context) {
		token, err := c.Cookie("GoSoftToken")
		if err != nil {
			c.AbortWithStatus(400)
			return
		}
		if DBMS.ValidateToken(token) {
			cart, err := DBMS.CartGet(token)
			if err != nil {
				return
			}
			var subtotal []float64
			sum := 0.0
			for _, v := range cart {
				price := v.Product.Price * float64(v.Count)
				sum += price
				subtotal = append(subtotal, price)
			}
			c.HTML(http.StatusOK, "cart.html", gin.H{
				"cart":     cart,
				"sum":      sum,
				"subtotal": subtotal,
			})
			return
		} else {
			c.HTML(http.StatusOK, "register.html", gin.H{})
			return
		}
	})
	router.Run()
}
