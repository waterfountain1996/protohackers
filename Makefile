PROGRAMS := $(shell find ./problems -mindepth 1 -maxdepth 1 -type d | xargs -I {} basename {})

PORT ?= 10000

all: $(PROGRAMS)

%: problems/%
	@go build -o ./bin/$@ ./$<

tunnel:
	@ssh -g -N -R 0.0.0.0:$(PORT):127.0.0.1:$(PORT) $(SSH_HOST)
