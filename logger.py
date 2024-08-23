import logging
import sys
from colorama import Fore, Style, init

init(autoreset=True)

# discard pytube logger
pytube_logger = logging.getLogger('pytube')
pytube_logger.setLevel(logging.ERROR)

logger = logging.getLogger()
logger.setLevel(logging.INFO)

console_handler = logging.StreamHandler(sys.stdout)
console_handler.setLevel(logging.DEBUG)

colors = {
    'DEBUG': Fore.BLUE,
    'INFO': Fore.GREEN,
    'WARNING': Fore.YELLOW,
    'ERROR': Fore.RED,
    'CRITICAL': Fore.MAGENTA,
}


class ColoredFormatter(logging.Formatter):
    def format(self, record):
        levelname = record.levelname
        colored_levelname = f"{colors.get(levelname, Fore.RESET)}{levelname}{Style.RESET_ALL}"
        record.levelname = colored_levelname
        return super().format(record)


console_formatter = ColoredFormatter('[%(levelname)s] [%(asctime)s]: %(message)s', datefmt='%d-%m-%Y %H:%M:%S')
console_handler.setFormatter(console_formatter)

for handler in logger.handlers:
    logger.removeHandler(handler)

logger.addHandler(console_handler)
logger.setLevel('INFO')


def get_custom_logger():
    return logger
