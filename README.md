DB 

User 
- username
- email
- password
- session
- conversations (UserID, ID, Title, echanges, dates)
- favoris (conversations ID)*

install necessaire pour la db

go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
go get -u github.com/joho/godotenv

commandes pour initialiser la DB

sudo -u postgres psql
CREATE USER mori WITH PASSWORD 'mori123456';
CREATE DATABASE mori;
GRANT ALL PRIVILEGES ON DATABASE mori TO mori;
\q
