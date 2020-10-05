SHELL:=/bin/bash
IGNORE := $(shell bash -c "source env.sh; env | sed 's/=/:=/' | sed 's/^/export /' > makeenv")                         
include makeenv

keygen:
	./admin keygen ~/.avatarlysis

db-up:
	@docker-compose -f docker-compose.server.yml up -d postgres 
	

db-down:
	@docker-compose -f docker-compose.server.yml down 
	
db-log:
	@docker logs -f postgres

server:
	./app

web-up:
	@docker load -i frontend.tar
	@docker-compose -f docker-compose.web.yml up -d

web-down:
	@docker-compose -f docker-compose.web.yml down

clean: db-down
	@docker-compose rm
	@docker system prune --volumes