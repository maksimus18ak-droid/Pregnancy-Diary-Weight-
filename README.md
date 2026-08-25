🤰 Pregnancy Diary (Weight) — Multi‑Language Weight Tracker
8 languages, one comprehensive pregnancy weight tracker – log your weight, track BMI, monitor weight gain, visualize trends, and export your data – right from your terminal.

✨ Features
⚖️ Log weight entries – with date and optional notes

📐 BMI calculation – based on height (cm) and weight (kg)

📈 Weight gain tracking – compare to pre‑pregnancy weight

📊 Statistics – current, average, min, max, total gain

📉 Trend analysis – moving average to smooth out fluctuations

📋 List entries – view all logged weights with dates

📤 Export to CSV – for further analysis in spreadsheets

💾 Persistent storage – all data saved in pregnancy_weight.json

🚀 Quick Start
All implementations follow the same CLI pattern:

bash
# Set your height (required before logging)
<command> height 165

# Set your pre‑pregnancy weight (optional, for gain tracking)
<command> preweight 60.5

# Log a weight entry
<command> add 62.3

# Log with note and date
<command> add 63.1 --note "Feeling good" --date 2026-08-25

# Show statistics
<command> stats

# Show BMI
<command> bmi

# Show weight gain
<command> gain

# Show trend (5‑day moving average)
<command> trend

# List all entries
<command> list

# Export to CSV
<command> export weights.csv
Commands:

height <cm> – set height (required)

preweight <kg> – set pre‑pregnancy weight (optional)

add <weight> [--date DATE] [--note TEXT] – log a weight

list – show all entries

stats – show statistics

bmi – show current BMI

gain – show weight gain

trend – show moving average trend

export [filename] – export to CSV

📸 Example Output
text
🤰 Pregnancy Weight Tracker
Height: 165.0 cm
Pre‑pregnancy: 60.5 kg
Current weight: 63.1 kg

📊 Statistics:
  Entries: 5
  Current: 63.1 kg
  Average: 62.4 kg
  Min: 61.2 kg
  Max: 63.1 kg
  Total gain: +2.6 kg

📐 BMI: 23.2 (Normal)
📈 Trend (5‑day avg): 62.8 kg
📁 Repository Structure
text
.
├── README.md
├── python/
│   └── pregnancy_weight.py
├── go/
│   └── pregnancy_weight.go
├── javascript/
│   └── pregnancy_weight.js
├── ruby/
│   └── pregnancy_weight.rb
├── php/
│   └── pregnancy_weight.php
├── java/
│   └── PregnancyWeight.java
├── csharp/
│   └── PregnancyWeight.cs
└── cpp/
    └── pregnancy_weight.cpp
