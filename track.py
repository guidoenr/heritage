from pytube import Search
from logger import get_custom_logger

logger = get_custom_logger()


class Track:
    def __init__(self, track_id=None, name=None):
        self.id = track_id
        self.name = name.replace("/", "-").replace("'", "").replace("(", "").replace("(", "")
        self.watch_url = "empty"

    def get_watch_url(self):
        first = None
        try:
            s = Search(self.name)
            first = s.results[0]
            s.get_next_results()

        except Exception as e:
            logger.error(str(e))

        # set the first result
        self.watch_url = first.watch_url

    # implement the __eq__ method for equality comparison (happens when tracks have the same name)
    def __eq__(self, other):
        if isinstance(other, Track):
            return self.name == other.name
        return False

    # implement the __hash__ method for hashing
    def __hash__(self):
        return hash(self.name)
