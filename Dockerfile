# python lightweight
FROM python:3.11.0-alpine

# set the working dir
WORKDIR /app
# copy everything
COPY . .

# install the requirements
RUN pip install -r requirements.txt

# install ffmpeg
RUN apk add ffmpeg

# Copy the compiled binary into the bin dir to use it directly
# https://github.com/yt-dlp/yt-dlp#release-files
COPY yt-dlp /usr/local/bin/

# Give permissions
RUN chmod 777 /usr/local/bin/yt-dlp

# Set spotify creds (TODO)
ENV SPOTIPY_CLIENT_ID=""
ENV SPOTIPY_CLIENT_SECRET=""
ENV SPOTIPY_REDIRECT_URI=""
ENV SPOTIPY_USERNAME=""

ENTRYPOINT ["python", "main.py"]
