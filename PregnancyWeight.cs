// PregnancyWeight.cs
using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text.Json;
using System.Text.Json.Serialization;

class Entry
{
    [JsonPropertyName("date")] public string Date { get; set; }
    [JsonPropertyName("weight")] public double Weight { get; set; }
    [JsonPropertyName("note")] public string Note { get; set; }
    [JsonPropertyName("bmi")] public double BMI { get; set; }
}

class Data
{
    [JsonPropertyName("height")] public double Height { get; set; }
    [JsonPropertyName("preWeight")] public double PreWeight { get; set; }
    [JsonPropertyName("entries")] public List<Entry> Entries { get; set; } = new List<Entry>();
}

class PregnancyWeight
{
    private Data data = new Data();
    private readonly string dataFile = "pregnancy_weight.json";
    private readonly JsonSerializerOptions options = new JsonSerializerOptions { WriteIndented = true };

    public PregnancyWeight() => Load();

    private void Load()
    {
        if (!File.Exists(dataFile)) return;
        string json = File.ReadAllText(dataFile);
        data = JsonSerializer.Deserialize<Data>(json) ?? new Data();
    }

    private void Save()
    {
        string json = JsonSerializer.Serialize(data, options);
        File.WriteAllText(dataFile, json);
    }

    public void SetHeight(double cm)
    {
        data.Height = cm;
        Save();
        Console.WriteLine($"✅ Height set to {cm} cm");
    }

    public void SetPreWeight(double kg)
    {
        data.PreWeight = kg;
        Save();
        Console.WriteLine($"✅ Pre‑pregnancy weight set to {kg} kg");
    }

    private double BMI(double weight)
    {
        if (data.Height <= 0) return 0;
        double h = data.Height / 100.0;
        return weight / (h * h);
    }

    private string BMICategory(double bmiVal)
    {
        if (bmiVal < 18.5) return "Underweight";
        if (bmiVal < 25.0) return "Normal";
        if (bmiVal < 30.0) return "Overweight";
        return "Obese";
    }

    public void Add(double weight, string dateStr, string note)
    {
        if (data.Height <= 0)
        {
            Console.WriteLine("❌ Please set your height first: height <cm>");
            return;
        }
        if (string.IsNullOrEmpty(dateStr)) dateStr = DateTime.Now.ToString("yyyy-MM-dd");
        var entry = new Entry
        {
            Date = dateStr,
            Weight = weight,
            Note = note ?? "",
            BMI = BMI(weight)
        };
        data.Entries.Add(entry);
        data.Entries = data.Entries.OrderBy(e => e.Date).ToList();
        Save();
        Console.WriteLine($"✅ Logged {weight} kg on {dateStr}");
    }

    public void List()
    {
        if (!data.Entries.Any())
        {
            Console.WriteLine("No entries.");
            return;
        }
        Console.WriteLine("\n📋 Weight Entries:");
        foreach (var e in data.Entries)
        {
            string note = string.IsNullOrEmpty(e.Note) ? "" : $" – {e.Note}";
            string bmiStr = e.BMI > 0 ? $" (BMI: {e.BMI:F1})" : "";
            Console.WriteLine($"  {e.Date}: {e.Weight:F1} kg{bmiStr}{note}");
        }
    }

    public void Stats()
    {
        if (!data.Entries.Any())
        {
            Console.WriteLine("No entries.");
            return;
        }
        var weights = data.Entries.Select(e => e.Weight).ToList();
        int total = weights.Count;
        double current = weights.Last();
        double avg = weights.Average();
        double mn = weights.Min();
        double mx = weights.Max();
        double? gain = data.PreWeight > 0 ? current - data.PreWeight : (double?)null;

        Console.WriteLine("\n📊 Statistics:");
        Console.WriteLine($"  Entries: {total}");
        Console.WriteLine($"  Current: {current:F1} kg");
        Console.WriteLine($"  Average: {avg:F1} kg");
        Console.WriteLine($"  Min: {mn:F1} kg");
        Console.WriteLine($"  Max: {mx:F1} kg");
        if (gain.HasValue) Console.WriteLine($"  Total gain: {gain.Value:+F1;-F1} kg");
    }

