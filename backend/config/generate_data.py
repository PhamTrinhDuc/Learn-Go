import random
import sys
import uuid
import psycopg2
from psycopg2.extras import execute_values
from datetime import date, datetime, timedelta
from faker import Faker

fake = Faker("vi_VN")
random.seed(42)

# --- DB CONFIG (MOCK) ---
DB_CONFIG = {
    "host": "localhost",
    "port": 5432,
    "database": "hair_salon",
    "user": "postgres",
    "password": "your_password"
}

def get_db_connection():
    try:
        conn = psycopg2.connect(**DB_CONFIG)
        return conn
    except Exception as e:
        print(f"Error connecting to database: {e}")
        sys.exit(1)

def uid(): return str(uuid.uuid4())

def insert_to_db(cursor, table, cols, rows):
    print(f"Inserting into {table}...")
    query = f"INSERT INTO {table} ({', '.join(cols)}) VALUES %s"
    data = [[r[col] for col in cols] for r in rows]
    execute_values(cursor, query, data)

def main():
    conn = get_db_connection()
    cur = conn.cursor()
    try:
        # ── 1. BRANCH ────────────────────────────────────────────────────────────────
        branches = [
            {"id":uid(),"name":"Chi nhánh Quận 1",     "address":"123 Nguyễn Huệ, Q.1, TP.HCM",           "phone":"02838001001","opening_hours":"08:00-20:00","is_active":True},
            {"id":uid(),"name":"Chi nhánh Quận 3",     "address":"45 Võ Văn Tần, Q.3, TP.HCM",            "phone":"02838002002","opening_hours":"08:00-20:00","is_active":True},
            {"id":uid(),"name":"Chi nhánh Bình Thạnh", "address":"78 Đinh Bộ Lĩnh, Q.Bình Thạnh, TP.HCM","phone":"02838003003","opening_hours":"08:30-19:30","is_active":True},
        ]
        insert_to_db(cur, "branch", ["id","name","address","phone","opening_hours","is_active"], branches)
        branch_ids = [b["id"] for b in branches]

        # ── 2. SERVICE ───────────────────────────────────────────────────────────────
        svc_data = [
            ("Cắt tóc nam cơ bản",   "cut",       "Cắt tạo kiểu cơ bản cho nam",       30,  30),
            ("Cắt tóc nam tạo kiểu", "cut",       "Cắt + tạo kiểu nâng cao cho nam",    45,  50),
            ("Cắt tóc nữ ngắn",      "cut",       "Cắt tóc ngắn cho nữ",                45,  50),
            ("Cắt tóc nữ dài",       "cut",       "Cắt + tỉa tóc dài cho nữ",           60,  70),
            ("Nhuộm màu cơ bản",     "color",     "Nhuộm 1 màu toàn đầu",               90, 100),
            ("Nhuộm highlight",      "color",     "Nhuộm highlight / balayage",         120, 135),
            ("Gội đầu + massage",    "treatment", "Gội đầu thư giãn + massage da đầu",   30,  30),
            ("Ủ tóc phục hồi",       "treatment", "Ủ keratin phục hồi tóc hư tổn",       60,  70),
            ("Uốn tóc",              "treatment", "Uốn xoăn / uốn lơi",                120, 140),
            ("Duỗi / Thẳng tóc",     "treatment", "Duỗi nhiệt hoặc hoá học",            120, 150),
        ]
        services = [{"id":uid(),"name":n,"category":c,"description":d,"duration_minutes":dur,"estimated_duration":est,"is_active":True}
                    for n,c,d,dur,est in svc_data]
        insert_to_db(cur, "service", ["id","name","category","description","duration_minutes","estimated_duration","is_active"], services)

        # ── 3. BRANCH_SERVICE_PRICE ──────────────────────────────────────────────────
        base_prices  = [80,120,150,180,250,350,60,180,300,280]  # nghìn đồng
        multipliers  = [1.00, 0.95, 0.90]
        bsp_rows = []
        for b, m in zip(branches, multipliers):
            for svc, base in zip(services, base_prices):
                bsp_rows.append({"id":uid(),"branch_id":b["id"],"service_id":svc["id"],
                                 "price": round(base * m) * 1000, "is_available": True})
        insert_to_db(cur, "branch_service_price", ["id","branch_id","service_id","price","is_available"], bsp_rows)

        # ── 4. STYLIST ───────────────────────────────────────────────────────────────
        sty_data = [
            ("Nguyễn Minh Tuấn","0901111001",0), ("Trần Thị Lan","0901111002",0),
            ("Lê Hoàng Nam","0901111003",0),     ("Phạm Thu Hà","0901111004",0),
            ("Võ Đình Khoa","0902222001",1),     ("Nguyễn Thị Mai","0902222002",1),
            ("Đặng Văn Hùng","0902222003",1),
            ("Bùi Thanh Tùng","0903333001",2),  ("Lý Thị Ngọc","0903333002",2),
            ("Trường Văn An","0903333003",2),
        ]
        stylists = [{"id":uid(),"branch_id":branch_ids[bi],"name":n,"phone":p,"is_active":True}
                    for n,p,bi in sty_data]
        insert_to_db(cur, "stylist", ["id","branch_id","name","phone","is_active"], stylists)

        branch_stylist_map: dict[str,list] = {b:[] for b in branch_ids}
        for s,(_,_,bi) in zip(stylists, sty_data):
            branch_stylist_map[branch_ids[bi]].append(s["id"])

        # ── 5. STYLIST_SCHEDULE ──────────────────────────────────────────────────────
        off_groups = [[1,4],[2,5],[0,3]]
        sched_rows = []
        for i, s in enumerate(stylists):
            off = off_groups[i % 3]
            for dow in range(7):
                if dow in off: continue
                weekend = dow in (0, 6)
                sched_rows.append({"id":uid(),"stylist_id":s["id"],"day_of_week":dow,
                                    "start_time":"08:00:00" if weekend else "09:00:00",
                                    "end_time":  "19:00:00" if weekend else "18:00:00",
                                    "is_active": True})
        insert_to_db(cur, "stylist_schedule", ["id","stylist_id","day_of_week","start_time","end_time","is_active"], sched_rows)

        # ── 6. USERS ─────────────────────────────────────────────────────────────────
        segs = ["new"]*30 + ["regular"]*45 + ["vip"]*10 + ["dormant"]*15
        now  = datetime(2026, 4, 26, 12, 0, 0)
        lv   = {"new":(1,30),"regular":(14,45),"vip":(7,21),"dormant":(60,180)}
        pts  = {"new":(0,50),"regular":(50,300),"vip":(300,1000),"dormant":(0,100)}

        used = {"0900000001","0900000002","0900000003","0900000004"}
        users = [
            {"id":uid(),"name":"Chủ chuỗi","phone":"0900000001","birthday":None,"address":"TP.HCM",
             "role":"owner","loyalty_points":0,"preferred_branch_id":None,"last_visit_at":None},
            {"id":uid(),"name":"Manager Q1","phone":"0900000002","birthday":None,"address":branches[0]["address"],
             "role":"manager","loyalty_points":0,"preferred_branch_id":branch_ids[0],"last_visit_at":None},
            {"id":uid(),"name":"Manager Q3","phone":"0900000003","birthday":None,"address":branches[1]["address"],
             "role":"manager","loyalty_points":0,"preferred_branch_id":branch_ids[1],"last_visit_at":None},
            {"id":uid(),"name":"Manager BT","phone":"0900000004","birthday":None,"address":branches[2]["address"],
             "role":"manager","loyalty_points":0,"preferred_branch_id":branch_ids[2],"last_visit_at":None},
        ]
        customers = []
        for _ in range(100):
            while True:
                phone = "09" + "".join(str(random.randint(0,9)) for _ in range(8))
                if phone not in used: used.add(phone); break
            seg  = random.choice(segs)
            lo,hi = lv[seg]
            lv_at = (now - timedelta(days=random.randint(lo,hi))).strftime("%Y-%m-%d %H:%M:%S+07")
            p     = random.randint(*pts[seg])
            bday  = date(random.randint(1980,2005), random.randint(1,12), random.randint(1,28))
            c = {"id":uid(),"name":fake.name(),"phone":phone,"birthday":bday.isoformat(),
                 "address":fake.address().replace("\n",", "),"role":"customer",
                 "loyalty_points":p,"preferred_branch_id":random.choice(branch_ids),"last_visit_at":lv_at}
            users.append(c); customers.append(c)

        insert_to_db(cur, "users", ["id","name","phone","birthday","address","role","loyalty_points","preferred_branch_id","last_visit_at"], users)

        # ── 7. BOOKING ───────────────────────────────────────────────────────────────
        wday_w   = [1,2,3,3,3,5,5,3,2,1,1,1]
        wend_w   = [1,2,3,4,5,5,4,3,2,1,1,1]
        n_visits = {"new":2,"regular":8,"vip":18,"dormant":1}

        bookings = []
        for c in customers:
            p = c["loyalty_points"]
            seg = "vip" if p>=300 else ("regular" if p>=50 else ("dormant" if (c["last_visit_at"] and int(c["last_visit_at"][:4])<2026) else "new"))
            n   = max(1, int(n_visits[seg] * random.uniform(0.6,1.4)))
            br  = c["preferred_branch_id"]
            pool= branch_stylist_map[br]
            for _ in range(n):
                dt  = now - timedelta(days=random.randint(1,180))
                wkd = dt.weekday() >= 5
                hr  = random.choices(range(8,20), weights=wend_w if wkd else wday_w)[0]
                mn  = random.choice([0,15,30,45])
                sch = dt.replace(hour=hr, minute=mn, second=0)
                svc = random.choice(services[4:] if seg=="vip" else (services[1:7] if seg=="regular" else services[:4]))
                status = random.choices(["completed","cancelled","no_show"], weights=[80,15,5])[0]
                bookings.append({
                    "id":uid(),"user_id":c["id"],"branch_id":br,
                    "stylist_id":random.choice(pool),"service_id":svc["id"],
                    "scheduled_at":sch.strftime("%Y-%m-%d %H:%M:%S+07"),
                    "duration_minutes":svc["duration_minutes"],
                    "estimated_duration":svc["estimated_duration"],
                    "actual_duration":None,"status":status,
                    "cancel_reason":"Khách bận đột xuất" if status=="cancelled" else None,
                    "source":random.choices(["zalo","web","agent","manual"],weights=[40,30,20,10])[0],
                    "notes":None,"check_in_at":None,"completed_at":None,
                })

        insert_to_db(cur, "booking",
            ["id","user_id","branch_id","stylist_id","service_id","scheduled_at",
             "duration_minutes","estimated_duration","actual_duration","status",
             "cancel_reason","source","notes","check_in_at","completed_at"],
            bookings)

        conn.commit()
        print("Data seeded successfully!")

    except Exception as e:
        conn.rollback()
        print(f"Error seeding data: {e}")
        import traceback
        traceback.print_exc()
    finally:
        cur.close()
        conn.close()

if __name__ == "__main__":
    main()
