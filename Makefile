check:
	encore check
	
dev:
	encore run --listen 0.0.0.0:4000


generate: 
	encore gen client -e staging --lang=typescript devweek-k65i > client/api.ts
	encore gen client -e staging --lang=go devweek-k65i > blogsync/client/goclient.go

.PHONY: blogsync
blogsync:
	cd blogsync && go build

sync:
	cd content && ../blogsync/blogsync push