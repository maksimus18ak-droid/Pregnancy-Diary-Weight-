# pregnancy_weight.py
import json
import os
import sys
import argparse
from datetime import datetime, date, timedelta

DATA_FILE = "pregnancy_weight.json"

class PregnancyWeight:
    def __init__(self):
        self.data = self.load()
        self.height = self.data.get("height")
        self.pre_weight = self.data.get("pre_weight")
        self.entries = self.data.get("entries", [])

    def load(self):
        if os.path.exists(DATA_FILE):
            with open(DATA_FILE, "r") as f:
                return json.load(f)
        return {"height": None, "pre_weight": None, "entries": []}

    def save(self):
        self.data = {
            "height": self.height,
            "pre_weight": self.pre_weight,
            "entries": self.entries
        }
        with open(DATA_FILE, "w") as f:
            json.dump(self.data, f, indent=2)

    def set_height(self, cm):
        self.height = cm
        self.save()
        print(f"✅ Height set to {cm} cm")

    def set_pre_weight(self, kg):
        self.pre_weight = kg
        self.save()
        print(f"✅ Pre‑pregnancy weight set to {kg} kg")

    def bmi(self, weight):
        if not self.height:
            return None
        h = self.height / 100.0
        return weight / (h * h)

    def bmi_category(self, bmi_val):
        if bmi_val < 18.5:
            return "Underweight"
        elif bmi_val < 25.0:
            return "Normal"
        elif bmi_val < 30.0:
            return "Overweight"
        else:
            return "Obese"

    def add(self, weight, date_str=None, note=""):
        if not self.height:
            print("❌ Please set your height first: height <cm>")
            return
        if date_str is None:
            date_str = datetime.now().strftime("%Y-%m-%d")
        entry = {
            "date": date_str,
            "weight": weight,
            "note": note,
            "bmi": self.bmi(weight)
        }
        self.entries.append(entry)
        self.entries.sort(key=lambda x: x["date"])
        self.save()
        print(f"✅ Logged {weight} kg on {date_str}")

    def list(self):
        if not self.entries:
            print("No entries.")
            return
        print("\n📋 Weight Entries:")
        for e in self.entries:
            note = f" – {e['note']}" if e.get("note") else ""
            bmi_str = f" (BMI: {e['bmi']:.1f})" if e.get("bmi") else ""
            print(f"  {e['date']}: {e['weight']:.1f} kg{bmi_str}{note}")

    def stats(self):
        if not self.entries:
            print("No entries.")
            return
        weights = [e["weight"] for e in self.entries]
        total = len(weights)
        current = weights[-1]
        avg = sum(weights) / total
        mn = min(weights)
        mx = max(weights)
        gain = current - self.pre_weight if self.pre_weight else None

        print("\n📊 Statistics:")
        print(f"  Entries: {total}")
        print(f"  Current: {current:.1f} kg")
        print(f"  Average: {avg:.1f} kg")
        print(f"  Min: {mn:.1f} kg")
        print(f"  Max: {mx:.1f} kg")
        if gain is not None:
            print(f"  Total gain: {gain:+.1f} kg")

    def show_bmi(self):
        if not self.entries:
            print("No entries.")
            return
        if not self.height:
            print("❌ Height not set.")
            return
        current = self.entries[-1]["weight"]
        bmi_val = self.bmi(current)
        cat = self.bmi_category(bmi_val)
        print(f"\n📐 Current BMI: {bmi_val:.1f} ({cat})")
        print(f"BMI Range: 18.5 - 24.9")

    def show_gain(self):
        if not self.entries:
            print("No entries.")
            return
        if self.pre_weight is None:
            print("❌ Pre‑pregnancy weight not set. Use: preweight <kg>")
            return
        current = self.entries[-1]["weight"]
        gain = current - self.pre_weight
        first = self.entries[0]["weight"]
        print(f"\n📈 Weight Gain:")
        print(f"  Pre‑pregnancy: {self.pre_weight:.1f} kg")
        print(f"  First logged: {first:.1f} kg")
        print(f"  Current: {current:.1f} kg")
        print(f"  Total gain: {gain:+.1f} kg")

    def trend(self, days=5):
        if len(self.entries) < days:
            print(f"Need at least {days} entries for trend.")
            return
        recent = self.entries[-days:]
        avg = sum(e["weight"] for e in recent) / len(recent)
        print(f"\n📈 Trend (last {days} days):")
        print(f"  Average weight: {avg:.1f} kg")
        print("  Recent entries:")
        for e in recent[-5:]:
            print(f"    {e['date']}: {e['weight']:.1f} kg")

    def export(self, filename="weights.csv"):
        import csv
        if not self.entries:
            print("No entries.")
            return
        with open(filename, 'w', newline='') as f:
            writer = csv.writer(f)
            writer.writerow(["Date", "Weight (kg)", "BMI", "Note"])
            for e in self.entries:
                writer.writerow([e["date"], e["weight"], e.get("bmi", ""), e.get("note", "")])
        print(f"✅ Exported {len(self.entries)} entries to {filename}")

def main():
    parser = argparse.ArgumentParser(description="Pregnancy Weight Tracker")
    subparsers = parser.add_subparsers(dest="cmd", required=True)

    height_parser = subparsers.add_parser("height")
    height_parser.add_argument("cm", type=float)

    pre_parser = subparsers.add_parser("preweight")
    pre_parser.add_argument("kg", type=float)

    add_parser = subparsers.add_parser("add")
    add_parser.add_argument("weight", type=float)
    add_parser.add_argument("--date", help="YYYY-MM-DD")
    add_parser.add_argument("--note", default="")

    subparsers.add_parser("list")
    subparsers.add_parser("stats")
    subparsers.add_parser("bmi")
    subparsers.add_parser("gain")
    subparsers.add_parser("trend")

    export_parser = subparsers.add_parser("export")
    export_parser.add_argument("filename", default="weights.csv", nargs="?")

    args = parser.parse_args()
    app = PregnancyWeight()

    if args.cmd == "height":
        app.set_height(args.cm)
    elif args.cmd == "preweight":
        app.set_pre_weight(args.kg)
    elif args.cmd == "add":
        app.add(args.weight, args.date, args.note)
    elif args.cmd == "list":
        app.list()
    elif args.cmd == "stats":
        app.stats()
    elif args.cmd == "bmi":
        app.show_bmi()
    elif args.cmd == "gain":
        app.show_gain()
    elif args.cmd == "trend":
        app.trend()
    elif args.cmd == "export":
        app.export(args.filename)

if __name__ == "__main__":
    main()
