// PregnancyWeight.java
import java.io.*;
import java.nio.file.*;
import java.time.*;
import java.time.format.*;
import java.util.*;
import com.google.gson.*;

class Entry {
    String date;
    double weight;
    String note;
    double bmi;
}

class Data {
    double height;
    double preWeight;
    List<Entry> entries = new ArrayList<>();
}

public class PregnancyWeight {
    private Data data = new Data();
    private final String dataFile = "pregnancy_weight.json";
    private final Gson gson = new GsonBuilder().setPrettyPrinting().create();

    public PregnancyWeight() { load(); }

    private void load() {
        try {
            Path path = Paths.get(dataFile);
            if (Files.exists(path)) {
                String json = new String(Files.readAllBytes(path));
                data = gson.fromJson(json, Data.class);
            }
        } catch (Exception e) {}
    }

    private void save() {
        try {
            Files.write(Paths.get(dataFile), gson.toJson(data).getBytes());
        } catch (Exception e) {}
    }

    public void setHeight(double cm) {
        data.height = cm;
        save();
        System.out.printf("✅ Height set to %.1f cm%n", cm);
    }

    public void setPreWeight(double kg) {
        data.preWeight = kg;
        save();
        System.out.printf("✅ Pre‑pregnancy weight set to %.1f kg%n", kg);
    }

    private double bmi(double weight) {
        if (data.height <= 0) return 0;
        double h = data.height / 100.0;
        return weight / (h * h);
    }

    private String bmiCategory(double bmiVal) {
        if (bmiVal < 18.5) return "Underweight";
        if (bmiVal < 25.0) return "Normal";
        if (bmiVal < 30.0) return "Overweight";
        return "Obese";
    }

    public void add(double weight, String dateStr, String note) {
        if (data.height <= 0) {
            System.out.println("❌ Please set your height first: height <cm>");
            return;
        }
        if (dateStr == null) dateStr = LocalDate.now().toString();
        Entry entry = new Entry();
        entry.date = dateStr;
        entry.weight = weight;
        entry.note = note != null ? note : "";
        entry.bmi = bmi(weight);
        data.entries.add(entry);
        data.entries.sort((a, b) -> a.date.compareTo(b.date));
        save();
        System.out.printf("✅ Logged %.1f kg on %s%n", weight, dateStr);
    }

    public void list() {
        if (data.entries.isEmpty()) {
            System.out.println("No entries.");
            return;
        }
        System.out.println("\n📋 Weight Entries:");
        for (Entry e : data.entries) {
            String note = e.note.isEmpty() ? "" : " – " + e.note;
            String bmiStr = e.bmi > 0 ? String.format(" (BMI: %.1f)", e.bmi) : "";
            System.out.printf("  %s: %.1f kg%s%s%n", e.date, e.weight, bmiStr, note);
        }
    }

    public void stats() {
        if (data.entries.isEmpty()) {
            System.out.println("No entries.");
            return;
        }
        List<Double> weights = new ArrayList<>();
        for (Entry e : data.entries) weights.add(e.weight);
        int total = weights.size();
        double current = weights.get(total - 1);
        double avg = weights.stream().mapToDouble(Double::doubleValue).average().orElse(0);
        double mn = weights.stream().mapToDouble(Double::doubleValue).min().orElse(0);
        double mx = weights.stream().mapToDouble(Double::doubleValue).max().orElse(0);
        Double gain = data.preWeight > 0 ? current - data.preWeight : null;

        System.out.println("\n📊 Statistics:");
        System.out.printf("  Entries: %d%n", total);
        System.out.printf("  Current: %.1f kg%n", current);
        System.out.printf("  Average: %.1f kg%n", avg);
        System.out.printf("  Min: %.1f kg%n", mn);
        System.out.printf("  Max: %.1f kg%n", mx);
        if (gain != null) System.out.printf("  Total gain: %+.1f kg%n", gain);
    }

    public void showBMI() {
        if (data.entries.isEmpty()) {
            System.out.println("No entries.");
            return;
        }
        if (data.height <= 0) {
            System.out.println("❌ Height not set.");
            return;
        }
        Entry current = data.entries.get(data.entries.size() - 1);
        double bmiVal = bmi(current.weight);
        String cat = bmiCategory(bmiVal);
        System.out.printf("\n📐 Current BMI: %.1f (%s)%n", bmiVal, cat);
        System.out.println("BMI Range: 18.5 - 24.9");
    }

