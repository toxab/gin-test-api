# gin-test-api
example of using the gin framework with mysql8 via docker

# temporally if have issue with conection to db through vs code:
docker exec -it mysql-db bash
mysql -uroot -p
ALTER USER 'gouser'@'%' IDENTIFIED WITH mysql_native_password BY 'gopassword';
FLUSH PRIVILEGES;

docker-compose restart mysql


run: http://localhost:8054/