    public void ShowBMI()
    {
        if (!data.Entries.Any())
        {
            Console.WriteLine("No entries.");
            return;
        }
        if (data.Height <= 0)
        {
            Console.WriteLine("❌ Height not set.");
            return;
        }
        var current = data.Entries.Last();
        double bmiVal = BMI(current.Weight);
        string cat = BMICategory(bmiVal);
        Console.WriteLine($"\n📐 Current BMI: {bmiVal:F1} ({cat})");
        Console.WriteLine("BMI Range: 18.5 - 24.9");
    }

    public void ShowGain()
    {
        if (!data.Entries.Any())
        {
            Console.WriteLine("No entries.");
            return;
        }
        if (data.PreWeight <= 0)
        {
            Console.WriteLine("❌ Pre‑pregnancy weight not set. Use: preweight <kg>");
            return;
        }
        var current = data.Entries.Last();
        var first = data.Entries.First();
        double gain = current.Weight - data.PreWeight;
        Console.WriteLine("\n📈 Weight Gain:");
        Console.WriteLine($"  Pre‑pregnancy: {data.PreWeight:F1} kg");
        Console.WriteLine($"  First logged: {first.Weight:F1} kg");
        Console.WriteLine($"  Current: {current.Weight:F1} kg");
        Console.WriteLine($"  Total gain: {gain:+F1;-F1} kg");
    }

    public void Trend(int days = 5)
    {
        if (data.Entries.Count < days)
        {
            Console.WriteLine($"Need at least {days} entries for trend.");
            return;
        }
        var recent = data.Entries.Skip(data.Entries.Count - days).ToList();
        double avg = recent.Average(e => e.Weight);
        Console.WriteLine($"\n📈 Trend (last {days} days):");
        Console.WriteLine($"  Average weight: {avg:F1} kg");
        Console.WriteLine("  Recent entries:");
        foreach (var e in recent.TakeLast(5))
            Console.WriteLine($"    {e.Date}: {e.Weight:F1} kg");
    }

    public void Export(string filename)
    {
        if (!data.Entries.Any())
        {
            Console.WriteLine("No entries.");
            return;
        }
        using var writer = new StreamWriter(filename);
        writer.WriteLine("Date,Weight (kg),BMI,Note");
        foreach (var e in data.Entries)
        {
            string bmiStr = e.BMI > 0 ? e.BMI.ToString("F1") : "";
            writer.WriteLine($"{e.Date},{e.Weight:F1},{bmiStr},{e.Note}");
        }
        Console.WriteLine($"✅ Exported {data.Entries.Count} entries to {filename}");
    }

    static void Main(string[] args)
    {
        if (args.Length < 1)
        {
            Console.WriteLine("Usage: PregnancyWeight <command> [options]");
            return;
        }
        var app = new PregnancyWeight();
        var parsed = ParseArgs(args);
        string cmd = args[0];
        switch (cmd)
        {
            case "height":
                if (!parsed.ContainsKey("arg")) { Console.WriteLine("height <cm>"); return; }
                app.SetHeight(double.Parse(parsed["arg"]));
                break;
            case "preweight":
                if (!parsed.ContainsKey("arg")) { Console.WriteLine("preweight <kg>"); return; }
                app.SetPreWeight(double.Parse(parsed["arg"]));
                break;
            case "add":
                if (!parsed.ContainsKey("arg")) { Console.WriteLine("add <weight> [--date DATE] [--note TEXT]"); return; }
                double w = double.Parse(parsed["arg"]);
                string date = parsed.GetValueOrDefault("date");
                string note = parsed.GetValueOrDefault("note");
                app.Add(w, date, note);
                break;
            case "list":
                app.List();
                break;
            case "stats":
                app.Stats();
                break;
            case "bmi":
                app.ShowBMI();
                break;
            case "gain":
                app.ShowGain();
                break;
            case "trend":
                int days = parsed.ContainsKey("days") ? int.Parse(parsed["days"]) : 5;
                app.Trend(days);
                break;
            case "export":
                string filename = parsed.GetValueOrDefault("arg", "weights.csv");
                app.Export(filename);
                break;
            default:
                Console.WriteLine("Unknown command.");
                break;
        }
    }

    static Dictionary<string, string> ParseArgs(string[] args)
    {
        var dict = new Dictionary<string, string>();
        for (int i = 1; i < args.Length; i++)
        {
            if (args[i].StartsWith("--") && i + 1 < args.Length)
                dict[args[i].Substring(2)] = args[++i];
            else if (!args[i].StartsWith("--"))
                dict["arg"] = args[i];
        }
        return dict;
    }
}
