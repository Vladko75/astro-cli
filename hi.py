from datetime import datetime
from zoneinfo import ZoneInfo

def show_current_time_marmaris():
    tz = ZoneInfo('Europe/Istanbul')
    now = datetime.now(tz)
    print(f"Current time in Marmaris, Turkey: {now.strftime('%Y-%m-%d %H:%M:%S %Z')}")

if __name__ == "__main__":
    show_current_time_marmaris()
