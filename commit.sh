set -e

echo "Staging changes"
git add .

read -p "Enter commit message: " message

git commit -m "$message"

echo "Pushing changes"
git push