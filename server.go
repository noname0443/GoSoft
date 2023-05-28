package main

import (
	"GoSoft/DBMS"
	"GoSoft/graph"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/lib/pq"
	"github.com/plutov/paypal/v4"
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

	paypal_handler, err := paypal.NewClient("AYZy3f-tqV9_YZOIavMunW1q4RXxrWVZLaeuHk2sX7BIzS7wHfEB3pLZmZtaqv4bDTM5XysCytmeGaNg",
		"EH6ImKiJycjZ-UKQE2C9-6ULvMj2b_DKDYxPGI-AFMG02N5Gs0nyAXFXZY3KJz3xAZvbZ5CL3jQOYOLH", paypal.APIBaseSandBox)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = paypal_handler.GetAccessToken(context.Background())
	if err != nil {
		log.Fatalln(err)
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
			history, err := DBMS.PurchasedSoftware(token)
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
	router.GET("/admin", func(c *gin.Context) {
		token, err := c.Cookie("GoSoftToken")
		if err != nil {
			c.AbortWithStatus(400)
			return
		}
		if DBMS.ValidatePrivileges(token, "admin") {
			products, err := DBMS.SearchProducts("", "", 0, 1e32)
			if err != nil {
				return
			}
			c.HTML(http.StatusOK, "admin.html", gin.H{
				"products": products,
			})
			return
		} else {
			c.AbortWithStatus(403)
			return
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
			var productids []int
			sum := 0.0
			for _, v := range cart {
				price := v.Product.Price * float64(v.Count)
				sum += price
				subtotal = append(subtotal, price)

				parseInt, err := strconv.ParseInt(v.Product.ID, 10, 32)
				if err != nil {
					return
				}
				productids = append(productids, int(parseInt))
			}
			c.HTML(http.StatusOK, "cart.html", gin.H{
				"cart":       cart,
				"sum":        sum,
				"subtotal":   subtotal,
				"productids": productids,
			})
			return
		} else {
			c.HTML(http.StatusOK, "register.html", gin.H{})
			return
		}
	})
	router.POST("/paypal/create-paypal-order", func(c *gin.Context) {
		token, err := c.Cookie("GoSoftToken")
		if err != nil {
			c.AbortWithStatus(400)
			return
		}
		if !DBMS.ValidateToken(token) {
			c.AbortWithStatus(400)
			return
		}
		cart, err := DBMS.CartGet(token)
		if err != nil {
			c.AbortWithStatus(400)
			return
		}

		sum := 0.0
		for _, v := range cart {
			price := v.Product.Price * float64(v.Count)
			sum += price
		}
		order, err := paypal_handler.CreateOrder(c, "CAPTURE",
			[]paypal.PurchaseUnitRequest{{
				Amount: &paypal.PurchaseUnitAmount{
					Currency:  "USD",
					Value:     strconv.FormatFloat(sum, 'f', -1, 64),
					Breakdown: nil,
				},
			},
			}, nil, nil)
		if err != nil {
			c.AbortWithStatus(400)
			return
		}

		err = DBMS.CartPurchase(token, order.ID)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(400)
			return
		}

		c.JSON(200, order)
	})
	router.POST("/paypal/capture-paypal-order", func(c *gin.Context) {
		token, err := c.Cookie("GoSoftToken")
		if err != nil {
			c.AbortWithStatus(400)
			return
		}
		type Order struct {
			OrderID string `json:"orderID"`
		}
		var order Order
		err = c.Bind(&order)
		if err != nil {
			return
		}
		getOrder, err := paypal_handler.GetOrder(c, order.OrderID)
		if err != nil {
			return
		}
		log.Println(getOrder)
		log.Println(getOrder.Status)
		if getOrder.Status != "APPROVED" {
			c.AbortWithStatus(400)
		}
		captureOrder, err := paypal_handler.CaptureOrder(c, order.OrderID, paypal.CaptureOrderRequest{})
		if err != nil {
			c.AbortWithStatus(400)
			return
		}
		err = DBMS.CartMakePaid(token, captureOrder.ID)
		if err != nil {
			c.AbortWithStatus(400)
			return
		}

		c.Status(200)
	})
	router.Run()
}
