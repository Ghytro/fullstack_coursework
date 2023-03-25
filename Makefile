start.local:
	cd cmd/github.com/Ghytro/galleryapp && go build -o app && cd ../../ && \
	DB_URL="postgres://postgres:mydbpassword@dockerdev.db:5432/postgres?sslmode=disable&" \
	./cmd/github.com/Ghytro/galleryapp/app

start.db:
	cd deployments && docker compose up -d db

up:
	cd deployments && docker compose up

up.detached:
	cd deployments && docker compose up -d

up.build:
	cd deployments && docker compose up --build

restart.nginx:
	docker restart myapp-nginx

down:
	cd deployments && docker compose down
