docker stop mw8_web_1
docker stop mw8_db_1
docker rmi mw8_web -f
docker rm mw8_web_1  -f
docker-compose up -d

