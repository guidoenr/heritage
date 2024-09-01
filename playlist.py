from util import format_track_name
from spotipy import SpotifyOAuth
from track import Track
from logger import get_custom_logger
from pytube import Playlist as YoutubePlaylist
from pytube import YouTube
from enum import Enum
from os.path import splitext

import os
import spotipy
import threading

# to read
scopes = ["playlist-read-collaborative", "playlist-read-private", "playlist-modify-private", "playlist-modify-public"]
fields = 'items(track(uri,name,artists(name))),next'
logger = get_custom_logger()

# Spotify Auth
sp_oauth = SpotifyOAuth(client_id=os.environ["SPOTIPY_CLIENT_ID"],
                        client_secret=os.environ["SPOTIPY_CLIENT_SECRET"],
                        redirect_uri=os.environ["SPOTIPY_REDIRECT_URI"],
                        username=os.environ["SPOTIPY_USERNAME"],
                        scope=scopes,
                        open_browser=False)


def authorize():
    token_info = sp_oauth.get_cached_token()
    if not token_info:
        # Prompt the user to authenticate and authorize the application
        auth_url = sp_oauth.get_authorize_url()
        print(f"Please visit this URL to authorize the application: {auth_url}")

        # Wait for the user to complete the authentication process
        response = input("Enter the URL you were redirected to: ")
        code = sp_oauth.parse_response_code(response)

        # Exchange the authorization code for a token
        token_info = sp_oauth.get_access_token(code)
        if token_info:
            return token_info['access_token']
        else:
            return None


# get access_token
access_token = authorize()

# Spotify Client
spotify_client = spotipy.Spotify(auth_manager=sp_oauth, auth=access_token)


class PlaylistType(Enum):
    SPOTIFY = "spotify"
    YOUTUBE = "youtube"


class Playlist:
    def __init__(self, playlist_id, owner):
        self.owner = owner
        self.playlist_id = playlist_id
        self.tracks_local = []
        self.tracks = []

        if playlist_id.startswith("youtube.com"):
            self.type = PlaylistType.YOUTUBE
            self.playlist = YoutubePlaylist(self.playlist_id)
            self.name = self.playlist.title.replace(" ", "-").replace("/", "-")
        else:
            self.type = PlaylistType.SPOTIFY
            self.playlist = spotify_client.playlist(self.playlist_id)
            self.name = self.playlist['name'].replace(" ", "-").replace("/", "-")

        self.path = os.path.join(f"/app/playlists/{self.owner}", self.name)

    def load(self):
        # if the playlist don't exist yet
        if not os.path.exists(self.path):
            os.makedirs(self.path)
            logger.debug(f"created playlist dir {self.path}")
        # get elements inside dir
        elements = os.listdir(self.path)
        if len(elements) == 0:
            logger.warning(f"0 tracks on {self.path} ")
        for el in elements:
            if os.path.isfile(os.path.join(self.path, el)):
                # append track removing .mp3/.wav ext
                self.tracks_local.append(Track(name=splitext(el)[0]))

        logger.debug(f'getting tracks for playlist: {self.name}')
        if self.type is PlaylistType.YOUTUBE:
            return self.get_youtube_tracks()
        if self.type is PlaylistType.SPOTIFY:
            return self.get_spotify_tracks()

    def clean(self):
        logger.info(f"cleaning playlist {self.name} ...")
        with os.scandir(self.path) as it:
            for entry in it:
                if entry.is_file() and not entry.name.endswith('.mp3') and not entry.name.endswith(".wav"):
                    filepath = os.path.join(self.path, entry.name)
                    os.remove(filepath)
                    logger.info(f"removed file: {entry.name}")

    def get_tracks_to_download(self):
        return list(set(self.tracks) - set(self.tracks_local))

    def get_tracks_to_remove(self):
        return list(set(self.tracks_local) - set(self.tracks))

    def get_youtube_tracks(self):
        threads = []
        for video_url in self.playlist.video_urls:
            thread = threading.Thread(target=process_yt_track, args=(video_url, self.tracks))
            threads.append(thread)
            thread.start()

        # Wait for all threads to complete
        for thread in threads:
            thread.join()

        return self.tracks

    def get_spotify_tracks(self):
        results = spotify_client.playlist_items(self.playlist_id, fields=fields)

        # (Spotify uses pagination - 100 items per page)
        while True:
            if results['next']:
                # aux
                first = results
                results = spotify_client.next(results)

                # extend
                results['items'].extend(first['items'])
            else:
                break

        total_tracks = len(results["items"])
        logger.debug(f'found {total_tracks} tracks on playlist {self.name}')

        # Get track IDs
        track_ids = [item['track']['uri'].split(':')[-1] for item in results['items']]

        # Get audio features for all tracks
        audio_features = spotify_client.audio_features(track_ids)

        threads = []

        for item, features in zip(results['items'], audio_features):
            thread = threading.Thread(target=process_sp_track, args=(item, self.tracks))
            threads.append(thread)
            thread.start()

        # Wait for all threads to complete
        for thread in threads:
            thread.join()

        return self.tracks


def process_sp_track(item, results_list):
    track_id = item['track']['uri'].split(':')[-1]
    track_name = format_track_name(item['track'])

    # create the track
    track = Track(track_id=track_id,
                  name=track_name)

    results_list.append(track)


def process_yt_track(video_url, tracks):
    # in this case, the playlist_id is the URL
    video = YouTube(video_url)

    new_track = Track(track_id=video.video_id, name=video.title)
    new_track.watch_url = video.watch_url

    tracks.append(new_track)
