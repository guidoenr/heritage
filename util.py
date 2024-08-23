import datetime


def get_now():
    # set the timezone to Argentina UTC-3
    argentina_tz = datetime.timezone(datetime.timedelta(hours=-3))

    # get the current date and time in Argentina UTC-3
    now = datetime.datetime.now(argentina_tz)

    # format the date as a string in the desired format
    return now.strftime("[%Y-%m-%d]")


# format_track_name formats the spotify track name
def format_track_name(track_spotify):
    name_track = track_spotify['name'].replace('"', '').replace("/", "-").strip()
    name_artists = ", ".join(artist['name'] for artist in track_spotify['artists'])

    return f"{name_track} - {name_artists}"
