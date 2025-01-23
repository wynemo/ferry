用docker compose up -d 启动服务


docker exec -i ferry_mysql mysql -uferry -p123456 ferry < config/db.sql
docker exec -i ferry_mysql mysql -uferry -p123456 ferry < config/ferry.sql
