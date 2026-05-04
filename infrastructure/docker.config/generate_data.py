import random
import sys
import uuid
import psycopg2
from psycopg2.extras import execute_values
from datetime import date, datetime, timedelta
from faker import Faker
from werkzeug.security import generate_password_hash

fake = Faker("vi_VN")
random.seed(42)

# --- DB CONFIG (MOCK) ---
DB_CONFIG = {
    "host": "localhost",
    "port": 5433,
    "database": "salon_chain",
    "user": "mcp_user",
    "password": "mcp_password"
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
            {"id":uid(),"name":"Chủ chuỗi","phone":"0900000001","email":"owner@salon.com","username":"owner","password_plain":"password123","birthday":None,"address":"TP.HCM",
             "role":"owner","loyalty_points":0,"preferred_branch_id":None,"last_visit_at":None},
            {"id":uid(),"name":"Manager Q1","phone":"0900000002","email":"managerq1@salon.com","username":"managerq1","password_plain":"password123","birthday":None,"address":branches[0]["address"],
             "role":"manager","loyalty_points":0,"preferred_branch_id":branch_ids[0],"last_visit_at":None},
            {"id":uid(),"name":"Manager Q3","phone":"0900000003","email":"managerq3@salon.com","username":"managerq3","password_plain":"password123","birthday":None,"address":branches[1]["address"],
             "role":"manager","loyalty_points":0,"preferred_branch_id":branch_ids[1],"last_visit_at":None},
            {"id":uid(),"name":"Manager BT","phone":"0900000004","email":"managerbt@salon.com","username":"managerbt","password_plain":"password123","birthday":None,"address":branches[2]["address"],
             "role":"manager","loyalty_points":0,"preferred_branch_id":branch_ids[2],"last_visit_at":None},
        ]
        
        # Tạo thêm 10 khách hàng với thông tin thực tế hơn
        customers = []
        for i in range(10):
            while True:
                phone = "09" + "".join(str(random.randint(0,9)) for _ in range(8))
                if phone not in used: 
                    used.add(phone)
                    break
            
            full_name = fake.name()
            # Tạo username từ tên (không dấu)
            username_base = "".join(filter(str.isalnum, full_name.lower())).replace("đ", "d")
            username = f"{username_base[:10]}{random.randint(100, 999)}"
            
            bday = date(random.randint(1985, 2005), random.randint(1, 12), random.randint(1, 28))
            
            c = {
                "id": uid(),
                "name": full_name,
                "phone": phone,
                "email": fake.email(),
                "username": username,
                "password_plain": "customer123",
                "birthday": bday.isoformat(),
                "address": fake.address().replace("\n", ", "),
                "role": "customer",
                "loyalty_points": random.randint(0, 500),
                "preferred_branch_id": random.choice(branch_ids),
                "last_visit_at": (now - timedelta(days=random.randint(1, 60))).strftime("%Y-%m-%d %H:%M:%S+07")
            }
            users.append(c)
            customers.append(c)

        # Hash passwords before insertion and print for visibility
        print("\n--- GENERATED USERS CREDENTIALS ---")
        for u in users:
            u["password_hash"] = generate_password_hash(u["password_plain"])
            print(f"Role: {u['role']:<10} | Username: {u['username']:<15} | Password: {u['password_plain']}")
        print("-----------------------------------\n")

        insert_to_db(cur, "users", ["id","name","phone","email","username","password_hash","birthday","address","role","loyalty_points","preferred_branch_id","last_visit_at"], users)

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
                    "status":status,
                    "cancel_reason":"Khách bận đột xuất" if status=="cancelled" else None,
                    "source":random.choices(["zalo","web","agent","manual"],weights=[40,30,20,10])[0]
                })

        insert_to_db(cur, "booking",
            ["id","user_id","branch_id","stylist_id","service_id","scheduled_at",
             "duration_minutes","status",
             "cancel_reason","source"],
            bookings)

        # ── 8. PRODUCT ───────────────────────────────────────────────────────────────
        prod_data = [
            # --- WAX & POMADE (Retail/Both) ---
            ("Sáp Kevin Murphy Rough Rider", "wax", 550000, 680000, "both"),
            ("Sáp Hanz de Fuko Quicksand",   "wax", 500000, 620000, "both"),
            ("Sáp Hanz de Fuko Claymation",  "wax", 500000, 620000, "both"),
            ("Sáp Blumaan Cavalier Clay",    "wax", 480000, 580000, "retail"),
            ("Sáp Blumaan Monarch Matte Paste", "wax", 480000, 580000, "retail"),
            ("Sáp Shear Revival Northern Lights", "wax", 450000, 550000, "retail"),
            ("Pomade Reuzel Blue (Strong Hold)", "pomade", 350000, 450000, "both"),
            ("Pomade Reuzel Pink (Heavy Grease)", "pomade", 350000, 450000, "both"),
            ("Pomade Reuzel Fiber",          "pomade", 350000, 450000, "both"),
            ("Pomade Suavecito Original Hold", "pomade", 300000, 380000, "retail"),
            ("Pomade Layrite Superhold",     "pomade", 380000, 480000, "retail"),

            # --- SHAMPOO & CONDITIONER (Both/Internal) ---
            ("Dầu gội Morris Motley Treatment", "shampoo", 850000, 1100000, "both"),
            ("Dầu xả Morris Motley",         "conditioner", 850000, 1100000, "both"),
            ("Dầu gội Kevin Murphy Stimulate-Me", "shampoo", 550000, 720000, "both"),
            ("Dầu gội Tigi Bed Head (Đỏ)",   "shampoo", 320000, 450000, "retail"),
            ("Dầu gội bưởi Thorakao 500ml",  "shampoo", 60000, 85000, "retail"),
            ("Dầu gội công nghiệp (Can 5L)", "shampoo", 250000, 0, "internal"),

            # --- PRE-STYLING & SPRAY (Retail/Both) ---
            ("Bona Fide Texture Spray",      "pre-styling", 380000, 480000, "both"),
            ("Sidekick By Vilain",           "pre-styling", 420000, 520000, "retail"),
            ("Gôm xịt tóc Osis+ 3 Session",  "spray", 220000, 300000, "both"),
            ("Gôm xịt tóc 2Vee",             "spray", 180000, 250000, "retail"),

            # --- ACCESSORIES & TOOLS (Retail) ---
            ("Lược Kent Comb A81T",          "comb", 150000, 220000, "retail"),
            ("Lược tạo phồng Chaoba",        "comb", 40000, 70000, "retail"),
            ("Máy sấy tóc Chaoba 2800W",     "tool", 250000, 350000, "retail"),
            ("Dầu dưỡng râu Prospectors",    "beard-oil", 350000, 450000, "retail"),

            # --- SUPPLIES (Internal Only) ---
            ("Lưỡi dao lam Dorco (Hộp 100)", "supply", 120000, 0, "internal"),
            ("Bột talc làm sạch cổ",         "supply", 45000, 0, "internal"),
            ("Dung dịch sát khuẩn Barbicide", "supply", 550000, 0, "internal"),
            ("Khăn giấy cổ (Cuộn)",          "supply", 25000, 0, "internal"),
        ]
        
        products = []
        for n, c, pi, po, ut in prod_data:
            products.append({
                "id": uid(), "name": n, "category": c, 
                "price_in": pi, "price_out": po, "usage_type": ut,
                "low_stock_threshold_retail": 5, "low_stock_threshold_internal": 3,
                "is_active": True
            })
        insert_to_db(cur, "product", ["id","name","category","price_in","price_out","usage_type","low_stock_threshold_retail","low_stock_threshold_internal","is_active"], products)

        # ── 9. INVENTORY ─────────────────────────────────────────────────────────────
        inventory_rows = []
        inventory_logs = []
        manager_id = next(u["id"] for u in users if u["role"] == "manager")

        for b_id in branch_ids:
            for p in products:
                # Mỗi chi nhánh nhập ngẫu nhiên số lượng
                qty_retail = random.randint(10, 30) if p["usage_type"] in ["retail", "both"] else 0
                qty_internal = random.randint(5, 15) if p["usage_type"] in ["internal", "both"] else 0
                inv_id = uid()
                
                inventory_rows.append({
                    "id": inv_id, "product_id": p["id"], "branch_id": b_id,
                    "quantity_total": qty_retail + qty_internal,
                    "quantity_retail": qty_retail,
                    "quantity_internal": qty_internal
                })
                
                # Log nhập kho ban đầu
                inventory_logs.append({
                    "id": uid(), "inventory_id": inv_id, "action_type": "import",
                    "qty_change": qty_retail + qty_internal, "note": "Nhập hàng đầu kỳ",
                    "performed_by": manager_id, "performer_role": "manager",
                    "created_at": (now - timedelta(days=200)).strftime("%Y-%m-%d %H:%M:%S+07")
                })

        insert_to_db(cur, "inventory", ["id","product_id","branch_id","quantity_total","quantity_retail","quantity_internal"], inventory_rows)
        
        # ── 10. ORDERS & ORDER ITEMS ──────────────────────────────────────────────────
        order_rows = []
        order_item_rows = []
        
        # Map để trừ tồn kho thực tế khi tạo order
        inv_map = {(r["branch_id"], r["product_id"]): r for r in inventory_rows}

        for c in customers:
            # Mỗi khách mua 0-3 đơn hàng
            for _ in range(random.randint(0, 3)):
                br_id = c["preferred_branch_id"]
                order_id = uid()
                order_date = now - timedelta(days=random.randint(1, 150))
                
                # Chọn ngẫu nhiên 1-3 sản phẩm (chỉ lấy loại bán lẻ)
                retail_prods = [p for p in products if p["usage_type"] in ["retail", "both"]]
                bought_prods = random.sample(retail_prods, k=random.randint(1, min(3, len(retail_prods))))
                
                total_amount = 0
                for p in bought_prods:
                    qty = random.randint(1, 2)
                    price = float(p["price_out"])
                    total_amount += price * qty
                    
                    order_item_rows.append({
                        "id": uid(), "order_id": order_id, "product_id": p["id"],
                        "quantity": qty, "unit_price": price
                    })
                    
                    # Log xuất kho bán hàng
                    inv_rec = inv_map.get((br_id, p["id"]))
                    if inv_rec:
                        inv_rec["quantity_retail"] -= qty
                        inv_rec["quantity_total"] -= qty
                        inventory_logs.append({
                            "id": uid(), "inventory_id": inv_rec["id"], "action_type": "sale",
                            "qty_change": -qty, "note": f"Bán hàng đơn {order_id[:8]}",
                            "performed_by": None, "performer_role": "agent",
                            "created_at": order_date.strftime("%Y-%m-%d %H:%M:%S+07")
                        })

                order_rows.append({
                    "id": order_id, "user_id": c["id"], "branch_id": br_id,
                    "total_amount": total_amount, "points_earned": int(total_amount // 10000),
                    "payment_status": True, "payment_method": random.choice(["cash", "transfer", "momo"]),
                    "created_at": order_date.strftime("%Y-%m-%d %H:%M:%S+07")
                })

        insert_to_db(cur, "orders", ["id","user_id","branch_id","total_amount","points_earned","payment_status","payment_method","created_at"], order_rows)
        insert_to_db(cur, "order_items", ["id","order_id","product_id","quantity","unit_price"], order_item_rows)
        
        # Cập nhật lại số lượng tồn kho sau khi trừ bán hàng
        print("Updating inventory quantities after sales...")
        for inv in inventory_rows:
            cur.execute("UPDATE inventory SET quantity_total = %s, quantity_retail = %s WHERE id = %s",
                        (inv["quantity_total"], inv["quantity_retail"], inv["id"]))
        
        # Insert logs cuối cùng
        insert_to_db(cur, "inventory_log", ["id","inventory_id","action_type","qty_change","note","performed_by","performer_role","created_at"], inventory_logs)

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