    public void showGain() {
        if (data.entries.isEmpty()) {
            System.out.println("No entries.");
            return;
        }
        if (data.preWeight <= 0) {
            System.out.println("❌ Pre‑pregnancy weight not set. Use: preweight <kg>");
            return;
        }
        Entry current = data.entries.get(data.entries.size() - 1);
        Entry first = data.entries.get(0);
        double gain = current.weight - data.preWeight;
        System.out.println("\n📈 Weight Gain:");
        System.out.printf("  Pre‑pregnancy: %.1f kg%n", data.preWeight);
        System.out.printf("  First logged: %.1f kg%n", first.weight);
        System.out.printf("  Current: %.1f kg%n", current.weight);
        System.out.printf("  Total gain: %+.1f kg%n", gain);
    }

    public void trend(int days) {
        if (data.entries.size() < days) {
            System.out.printf("Need at least %d entries for trend.%n", days);
            return;
        }
        List<Entry> recent = data.entries.subList(data.entries.size() - days, data.entries.size());
        double avg = recent.stream().mapToDouble(e -> e.weight).average().orElse(0);
        System.out.printf("\n📈 Trend (last %d days):%n", days);
        System.out.printf("  Average weight: %.1f kg%n", avg);
        System.out.println("  Recent entries:");
        for (int i = Math.max(0, recent.size() - 5); i < recent.size(); i++) {
            Entry e = recent.get(i);
            System.out.printf("    %s: %.1f kg%n", e.date, e.weight);
        }
    }

    public void export(String filename) throws Exception {
        if (data.entries.isEmpty()) {
            System.out.println("No entries.");
            return;
        }
        try (BufferedWriter writer = Files.newBufferedWriter(Paths.get(filename))) {
            writer.write("Date,Weight (kg),BMI,Note\n");
            for (Entry e : data.entries) {
                String bmiStr = e.bmi > 0 ? String.format("%.1f", e.bmi) : "";
                writer.write(String.format("%s,%.1f,%s,%s%n", e.date, e.weight, bmiStr, e.note));
            }
        }
        System.out.printf("✅ Exported %d entries to %s%n", data.entries.size(), filename);
    }

    public static void main(String[] args) throws Exception {
        if (args.length < 1) {
            System.out.println("Usage: PregnancyWeight <command> [options]");
            return;
        }
        PregnancyWeight app = new PregnancyWeight();
        String cmd = args[0];
        Map<String, String> params = new HashMap<>();
        for (int i=1; i<args.length; i++) {
            if (args[i].startsWith("--") && i+1 < args.length) {
                params.put(args[i].substring(2), args[++i]);
            } else if (!args[i].startsWith("--")) {
                params.put("arg", args[i]);
            }
        }
        switch (cmd) {
            case "height":
                if (!params.containsKey("arg")) { System.out.println("height <cm>"); return; }
                app.setHeight(Double.parseDouble(params.get("arg")));
                break;
            case "preweight":
                if (!params.containsKey("arg")) { System.out.println("preweight <kg>"); return; }
                app.setPreWeight(Double.parseDouble(params.get("arg")));
                break;
            case "add":
                if (!params.containsKey("arg")) { System.out.println("add <weight> [--date DATE] [--note TEXT]"); return; }
                double weight = Double.parseDouble(params.get("arg"));
                String date = params.get("date");
                String note = params.get("note");
                app.add(weight, date, note);
                break;
            case "list":
                app.list();
                break;
            case "stats":
                app.stats();
                break;
            case "bmi":
                app.showBMI();
                break;
            case "gain":
                app.showGain();
                break;
            case "trend":
                int days = params.containsKey("days") ? Integer.parseInt(params.get("days")) : 5;
                app.trend(days);
                break;
            case "export":
                String filename = params.getOrDefault("arg", "weights.csv");
                app.export(filename);
                break;
            default:
                System.out.println("Unknown command.");
        }
    }
}
