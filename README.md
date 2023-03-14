# funcaptcha-api
API takes form body

gamevariant - name of the game

file - mp3 challenge 

file is split into the 3 challenges then most likely index 1- 3 is returned

# Setup
ffmpeg is required but will be handled by docker.
clone the repo.

build docker container
```
sudo docker build -t funcaptcha .
```

run docker container
```
sudo docker run -p your-choice:8080 funcaptcha
```

