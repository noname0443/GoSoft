# GoSoft

## Installation Process
### Ubuntu
```
sudo apt-get install postgresql
sudo apt-get install go
mkdir GolangChat
cd GolangChat
git clone https://github.com/noname0443/GoSoft.git
psql -U myuser -f FirstRun.sql
go build .
go run server.go
```
### Windows
1. Install golang compiler on your device.
2. Install PostgreSQL.
3. ```git clone https://github.com/noname0443/GoSoft.git```
4. Launch psql.exe -U my_postgres_user -f FirstRun.sql
5. ```go build .```
6. ```go run server.go```

## Business Process
There are two actors - the user, GoSoft - the company selling application programs, companies developing software solutions.

GoSoft's goal is to make a profit from the sale of a license to use of their products and products of those companies that signed an agreement with them.
The purpose of the user is to use the license for a specific product/products for the implementation of the product/products. Target
solution companies to make a profit using GoSoft distributor.

To accomplish this goal, GoSoft provides two solutions. For users are provided with a website that allows them to find the necessary product, place an order, buy it and download it.
For companies an open API is provided, which, in particular, allows you to check whether the specified user has purchased their software solution.
Transfers between GoSoft and the companies are governed by private agreements and this the business process does not describe this interaction.
All funds received from sales are transferred primarily to GoSoft.

To achieve the goal, the user must register an account on the GoSoft website, thereby gaining access to the creation of orders and their payment.
Then create an order by selecting the necessary products, and complete payment on the shopping cart page via Paypal. After successful payment, all
ordered products will appear on the user's profile page and the will be able to successfully download them and apply them in practice.

Once a customer makes a payment, GoSoft receives the money on their
account and can fulfill the conditions stipulated in the relevant agreements
with software developers.

It is important to note that all active products on GoSoft must
be built-in protection that allows you to check the relevance of the user
subscriptions. Verification can be done using the GoSoft API by asking
user email and password.

## Technologies

The server technology stack is based primarily on the language Golang using gin framework to implement http server, which implements the project website.

To implement the API, the graphql standard is used, which was implemented by the gplgen library and used in all requests requiring information about the state of the object.

It is used to store information about the states of systems. relational PostgreSQL database.

In terms of user interaction with components site uses the built-in mechanisms of browsers to work with html/css/javascript. They are also used to send web requests to server side of the project.

## Roles
### Guest
This role describes the capabilities of a user who is not
registered and can perform the simplest actions:
- View assortment
- View comments
- Registration/authorization

