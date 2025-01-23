注意 docker/entrypoint.sh 要为linux格式，否则会报错， 不能是windows格式
可以用note++ 查看并转换格式

用docker compose up -d 启动服务


docker exec -i ferry_mysql mysql -uferry -p123456 ferry < config/db.sql
docker exec -i ferry_mysql mysql -uferry -p123456 ferry < config/ferry.sql

数据库在mysql/db 文件夹
前端目录在ferry_web

如果修改了代码，需要重新构建镜像：
docker compose down && docker rmi ferry:latest || true && docker compose up
