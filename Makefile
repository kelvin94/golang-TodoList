appname=golang-todolist
appnameProduction=golang-todolist-production

dockerProductionFile=Dockerfile.production
dockerImageTag=golang-todolist

dockerImageProductionTag=golang-todolist-production
# build: ## Build the container
# 	docker build -t $(dockerImageTag) .
build-production:
	docker build -t $(dockerImageProductionTag) -f $(dockerProductionFile) .
# run:
# 	docker run --name $(appname) -p 3000:8080 golang-todolist:latest 
run-production:
	docker run --name $(appnameProduction) -p 3000:8080 $(dockerImageProductionTag) 
# stop:
# 	docker stop $(appname); docker rm $(appname)
stop-production:
	docker stop $(appnameProduction); docker rm $(appnameProduction)