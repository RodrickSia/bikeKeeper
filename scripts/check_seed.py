import urllib.request, json

BASE = "http://localhost:32065/api/v1"

data = json.dumps({"email": "admin@parksmart.vn", "password": "demo123"}).encode()
req = urllib.request.Request(f"{BASE}/auth/login", data=data, headers={"Content-Type": "application/json"})
token = json.load(urllib.request.urlopen(req))["token"]
print(f"Token OK: {token[:20]}...")

headers = {"Authorization": f"Bearer {token}"}

for path, label in [
    ("/sessions?limit=1", "Sessions"),
    ("/members?limit=1", "Members"),
    ("/admin/users?limit=1", "Users"),
    ("/cards?limit=1", "Cards"),
    ("/shifts?limit=1", "Shifts"),
    ("/parking-lots?limit=1", "Parking Lots"),
    ("/incidents?limit=1", "Incidents"),
    ("/devices?limit=1", "Devices"),
    ("/notifications?limit=1", "Notifications"),
    ("/support-tickets?limit=1", "Support Tickets"),
    ("/card-requests?limit=1", "Card Requests"),
    ("/visitor-passes?limit=1", "Visitor Passes"),
]:
    try:
        resp = urllib.request.urlopen(urllib.request.Request(f"{BASE}{path}", headers=headers))
        data = json.load(resp)
        count = len(data) if isinstance(data, list) else 1
        print(f"  {label}: {count}")
    except Exception as e:
        status = getattr(e, "code", None) or str(e)
        print(f"  {label}: {status}")
