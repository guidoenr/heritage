import os
from colorama import init, Fore

def show_banner():
    os.system('cls' if os.name == 'nt' else 'clear') # clear
    init(autoreset=True)  # this make windows colored
    banner = f"""
    {Fore.RED}
    .__                 .__  __                         
    |  |__   ___________|__|/  |______     ____   ____  
    |  |  \_/ __ \_  __ \  \   __\__  \   / ___\_/ __ \ 
    |   Y  \  ___/|  | \/  ||  |  / __ \_/ /_/  >  ___/ 
    |___|  /\___  >__|  |__||__| (____  /\___  / \___  >
         \/     \/                    \//_____/      \/
            
    {Fore.WHITE}@ github.com/guidoenr/heritage
    (v-2024)
        """
    print(banner)

