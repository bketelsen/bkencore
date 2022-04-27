check:
	encore check
	
dev:
	encore run --listen 0.0.0.0:4000


generate: 
	encore gen client -e local --lang=typescript devweek-k65i > client/api.ts
	encore gen client -e local --lang=go devweek-k65i > goclient.gox
	rm bkml/client/goclient.go
	mv goclient.gox bkml/client/goclient.go

.PHONY: bkml 
bkml:
	cd bkml && go install

sync: bkml
	cd content && bkml push -e staging

synclocal: bkml
	cd content && bkml push -e local