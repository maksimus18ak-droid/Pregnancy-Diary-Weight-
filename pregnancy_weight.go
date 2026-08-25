// pregnancy_weight.go
package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"
)

type Entry struct {
	Date   string  `json:"date"`
	Weight float64 `json:"weight"`
	Note   string  `json:"note"`
	BMI    float64 `json:"bmi"`
}

type Data struct {
	Height    float64 `json:"height"`
	PreWeight float64 `json:"pre_weight"`
	Entries   []Entry `json:"entries"`
}

type App struct {
	data Data
	file string
}

func NewApp(file string) *App {
	app := &App{file: file}
	app.load()
	return app
}

func (a *App) load() {
	f, err := os.ReadFile(a.file)
	if err != nil {
		return
	}
	json.Unmarshal(f, &a.data)
}

func (a *App) save() {
	data, _ := json.MarshalIndent(a.data, "", "  ")
	os.WriteFile(a.file, data, 0644)
}

func (a *App) setHeight(cm float64) {
	a.data.Height = cm
	a.save()
	fmt.Printf("✅ Height set to %.1f cm\n", cm)
}

func (a *App) setPreWeight(kg float64) {
	a.data.PreWeight = kg
	a.save()
	fmt.Printf("✅ Pre‑pregnancy weight set to %.1f kg\n", kg)
}

func (a *App) bmi(weight float64) float64 {
	if a.data.Height <= 0 {
		return 0
	}
	h := a.data.Height / 100.0
	return weight / (h * h)
}

func (a *App) bmiCategory(bmiVal float64) string {
	if bmiVal < 18.5 {
		return "Underweight"
	} else if bmiVal < 25.0 {
		return "Normal"
	} else if bmiVal < 30.0 {
		return "Overweight"
	} else {
		return "Obese"
	}
}

func (a *App) add(weight float64, dateStr string, note string) {
	if a.data.Height <= 0 {
		fmt.Println("❌ Please set your height first: height <cm>")
		return
	}
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}
	bmiVal := a.bmi(weight)
	entry := Entry{
		Date:   dateStr,
		Weight: weight,
		Note:   note,
		BMI:    bmiVal,
	}
	a.data.Entries = append(a.data.Entries, entry)
	sort.Slice(a.data.Entries, func(i, j int) bool {
		return a.data.Entries[i].Date < a.data.Entries[j].Date
	})
	a.save()
	fmt.Printf("✅ Logged %.1f kg on %s\n", weight, dateStr)
}

func (a *App) list() {
	if len(a.data.Entries) == 0 {
		fmt.Println("No entries.")
		return
	}
	fmt.Println("\n📋 Weight Entries:")
	for _, e := range a.data.Entries {
		note := ""
		if e.Note != "" {
			note = " – " + e.Note
		}
		bmiStr := ""
		if e.BMI > 0 {
			bmiStr = fmt.Sprintf(" (BMI: %.1f)", e.BMI)
		}
		fmt.Printf("  %s: %.1f kg%s%s\n", e.Date, e.Weight, bmiStr, note)
	}
}

func (a *App) stats() {
	if len(a.data.Entries) == 0 {
		fmt.Println("No entries.")
		return
	}
	weights := make([]float64, len(a.data.Entries))
	for i, e := range a.data.Entries {
		weights[i] = e.Weight
	}
	current := weights[len(weights)-1]
	total := len(weights)
	sum := 0.0
	min := weights[0]
	max := weights[0]
	for _, w := range weights {
		sum += w
		if w < min {
			min = w
		}
		if w > max {
			max = w
		}
	}
	avg := sum / float64(total)
	gain := 0.0
	if a.data.PreWeight > 0 {
		gain = current - a.data.PreWeight
	}
	fmt.Println("\n📊 Statistics:")
	fmt.Printf("  Entries: %d\n", total)
	fmt.Printf("  Current: %.1f kg\n", current)
	fmt.Printf("  Average: %.1f kg\n", avg)
	fmt.Printf("  Min: %.1f kg\n", min)
	fmt.Printf("  Max: %.1f kg\n", max)
	if a.data.PreWeight > 0 {
		fmt.Printf("  Total gain: %+.1f kg\n", gain)
	}
}

