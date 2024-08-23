import json
from updater import *
from playlist import *
from banner import show_banner

show_banner()
logger = get_custom_logger()

# start read config.json, found playlists per users and update them! (remove-download)
def start():
    logger.info("Loading all playlists...")
    playlists_by_owner = {}

    # load data
    with open('/app/config.json', 'r') as file:
        user_data = json.load(file)

    # get playlists by owner
    for owner, data in user_data['users'].items():
        playlists = [playlist['id'] for playlist in data['playlists']]
        playlists_by_owner[owner] = playlists

    # load playlists
    loaded_playlists = []
    for owner, playlists in playlists_by_owner.items():
        for playlist_id in playlists:
            logger.info(f"[{owner}] - {playlist_id} (loading)")
            p = Playlist(playlist_id, owner)
            p.load()
            loaded_playlists.append(p)

    # get the updater for playlists
    updater = Updater(audio_format="mp3")

    # update all of them
    for playlist in loaded_playlists:
        updater.update_playlist(playlist)

    for playlist in loaded_playlists:
        logger.info(f"[{playlist.name}] - local tracks: {len(playlist.tracks_local)}")

    logger.info("Finished...")


if __name__ == '__main__':
    start()