![image](https://github.com/noname0443/GoSoft/blob/master/git/Guest%20Business%20Process%20Diagram.jpg)

### Customer
This role describes the capabilities of the user who
registered. He can perform the same actions as the guest, but his
wider possibilities
- View assortment
- View comments
- Profile view
- Add a comment
- Add item to cart
- View Cart
- Remove item from cart
- Buy items from the cart
- View purchased items
- Download purchased products

![image](https://github.com/noname0443/GoSoft/blob/master/git/Customer%20Business%20Process%20Diagram.jpg)

### Administrator
This role describes the capabilities of an administrator. The role is required for
filling the store with new products. Its features include
the possibilities of previous roles, but they are added:
- Adding a product to the store
- Removing a product
- Product information update

![image](https://github.com/noname0443/GoSoft/blob/master/git/Administrator%20Business%20Process%20Diagram.jpg)

## Database
The basis for storing information is a PostgreSQL relational database.

![image](https://github.com/noname0443/GoSoft/blob/master/git/Database%20Entity%20Relationship%20Diagram.jpg)

Users:
- UserID – serial number of the registered user
- Email - unique email address
- Name – username
- Surname – user's last name
- Password - password in SHA1 hashed form
- Gender - user's gender
- Token – session token that is used for authorized
action
- TokenDate – time of token creation, after 72 hours it will automatically
becomes unusable and the user must
re-enter your details.
- RegistrationDate - registration date.
- Role - user role, can be "customer" or "admin".

Store:
- ProductID - serial number of the product
- Name - product name
- Desciption – product description (May contain Russian letters, but
cannot contain line breaks)
- Photo - path to the photo on the server (should be square NxN)
- File – path to the file on the server, access to which is opened
automatically upon purchase of the product and is valid until
the subscription will not expire. It is required, as well as a photo, to upload to the server before
product creation.
- Price - the price in USD that the buyer must pay for the unit.
SubscriptionType. In other words, if the user subscribes to
month (SubscriptionType = month), then the price is multiplied by the quantity
months.
- SubscriptionType - subscription type: month, year.
- Company - the name of the manufacturing company.

Cart:
- UserID – serial number of the user-owner of the cart
- ProductID - serial number of the ordered product
- Count – quantity of the ordered product. Cannot be <= 0 otherwise
this order is automatically deleted.

Categories:
- ProductID - serial number of the product
- Category - the name of the category to which the product belongs.

Purchase:
- UserID – serial number of the user who bought/ordered
product.
- DateTime – order/purchase time.
- ProductID - the serial number of the purchased product.
- Price - price per unit of the purchased product.
- OrderID - order number, generated by the paypal service that responds
for creating and paying orders.
- Paid - order payment indicator
- Count - the number of SubscriptionType units during which
work subscription.
- SubscriptionType - subscription type: month, year.

Comment:
- CommentID – serial number of the comment
- Date - time when the comment was sent
- UserID – serial number of the user who sent the comment
- ProductID – serial number of the product to which the
a comment
- Content - the content of the comment
Trigger on the scheme is engaged in cleaning the database from records in Cart,
the quantity value of which is less than or equal to zero.

## Website scheme

![image](https://github.com/noname0443/GoSoft/blob/master/git/Website%20Site%20Map%20Diagram.jpg)
Above is a graph that describes the possible behavior of the client on
site (without transitions with an explicit URL in the address bar). Website
consists of 4 main parts:
1. Index/Main - the main page that contains information about the project.
2. Store - store page where you can search for products, view
detailed information and add comments to products.
3. Profile / register - a page that allows you to register / log in
profile to guests or view active subscriptions and information about
buyer profile.
4. Cart - a page that displays the current state of the cart
buyer.

## API Description
API (application program interface) - a formal description of the method
interaction with the server part of the software. In this case, it is used
graphql, all queries of which are represented as query structures.

```
enum Roles {
    Customer
    Administrator
}

type User {
    id: ID!
    email: String!
    name: String!
    surname: String!
    gender: String!
    date: String!
    role: Roles!
}

type Product {
    id: ID!
    name: String!
    description: String!
    photo: String!
    file: String!
    price: Float!
    subscriptiontype: String!
    company: String!
}

type Purchase {
    product: Product!
    count: Int!
    date: String!
}

type CartItem {
    product: Product!
    count: Int!
}

type Comment {
    id: ID!
    userid: ID!
    date: String!
    productid: ID!
    content: String!
}

type ExtendedComment {
    id: ID!
    name: String!
    surname: String!
    role: String!
    date: String!
    productid: ID!
    content: String!
}

input NewProduct {
    name: String!
    description: String!
    photo: String!
    file: String!
    price: Float!
    subscriptiontype: String!
    company: String!
}

type Query {
    search(name: String, categories: String, lower_price: Float, highest_price: Float): [Product!]!
    product(id: Int!): Product!
    comments(productid: Int!): [ExtendedComment!]!

    register(email: String!, name: String!, surname: String!, gender: String!, password: String!): Boolean!
    login(email: String!, password: String!): Boolean!

    CartAdd(productid: Int!, count: Int!): Boolean!
    CartRemove(productid: Int!, count: Int!): Boolean!
    CartInspect(productid: Int!): CartItem!
    CartGet: [CartItem!]
    CartPurchase: Boolean!

    ProfileGet: User!
    ProfileUpdate(email: String, name: String, surname: String, gender: String, password: String): User!

    History: [Purchase!]
    isBought(email: String!, password: String!, productid: Int!): Boolean!

    CommentAdd(content: String!, productid: Int!): Boolean!
    CommentRemove(commentid: Int!): Boolean!
    CommentUpdate(commentid: Int!, content: String!): Boolean!

    StoreAdd(product: NewProduct!): Boolean!
    StoreRemove(productid: Int!): Boolean!
    StoreUpdate(productid: Int!, product: NewProduct!): Boolean!
}
```
All structures are a human-readable description. Structure
request looks like this:
```
MyRequest(arg1: type1, arg2: type2, ...): result_type
```
What about types:

- ```type``` - regular type
- ```[type]``` - an array of this type
- ```type!``` - the exclamation mark says that the type is not null.
- ```[type!]``` - an array of non-null types (although the array itself can be null)
- ```[type]!``` - array of non-null types

## Licensing
Product licenses apply at the time they are first activated on
side of the user after purchasing them. GoSoft manufactures
licensing only your software and only under conditions from the past
offers.

All products purchased by the user are guaranteed to remain
available during the subscription period. Further extension
carried out either under a private contract with a third party -
by the manufacturer of the software solution, or automatically using
subscription systems on the GoSoft website.
