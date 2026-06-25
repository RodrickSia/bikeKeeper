import urllib.request, json

BASE = "http://localhost:32065/api/v1"

# Login as student
data = json.dumps({"email": "student@university.edu.vn", "password": "demo123"}).encode()
req = urllib.request.Request(f"{BASE}/auth/login", data=data, headers={"Content-Type": "application/json"})
resp = json.load(urllib.request.urlopen(req))
token = resp["token"]
user = resp["user"]
print(f"Student: {user['email']} memberId={user.get('memberId')}")

headers = {"Authorization": f"Bearer {token}"}

# Try to deposit to a card that does NOT belong to this student
# Use a known casual card from seed - they belong to no member
deposit_data = json.dumps({"amount": 50000}).encode()
req2 = urllib.request.Request(
    f"{BASE}/cards/NFC-CSL-0001/deposit",
    data=deposit_data,
    headers={**headers, "Content-Type": "application/json"}
)
try:
    resp2 = urllib.request.urlopen(req2)
    txn = json.load(resp2)
    print(f"ERROR: Should have been blocked but got: {txn}")
except urllib.error.HTTPError as e:
    print(f"Correctly blocked: {e.code} {e.read().decode()}")

# Login as admin - should be able to deposit to any card
admin_data = json.dumps({"email": "admin@parksmart.vn", "password": "demo123"}).encode()
req3 = urllib.request.Request(f"{BASE}/auth/login", data=admin_data, headers={"Content-Type": "application/json"})
resp3 = json.load(urllib.request.urlopen(req3))
admin_token = resp3["token"]
admin_headers = {"Authorization": f"Bearer {admin_token}"}

req4 = urllib.request.Request(
    f"{BASE}/cards/NFC-CSL-0001/deposit",
    data=deposit_data,
    headers={**admin_headers, "Content-Type": "application/json"}
)
try:
    resp4 = urllib.request.urlopen(req4)
    txn = json.load(resp4)
    print(f"Admin deposit OK: {txn['id']}")
except urllib.error.HTTPError as e:
    print(f"Admin deposit blocked: {e.code} {e.read().decode()}")
