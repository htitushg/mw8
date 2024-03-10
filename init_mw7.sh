docker stop mw7_web_1
docker stop mw7_db_1
docker rmi mw7_web -f
docker rm mw7_web_1  -f
docker-compose up -d

