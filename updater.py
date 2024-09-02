import os
import threading
from logger import get_custom_logger

logger = get_custom_logger()


class Updater:
    def __init__(self, audio_format):
        self.audio_format = audio_format
        self.base_path = "playlists"
        logger.info(f'Updater: [{audio_format}] format')

    def update_playlist(self, playlist):
        tracks_to_download = playlist.get_tracks_to_download()
        tracks_to_remove = playlist.get_tracks_to_remove()

        # fuck others quality
        if playlist.owner == "guido":
            self.audio_format = "wav"

        if tracks_to_download or tracks_to_remove:
            logger.info("tracks to download: %s", '\n'.join(t.name for t in tracks_to_download))
            logger.warning("tracks to remove: %s", '\n'.join(t.name for t in tracks_to_remove))

        for tr in tracks_to_remove:
            logger.debug(f"removing track: {tr.name} ...")
            track_path = os.path.join(playlist.path, tr.name + f".{self.audio_format}")
            try:
                os.remove(track_path)
            except FileNotFoundError:
                logger.warning(f"track path {track_path} not found...")

        threads = []
        for td in tracks_to_download:
            thread = threading.Thread(target=self.download_track, args=(td, playlist.path))
            threads.append(thread)
            thread.start()

        
        # wait threads to finish
        for thread in threads:
            thread.join()

        playlist.clean()

    def download_track(self, track, dir_name):
        # find the watch url for that track
        track.get_watch_url()

        # download using binary
        bash_command = f"./yt-dlp -f bestaudio --extract-audio -k --audio-format {self.audio_format} -o '{dir_name}/{track.name}.%(ext)s' '{track.watch_url}'"

        try:
            # run sh
            return_code = os.system(bash_command)
            if return_code == 0:
                logger.debug(f'track {track.name} downloaded!')
            else:
                logger.critical(f'error while downloading track {track.name}, return code={return_code}')
        except Exception as e:
            logger.critical(f'cannot download track: {e}')
