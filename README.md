# **E-commerce**

## **Getting Started**

### **Quick Start**

To start the project, simply run:

```bash
make all
```

### **Ports**
- Main Server: `8084`
- Backup Server: `8085`

## **API Endpoints**

### **User Authentication**

#### **Sign Up**
**POST** `/users/signup`

Request:
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "password": "securepassword",
  "phone": "+1234567890"
}
```
Response:
```json
"Successfully Signed Up!"
```

#### **Log In**
**POST** `/users/login`

Request:
```json
{
  "email": "john.doe@example.com",
  "password": "securepassword"
}
```
Response:
```json
{
  "_id": "unique_user_id",
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "token": "JWT_TOKEN",
  "refresh_token": "REFRESH_TOKEN",
  "created_at": "2025-01-12T08:00:00Z",
  "updated_at": "2025-01-12T08:00:00Z",
  "user_cart": [],
  "address": [],
  "orders": []
}
```

### **Admin Operations**

#### **Add Product**
**POST** `/admin/products/add`

Request:
```json
{
  "product_name": "MacBook Pro",
  "price": 1999,
  "rating": 4.5,
  "image": "MacBook_pro.jpg"
}
```
Response:
```json
"Product added successfully!"
```

### **Product Operations**

#### **View Products**
**GET** `/users/productview`

Response:
```json
[
  {
    "product_id": "12345",
    "product_name": "MacBook Pro",
    "price": 1999,
    "rating": 4.5,
    "image": "MacBook_pro.jpg"
  },
  {
    "product_id": "67890",
    "product_name": "SmartWidget",
    "price": 299,
    "rating": 4.7,
    "image": "smartwidget.jpg"
  }
]
```

#### **Search Products**
**GET** `/users/search?name=widget`

Response:
```json
[
  {
    "product_id": "67890",
    "product_name": "SmartWidget",
    "price": 299,
    "rating": 4.7,
    "image": "smartwidget.jpg"
  }
]
```

### **Cart Operations**

#### **Add to Cart**
**GET** `/cart/add?id=product_id&user_id=user_id`

Response: `"Product added to cart."`

#### **Remove from Cart**
**GET** `/cart/remove?id=product_id&user_id=user_id`

Response: `"Product removed from cart."`

#### **View Cart**
**GET** `/cart/list?user_id=user_id`

Response:
```json
{
  "cart_items": [
    {
      "product_id": "12345",
      "product_name": "MacBook Pro",
      "price": 1999,
      "quantity": 1
    }
  ],
  "total_price": 1999
}
```

#### **Checkout Cart**
**GET** `/cart/checkout?user_id=user_id`

Response: `"Order placed successfully!"`

#### **Instant Buy**
**GET** `/cart/buy?user_id=user_id&product_id=product_id`

Response: `"Purchase completed successfully!"`

### **Address Management**

#### **Add Address**
**POST** `/address/add`

Request:
```json
{
  "house_name": "Green Villa",
  "street_name": "Maple Street",
  "city_name": "Metropolis",
  "pin_code": "123456"
}
```
Response: `"Address added successfully!"`

#### **Edit Home Address**
**PUT** `/address/edit/home`

#### **Edit Work Address**
**PUT** `/address/edit/work`

#### **Delete Address**
**DELETE** `/address/delete`

Response: `"Address deleted successfully!"`

## **Technology Stack**

- **Programming Language**: Go (Golang)
- **Database**: MongoDB
- **Containerization**: Docker (Docker Compose)
- **API Documentation**: Swagger

## **Project Structure**

The project follows a modular structure with clearly defined routes for users, admin, products, cart, and address management.

## **Running the Project**

1. Clone the repository.
2. Ensure Docker is installed and running.
3. Use the `make all` command to build and run the project.
