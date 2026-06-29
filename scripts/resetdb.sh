set -e
echo "Removing all docker containers and volumes"
docker rm -f $(docker ps -aq)
docker volume rm $(docker volume ls -q)
read -p "Do you want to start containers? (y/n) " choice
if [ "$choice" != "y" ]; then
    echo "container initialization cancelled"
else
    echo "Starting docker compose"
    docker compose up -d
fi