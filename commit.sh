set -e

echo "staging changes..."

git add .

read -p "Enter commit message: " message

git commit -m "$message"

git push

clear
