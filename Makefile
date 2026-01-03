.PHONY: setup env up down build logs clean

.DEFAULT_GOAL := up

setup: env
	@echo "Setup complete! Run 'make up' to start the application."

env:
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		JWT_SECRET=$$(openssl rand -base64 32); \
		if [ "$$(uname)" = "Darwin" ]; then \
			sed -i '' "s|^JWT_SECRET=.*|JWT_SECRET=$$JWT_SECRET|" .env; \
		else \
			sed -i "s|^JWT_SECRET=.*|JWT_SECRET=$$JWT_SECRET|" .env; \
		fi; \
		@echo "Created .env with generated JWT_SECRET"; \
	fi

up:
	$(MAKE) env
	docker-compose up -d

build:
	docker-compose up -d --build

down:
	docker-compose down

logs:
	docker-compose logs -f

clean:
	docker-compose down -v
