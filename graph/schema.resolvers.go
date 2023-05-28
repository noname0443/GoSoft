package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.31

import (
	"GoSoft/DBMS"
	"GoSoft/Utility"
	"GoSoft/graph/model"
	"context"
	"errors"
	"math"
)

// Search is the resolver for the search field.
func (r *queryResolver) Search(ctx context.Context, name *string, categories *string, lowerPrice *float64, highestPrice *float64) ([]*model.Product, error) {
	var preparedName string
	var preparedCategory string
	var preparedLowerPrice float64
	var preparedHighestPrice float64

	if name == nil {
		preparedName = ""
	} else {
		preparedName = *name
	}

	if categories == nil {
		preparedCategory = ""
	} else {
		preparedCategory = *categories
	}

	if lowerPrice == nil {
		preparedLowerPrice = 0
	} else {
		preparedLowerPrice = *lowerPrice
	}

	if highestPrice == nil {
		preparedHighestPrice = math.Inf(0)
	} else {
		preparedHighestPrice = *highestPrice
	}

	products, err := DBMS.SearchProducts(preparedName, preparedCategory, preparedLowerPrice, preparedHighestPrice)
	if err != nil {
		return nil, err
	}
	return products, nil
}

// Product is the resolver for the product field.
func (r *queryResolver) Product(ctx context.Context, id int) (*model.Product, error) {
	product, err := DBMS.GetProduct(id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

// Comments is the resolver for the comments field.
func (r *queryResolver) Comments(ctx context.Context, productid int) ([]*model.ExtendedComment, error) {
	comments, err := DBMS.GetComments(productid)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// Register is the resolver for the register field.
func (r *queryResolver) Register(ctx context.Context, email string, name string, surname string, gender string, password string) (bool, error) {
	gc, err := Utility.GinContextFromContext(ctx)
	if err != nil {
		return false, err
	}
	token, err := DBMS.RegisterCustomer(email, name, surname, gender, password)
	if err != nil || len(token) == 0 {
		return false, err
	}
	gc.SetCookie("GoSoftToken", token, 259200, "", "", true, true)
	return true, nil
}

// Login is the resolver for the login field.
func (r *queryResolver) Login(ctx context.Context, email string, password string) (bool, error) {
	gc, err := Utility.GinContextFromContext(ctx)
	if err != nil {
		return false, err
	}
	token, err := DBMS.LoginCustomer(email, password)
	if err != nil || len(token) == 0 {
		return false, err
	}
	gc.SetCookie("GoSoftToken", token, 259200, "", "", true, true)
	return true, nil
}

// CartAdd is the resolver for the CartAdd field.
func (r *queryResolver) CartAdd(ctx context.Context, productid int, count int) (bool, error) {
	gc, err := Utility.GinContextFromContext(ctx)
	if err != nil {
		return false, err
	}
	token, err := gc.Cookie("GoSoftToken")
	if err != nil {
		return false, err
	}
	if !DBMS.ValidateToken(token) {
		return false, errors.New("you have to log in")
	}
	err = DBMS.CartAdd(token, productid, count)
	if err != nil {
		return false, err
	}
	return true, nil
}

// CartRemove is the resolver for the CartRemove field.
func (r *queryResolver) CartRemove(ctx context.Context, productid int, count int) (bool, error) {
	gc, err := Utility.GinContextFromContext(ctx)
	if err != nil {
		return false, err
	}
	token, err := gc.Cookie("GoSoftToken")
	if err != nil {
		return false, err
	}
	if !DBMS.ValidateToken(token) {
		return false, errors.New("you have to log in")
	}
	err = DBMS.CartRemove(token, productid, count)
	if err != nil {
		return false, err
	}
	return true, nil
}

// CartInspect is the resolver for the CartInspect field.
func (r *queryResolver) CartInspect(ctx context.Context, productid int) (*model.CartItem, error) {
	gc, err := Utility.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}
	token, err := gc.Cookie("GoSoftToken")
	if err != nil {
		return nil, err
	}
	if !DBMS.ValidateToken(token) {
		return nil, errors.New("you have to log in")
	}
	item, err := DBMS.CartGetItem(token, productid)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// CartGet is the resolver for the CartGet field.
func (r *queryResolver) CartGet(ctx context.Context) ([]*model.CartItem, error) {
	gc, err := Utility.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}
	token, err := gc.Cookie("GoSoftToken")
	if err != nil {
		return nil, err
	}
	if !DBMS.ValidateToken(token) {
		return nil, errors.New("you have to log in")
	}
	cart, err := DBMS.CartGet(token)
	if err != nil {
		return nil, err
	}
	return cart, nil
}

// CartPurchase is the resolver for the CartPurchase field.
func (r *queryResolver) CartPurchase(ctx context.Context) (bool, error) {
	gc, err := Utility.GinContextFromContext(ctx)
	if err != nil {
		return false, err
	}
	token, err := gc.Cookie("GoSoftToken")
	if err != nil {
		return false, err
	}
	if !DBMS.ValidateToken(token) {
		return false, errors.New("you have to log in")
	}
	err = DBMS.CartPurchase(token, "1")
	if err != nil {
		return false, err
	}
	return true, nil
}

// ProfileGet is the resolver for the ProfileGet field.
func (r *queryResolver) ProfileGet(ctx context.Context) (*model.User, error) {
	gc, err := Utility.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}
	token, err := gc.Cookie("GoSoftToken")
	if err != nil {
		return nil, err
	}
	if !DBMS.ValidateToken(token) {
		return nil, errors.New("you have to log in")
	}
	profile, err := DBMS.GetProfile(token)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

// ProfileUpdate is the resolver for the ProfileUpdate field.
func (r *queryResolver) ProfileUpdate(ctx context.Context, email *string, name *string, surname *string, gender *string, password *string) (*model.User, error) {
	gc, err := Utility.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}
	token, err := gc.Cookie("GoSoftToken")
	if err != nil {
		return nil, err
	}

	var preparedEmail string
	var preparedName string
	var preparedSurname string
	var preparedGender string
	var preparedPassword string

	if email == nil {
		preparedEmail = ""
	} else {
		preparedEmail = *email
	}

	if name == nil {
		preparedName = ""
	} else {
		preparedName = *name
	}

	if surname == nil {
		preparedSurname = ""
	} else {
		preparedSurname = *surname
	}

	if gender == nil {
		preparedGender = ""
	} else {
		preparedGender = *gender
	}

	if password == nil {
		preparedPassword = ""
	} else {
		preparedPassword = *password
	}

	if !DBMS.ValidateToken(token) {
		return nil, errors.New("you have to log in")
	}
	err = DBMS.UpdateProfile(token, preparedEmail, preparedName, preparedSurname, preparedGender, preparedPassword)
	if err != nil {
		return nil, err
	}
	profile, err := DBMS.GetProfile(token)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

// History is the resolver for the History field.
func (r *queryResolver) History(ctx context.Context) ([]*model.Purchase, error) {
	gc, err := Utility.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}
	token, err := gc.Cookie("GoSoftToken")
	if err != nil {
		return nil, err
	}
	if !DBMS.ValidateToken(token) {
		return nil, errors.New("you have to log in")
	}
	history, err := DBMS.PurchasedSoftware(token)
	if err != nil {
		return nil, err
	}
	return history, nil
}

// IsBought is the resolver for the isBought field.
func (r *queryResolver) IsBought(ctx context.Context, email string, password string, productid int) (bool, error) {
	result, err := DBMS.CheckUsersProduct(email, password, productid)
	if err != nil {
		return false, err
	}
	return result, nil
}

// CommentAdd is the resolver for the CommentAdd field.
func (r *queryResolver) CommentAdd(ctx context.Context, content string, productid int) (bool, error) {
	gc, err := Utility.GinContextFromContext(ctx)
	if err != nil {
		return false, err
	}
	token, err := gc.Cookie("GoSoftToken")
	if err != nil {
		return false, err
	}
	if !DBMS.ValidateToken(token) {
		return false, errors.New("you have to log in")
	}
	err = DBMS.AddComment(token, productid, content)
	if err != nil {
		return false, err
	}
	return true, nil
}

// CommentRemove is the resolver for the CommentRemove field.
func (r *queryResolver) CommentRemove(ctx context.Context, commentid int) (bool, error) {
	gc, err := Utility.GinContextFromContext(ctx)
	if err != nil {
		return false, err
	}
	token, err := gc.Cookie("GoSoftToken")
	if err != nil {
		return false, err
	}
	if !DBMS.ValidateToken(token) {
		return false, errors.New("you have to log in")
	}
	err = DBMS.RemoveComment(token, commentid)
	if err != nil {
		return false, err
	}
	return true, nil
}

// CommentUpdate is the resolver for the CommentUpdate field.
func (r *queryResolver) CommentUpdate(ctx context.Context, commentid int, content string) (bool, error) {
	gc, err := Utility.GinContextFromContext(ctx)
	if err != nil {
		return false, err
	}
	token, err := gc.Cookie("GoSoftToken")
	if err != nil {
		return false, err
	}
	if !DBMS.ValidateToken(token) {
		return false, errors.New("you have to log in")
	}
	err = DBMS.UpdateComment(token, commentid, content)
	if err != nil {
		return false, err
	}
	return true, nil
}

// StoreAdd is the resolver for the StoreAdd field.
func (r *queryResolver) StoreAdd(ctx context.Context, product model.NewProduct) (bool, error) {
	gc, err := Utility.GinContextFromContext(ctx)
	if err != nil {
		return false, err
	}
	token, err := gc.Cookie("GoSoftToken")
	if err != nil {
		return false, err
	}
	if !DBMS.ValidatePrivileges(token, "admin") {
		return false, errors.New("you have to log in")
	}
	err = DBMS.StoreAdd(product)
	if err != nil {
		return false, err
	}
	return true, nil
}

// StoreRemove is the resolver for the StoreRemove field.
func (r *queryResolver) StoreRemove(ctx context.Context, productid int) (bool, error) {
	gc, err := Utility.GinContextFromContext(ctx)
	if err != nil {
		return false, err
	}
	token, err := gc.Cookie("GoSoftToken")
	if err != nil {
		return false, err
	}
	if !DBMS.ValidatePrivileges(token, "admin") {
		return false, errors.New("you have to log in")
	}
	err = DBMS.StoreRemove(productid)
	if err != nil {
		return false, err
	}
	return true, nil
}

// StoreUpdate is the resolver for the StoreUpdate field.
func (r *queryResolver) StoreUpdate(ctx context.Context, productid int, product model.NewProduct) (bool, error) {
	gc, err := Utility.GinContextFromContext(ctx)
	if err != nil {
		return false, err
	}
	token, err := gc.Cookie("GoSoftToken")
	if err != nil {
		return false, err
	}
	if !DBMS.ValidatePrivileges(token, "admin") {
		return false, errors.New("you have to log in")
	}
	err = DBMS.StoreUpdate(productid, product)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
