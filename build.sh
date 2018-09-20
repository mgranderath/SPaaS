make build_linux
cd frontend
./node_modules/.bin/webpack
cd ..
docker build -t mgranderath/spaas:latest .