CREATE DATABASE "GoSoft"
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    CONNECTION LIMIT = -1;

\c GoSoft

CREATE TABLE IF NOT EXISTS Users (
                                     UserID SERIAL PRIMARY KEY,
                                     Email TEXT NOT NULL UNIQUE,
                                     Name TEXT NOT NULL,
                                     Surname TEXT NOT NULL,
                                     Password TEXT NOT NULL,
                                     Gender TEXT NOT NULL,
                                     Token TEXT NOT NULL,
                                     TokenDate TIMESTAMP NOT NULL,
                                     RegistrationDate TIMESTAMP NOT NULL,
                                     Role TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS Store (
                                     ProductID SERIAL PRIMARY KEY,
                                     Name TEXT NOT NULL,
                                     Description TEXT NOT NULL,
                                     Photo TEXT NOT NULL,
                                     File TEXT NOT NULL,
                                     Price REAL NOT NULL,
                                     SubscriptionType TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS CART (
                                    UserID INTEGER REFERENCES Users(UserID) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
                                    ProductID INTEGER REFERENCES Store(ProductID) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
                                    Count INTEGER NOT NULL,
                                    CONSTRAINT unique_user_product UNIQUE (userid, productid),
                                    CONSTRAINT cartCountTrigger CHECK (Count > 0)
);

CREATE OR REPLACE FUNCTION cartClearer()
    RETURNS TRIGGER AS $$
BEGIN
    IF NEW.Count <= 0 THEN
        DELETE FROM CART WHERE UserID = NEW.UserID AND ProductID = NEW.ProductID;
        RETURN NULL;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER cartCountTrigger
    BEFORE INSERT OR UPDATE ON CART
    FOR EACH ROW
EXECUTE FUNCTION cartClearer();

CREATE TABLE IF NOT EXISTS Categories (
                                          ProductID INTEGER REFERENCES Store(ProductID) ON DELETE CASCADE ON UPDATE CASCADE,
                                          Category TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS Purchase (
                                        UserID INTEGER REFERENCES Users(UserID),
                                        DateTime TIMESTAMP NOT NULL,
                                        ProductID INTEGER REFERENCES Store(ProductID) ON DELETE CASCADE ON UPDATE CASCADE,
                                        Price REAL NOT NULL,
                                        OrderID TEXT NOT NULL,
                                        Paid BOOLEAN NOT NULL,
                                        Count INTEGER NOT NULL,
                                        SubscriptionType TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS Comment (
                                       CommentID SERIAL PRIMARY KEY,
                                       Date TIMESTAMP NOT NULL,
                                       UserID INTEGER REFERENCES Users(UserID) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
                                       ProductID INTEGER REFERENCES Store(ProductID) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
                                       Content TEXT NOT NULL
);