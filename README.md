# 💼 Wallet Project

🔧 Technologies & Libraries Used:
This project utilizes several well-known Go libraries and tools:
- **Ozzo Validation**: Input request validation.
- **Godotenv**: For loading environment variables.
- **jmoiron/sqlx**: MySQL driver.
- **Redis**: Caching with in-memory DB.
- **MySQL**: Relational Database.

# 📦 Setup Instructions (after cloning or unzipping):
```bash
cd wallet
go mod tidy
```
> **Note**: Make sure you update the `.env` file before running the project.

# 🗄️ Database:
- Inside the db folder, you’ll find a .sql file containing the table creation script.
- This project uses MySQL, so you can run the SQL file directly in your database management tool (e.g., phpMyAdmin, DBeaver, etc).

# 🧪 Running Unit Tests:
Unit tests are available in the usecase layer only.

To run tests:
```bash
cd path/to/usecase
go test
```
- If you're using VS Code, you can view the coverage:
  Open the test file → Right-click → Select "Go: Toggle Test Coverage in Current Package"

# ▶️ Running the Project:
Once your .env is properly configured with database and Redis details, simply run the project from the root directory:
- go run main.go

# 📮 API Collections:
You can import the provided Postman collection:
- Ewallet.postman_collection.json

Or try using curl with examples provided in the project folder.

## register
curl --location --request POST 'http://localhost:8080/api/user/create_user' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username":"usertest1"
}'

## login
curl --location --request POST 'http://localhost:8080/api/user/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username":"usertest1"
}'

## topup
curl --location --request POST 'http://localhost:8080/api/balance/balance_topup' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo2OCwid2FsbGV0X2lkIjo0MiwidXNlcm5hbWUiOiJ0ZXNhamExMTkxIiwiZXhwIjoxNzQ1NDAyOTg0fQ.Goi5ccbb358U04RhTt5707fUhtd6HAEWABcBG6PiyVA' \
--header 'Content-Type: application/json' \
--data-raw '{
    "amount":200000
}'

## Balance Read
curl --location --request GET 'http://localhost:8080/api/balance/balance_read' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo2OCwid2FsbGV0X2lkIjo0MiwidXNlcm5hbWUiOiJ0ZXNhamExMTkxIiwiZXhwIjoxNzQ1NDAyOTg0fQ.Goi5ccbb358U04RhTt5707fUhtd6HAEWABcBG6PiyVA' \
--header 'Content-Type: application/json' \
--data-raw '{
    "amount":10000
}'

## Transfer 
curl --location --request POST 'http://localhost:8080/api/transaction/transfer' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxNSwid2FsbGV0X2lkIjozLCJ1c2VybmFtZSI6ImpvZHlhbG1haWRhIiwiZXhwIjoxNzE2NzE2MzkxfQ.o0fgUyyQ46NkK7IJqa-nEgbsXXgse5OWCcYNoNPWoVk' \
--header 'Content-Type: application/json' \
--data-raw '{
    "to_username":"usertest1",
    "amount":50000
}'

## TopUser
curl --location --request GET 'http://localhost:8080/api/transaction/top_users' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo2OCwid2FsbGV0X2lkIjo0MiwidXNlcm5hbWUiOiJ0ZXNhamExMTkxIiwiZXhwIjoxNzQ1NDAyOTg0fQ.Goi5ccbb358U04RhTt5707fUhtd6HAEWABcBG6PiyVA'

## Top Transactions Per User
curl --location --request GET 'http://localhost:8080/api/transaction/top_transactions_per_user' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo2OCwid2FsbGV0X2lkIjo0MiwidXNlcm5hbWUiOiJ0ZXNhamExMTkxIiwiZXhwIjoxNzQ1NDAyOTg0fQ.Goi5ccbb358U04RhTt5707fUhtd6HAEWABcBG6PiyVA'
