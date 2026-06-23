set -e
echo "Removing all docker containers and volumes"
docker rm -f $(docker ps -aq)
docker volume rm $(docker volume ls -q)
echo "Starting docker compose"
docker compose up