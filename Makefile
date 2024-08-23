# load env
include .env

ifeq ($(OS),Windows_NT)
    # Windows
    VOLUME_PATH_PLAYLISTS := C:\Users\guido\Desktop\playlists
    VOLUME_PATH_CONFIG := C:\Users\guido\Documents\heritage\config.json
else
    # Linux
    VOLUME_PATH_PLAYLISTS := /home/guido/Desktop/playlists
    VOLUME_PATH_CONFIG := /home/guido/Documents/heritage/config.json
endif

build:
	docker build -t heritage .

run: # this runs both windows and linux, do not edit
	docker run -it \
		-v $(VOLUME_PATH_PLAYLISTS):/app/playlists \
		-v $(VOLUME_PATH_CONFIG):/app/config.json \
	    -e SPOTIPY_CLIENT_ID=$(SPOTIPY_CLIENT_ID) \
        -e SPOTIPY_CLIENT_SECRET=$(SPOTIPY_CLIENT_SECRET) \
        -e SPOTIPY_REDIRECT_URI=$(SPOTIPY_REDIRECT_URI) \
        -e SPOTIPY_USERNAME=$(SPOTIPY_USERNAME) \
		heritage

clean: # this works
	docker kill heritage
	docker rm heritage


# REMEMBER THIS, WINDOWS IS SLOW.
clean-all:
	docker rmi -f $$(docker images -q);
	docker rm -f $$(docker ps -a -q)
