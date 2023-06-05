eval `ssh-agent -s`
ssh-add /root/.ssh/github_krassor_rsa
git pull git@github.com:krassor/skygrow.git
docker-compose up

# echo "stopping app-tg-gpt-bot container"
# docker stop app-tg-gpt-bot
# echo "removing app-tg-gpt-bot container"
# docker rm app-tg-gpt-bot
# echo "removing app-tg-gpt-bot image"
# docker rmi app-tg-gpt-bot
# echo "building new app-tg-gpt-bot image"
# docker build -t app-tg-gpt-bot .
# docker run -it -d --name app-tg-gpt-bot app-tg-gpt-bot