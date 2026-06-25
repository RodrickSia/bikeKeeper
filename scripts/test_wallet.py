import urllib.request, json

BASE = "http://localhost:32065/api/v1"

# 1. Login as student
data = json.dumps({"email": "student@university.edu.vn", "password": "demo123"}).encode()
req = urllib.request.Request(f"{BASE}/auth/login", data=data, headers={"Content-Type": "application/json"})
resp = json.load(urllib.request.urlopen(req))
token = resp["token"]
user = resp["user"]
print(f"1. Logged in: {user['email']} role={user['role']} memberId={user.get('memberId')}")

headers = {"Authorization": f"Bearer {token}"}

# 2. Try to get cards by member
mid = user.get("memberId")
if mid:
    req2 = urllib.request.Request(f"{BASE}/members-cards/{mid}", headers=headers)
    try:
        resp2 = urllib.request.urlopen(req2)
        cards = json.load(resp2)
        print(f"2. Cards by member ({mid}): {len(cards)} found")
        for c in cards:
            print(f"   {c['cardUid']} balance={c['balance']} status={c['status']}")
        
        # 3. Try deposit
        if cards:
            card = cards[0]
            deposit_data = json.dumps({"amount": 50000}).encode()
            req3 = urllib.request.Request(
                f"{BASE}/cards/{card['cardUid']}/deposit",
                data=deposit_data,
                headers={**headers, "Content-Type": "application/json"}
            )
            try:
                resp3 = urllib.request.urlopen(req3)
                txn = json.load(resp3)
                print(f"3. Deposit OK: {txn['amount']} -> {txn['id']}")
            except urllib.error.HTTPError as e:
                print(f"3. Deposit error: {e.code} {e.read().decode()}")
        
        # 4. Get transactions
        req4 = urllib.request.Request(f"{BASE}/cards/{card['cardUid']}/transactions", headers=headers)
        try:
            resp4 = urllib.request.urlopen(req4)
            txns = json.load(resp4)
            print(f"4. Transactions: {len(txns)} found")
        except urllib.error.HTTPError as e:
            print(f"4. Transactions error: {e.code} {e.read().decode()}")
        
        # 5. Get balance
        req5 = urllib.request.Request(f"{BASE}/cards/{card['cardUid']}/balance", headers=headers)
        try:
            resp5 = urllib.request.urlopen(req5)
            bal = json.load(resp5)
            print(f"5. Balance: {bal['balance']}")
        except urllib.error.HTTPError as e:
            print(f"5. Balance error: {e.code} {e.read().decode()}")
            
    except urllib.error.HTTPError as e:
        print(f"2. ERROR: {e.code} {e.reason}")
        print(f"   Body: {e.read().decode()}")
else:
    print("2. No memberId - cannot fetch cards")