func (a *App) showBMI() {
	if len(a.data.Entries) == 0 {
		fmt.Println("No entries.")
		return
	}
	if a.data.Height <= 0 {
		fmt.Println("❌ Height not set.")
		return
	}
	current := a.data.Entries[len(a.data.Entries)-1].Weight
	bmiVal := a.bmi(current)
	cat := a.bmiCategory(bmiVal)
	fmt.Printf("\n📐 Current BMI: %.1f (%s)\n", bmiVal, cat)
	fmt.Println("BMI Range: 18.5 - 24.9")
}

func (a *App) showGain() {
	if len(a.data.Entries) == 0 {
		fmt.Println("No entries.")
		return
	}
	if a.data.PreWeight <= 0 {
		fmt.Println("❌ Pre‑pregnancy weight not set. Use: preweight <kg>")
		return
	}
	current := a.data.Entries[len(a.data.Entries)-1].Weight
	first := a.data.Entries[0].Weight
	gain := current - a.data.PreWeight
	fmt.Println("\n📈 Weight Gain:")
	fmt.Printf("  Pre‑pregnancy: %.1f kg\n", a.data.PreWeight)
	fmt.Printf("  First logged: %.1f kg\n", first)
	fmt.Printf("  Current: %.1f kg\n", current)
	fmt.Printf("  Total gain: %+.1f kg\n", gain)
}

func (a *App) trend(days int) {
	if len(a.data.Entries) < days {
		fmt.Printf("Need at least %d entries for trend.\n", days)
		return
	}
	recent := a.data.Entries[len(a.data.Entries)-days:]
	sum := 0.0
	for _, e := range recent {
		sum += e.Weight
	}
	avg := sum / float64(len(recent))
	fmt.Printf("\n📈 Trend (last %d days):\n", days)
	fmt.Printf("  Average weight: %.1f kg\n", avg)
	fmt.Println("  Recent entries:")
	start := len(recent) - 5
	if start < 0 {
		start = 0
	}
	for _, e := range recent[start:] {
		fmt.Printf("    %s: %.1f kg\n", e.Date, e.Weight)
	}
}

func (a *App) export(filename string) {
	if len(a.data.Entries) == 0 {
		fmt.Println("No entries.")
		return
	}
	f, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()
	w.Write([]string{"Date", "Weight (kg)", "BMI", "Note"})
	for _, e := range a.data.Entries {
		bmiStr := ""
		if e.BMI > 0 {
			bmiStr = fmt.Sprintf("%.1f", e.BMI)
		}
		w.Write([]string{e.Date, fmt.Sprintf("%.1f", e.Weight), bmiStr, e.Note})
	}
	fmt.Printf("✅ Exported %d entries to %s\n", len(a.data.Entries), filename)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: pregnancy_weight <command> [options]")
		return
	}
	app := NewApp("pregnancy_weight.json")
	cmd := os.Args[1]

	switch cmd {
	case "height":
		if len(os.Args) < 3 {
			fmt.Println("height <cm>")
			return
		}
		cm, _ := strconv.ParseFloat(os.Args[2], 64)
		app.setHeight(cm)
	case "preweight":
		if len(os.Args) < 3 {
			fmt.Println("preweight <kg>")
			return
		}
		kg, _ := strconv.ParseFloat(os.Args[2], 64)
		app.setPreWeight(kg)
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("add <weight> [--date DATE] [--note TEXT]")
			return
		}
		weight, _ := strconv.ParseFloat(os.Args[2], 64)
		dateStr := ""
		note := ""
		for i := 3; i < len(os.Args); i++ {
			if os.Args[i] == "--date" && i+1 < len(os.Args) {
				dateStr = os.Args[i+1]
				i++
			}
			if os.Args[i] == "--note" && i+1 < len(os.Args) {
				note = os.Args[i+1]
				i++
			}
		}
		app.add(weight, dateStr, note)
	case "list":
		app.list()
	case "stats":
		app.stats()
	case "bmi":
		app.showBMI()
	case "gain":
		app.showGain()
	case "trend":
		days := 5
		for i := 2; i < len(os.Args); i++ {
			if os.Args[i] == "--days" && i+1 < len(os.Args) {
				days, _ = strconv.Atoi(os.Args[i+1])
				i++
			}
		}
		app.trend(days)
	case "export":
		filename := "weights.csv"
		if len(os.Args) > 2 {
			filename = os.Args[2]
		}
		app.export(filename)
	default:
		fmt.Println("Unknown command. Use height, preweight, add, list, stats, bmi, gain, trend, export.")
	}
}
