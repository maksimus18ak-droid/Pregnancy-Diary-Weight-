// pregnancy_weight.cpp
#include <iostream>
#include <fstream>
#include <string>
#include <vector>
#include <map>
#include <ctime>
#include <iomanip>
#include <sstream>
#include <algorithm>
#include <nlohmann/json.hpp>
#include <getopt.h>

using namespace std;
using json = nlohmann::json;

struct Entry {
    string date;
    double weight;
    string note;
    double bmi;
};

class PregnancyWeight {
private:
    double height = 0;
    double preWeight = 0;
    vector<Entry> entries;
    string dataFile = "pregnancy_weight.json";

    void load() {
        ifstream f(dataFile);
        if (!f.is_open()) return;
        json j;
        f >> j;
        if (j.contains("height")) height = j["height"];
        if (j.contains("preWeight")) preWeight = j["preWeight"];
        if (j.contains("entries")) {
            for (auto& item : j["entries"]) {
                Entry e;
                e.date = item["date"];
                e.weight = item["weight"];
                e.note = item["note"];
                e.bmi = item["bmi"];
                entries.push_back(e);
            }
        }
    }

    void save() {
        json j;
        j["height"] = height;
        j["preWeight"] = preWeight;
        j["entries"] = json::array();
        for (auto& e : entries) {
            j["entries"].push_back({
                {"date", e.date},
                {"weight", e.weight},
                {"note", e.note},
                {"bmi", e.bmi}
            });
        }
        ofstream f(dataFile);
        f << setw(2) << j << endl;
    }

    double bmi(double weight) {
        if (height <= 0) return 0;
        double h = height / 100.0;
        return weight / (h * h);
    }

    string bmiCategory(double bmiVal) {
        if (bmiVal < 18.5) return "Underweight";
        if (bmiVal < 25.0) return "Normal";
        if (bmiVal < 30.0) return "Overweight";
        return "Obese";
    }

    string currentDate() {
        time_t t = time(nullptr);
        char buf[11];
        strftime(buf, sizeof(buf), "%Y-%m-%d", localtime(&t));
        return string(buf);
    }

public:
    PregnancyWeight() { load(); }

    void setHeight(double cm) {
        height = cm;
        save();
        cout << "✅ Height set to " << cm << " cm\n";
    }

    void setPreWeight(double kg) {
        preWeight = kg;
        save();
        cout << "✅ Pre‑pregnancy weight set to " << kg << " kg\n";
    }

    void add(double weight, const string& dateStr, const string& note) {
        if (height <= 0) {
            cout << "❌ Please set your height first: height <cm>\n";
            return;
        }
        string d = dateStr.empty() ? currentDate() : dateStr;
        Entry e;
        e.date = d;
        e.weight = weight;
        e.note = note;
        e.bmi = bmi(weight);
        entries.push_back(e);
        sort(entries.begin(), entries.end(), [](const Entry& a, const Entry& b) {
            return a.date < b.date;
        });
        save();
        cout << "✅ Logged " << weight << " kg on " << d << "\n";
    }

    void list() {
        if (entries.empty()) {
            cout << "No entries.\n";
            return;
        }
        cout << "\n📋 Weight Entries:\n";
        for (auto& e : entries) {
            string note = e.note.empty() ? "" : " – " + e.note;
            string bmiStr = e.bmi > 0 ? " (BMI: " + to_string(e.bmi).substr(0,4) + ")" : "";
            cout << "  " << e.date << ": " << e.weight << " kg" << bmiStr << note << "\n";
        }
    }

    void stats() {
        if (entries.empty()) {
            cout << "No entries.\n";
            return;
        }
        vector<double> weights;
        for (auto& e : entries) weights.push_back(e.weight);
        int total = weights.size();
        double current = weights.back();
        double avg = 0;
        double mn = weights[0], mx = weights[0];
        for (double w : weights) {
            avg += w;
            if (w < mn) mn = w;
            if (w > mx) mx = w;
        }
        avg /= total;
        double gain = preWeight > 0 ? current - preWeight : 0;

        cout << "\n📊 Statistics:\n";
        cout << "  Entries: " << total << "\n";
        cout << "  Current: " << current << " kg\n";
        cout << "  Average: " << avg << " kg\n";
        cout << "  Min: " << mn << " kg\n";
        cout << "  Max: " << mx << " kg\n";
        if (preWeight > 0) cout << "  Total gain: " << gain << " kg\n";
    }

    void showBMI() {
        if (entries.empty()) {
            cout << "No entries.\n";
            return;
        }
        if (height <= 0) {
            cout << "❌ Height not set.\n";
            return;
        }
        double current = entries.back().weight;
        double bmiVal = bmi(current);
        string cat = bmiCategory(bmiVal);
        cout << "\n📐 Current BMI: " << bmiVal << " (" << cat << ")\n";
        cout << "BMI Range: 18.5 - 24.9\n";
    }

    void showGain() {
        if (entries.empty()) {
            cout << "No entries.\n";
            return;
        }
        if (preWeight <= 0) {
            cout << "❌ Pre‑pregnancy weight not set. Use: preweight <kg>\n";
            return;
        }
        double current = entries.back().weight;
        double first = entries.front().weight;
        double gain = current - preWeight;
        cout << "\n📈 Weight Gain:\n";
        cout << "  Pre‑pregnancy: " << preWeight << " kg\n";
        cout << "  First logged: " << first << " kg\n";
        cout << "  Current: " << current << " kg\n";
        cout << "  Total gain: " << gain << " kg\n";
    }

    void trend(int days = 5) {
        if ((int)entries.size() < days) {
            cout << "Need at least " << days << " entries for trend.\n";
            return;
        }
        vector<Entry> recent(entries.end() - days, entries.end());
        double avg = 0;
        for (auto& e : recent) avg += e.weight;
        avg /= recent.size();
        cout << "\n📈 Trend (last " << days << " days):\n";
        cout << "  Average weight: " << avg << " kg\n";
        cout << "  Recent entries:\n";
        int start = max(0, (int)recent.size() - 5);
        for (int i = start; i < (int)recent.size(); i++) {
            cout << "    " << recent[i].date << ": " << recent[i].weight << " kg\n";
        }
    }

    void exportCSV(const string& filename) {
        if (entries.empty()) {
            cout << "No entries.\n";
            return;
        }
        ofstream f(filename);
        f << "Date,Weight (kg),BMI,Note\n";
        for (auto& e : entries) {
            string bmiStr = e.bmi > 0 ? to_string(e.bmi).substr(0,4) : "";
            f << e.date << "," << e.weight << "," << bmiStr << "," << e.note << "\n";
        }
        f.close();
        cout << "✅ Exported " << entries.size() << " entries to " << filename << "\n";
    }
};

int main(int argc, char* argv[]) {
    if (argc < 2) {
        cerr << "Usage: pregnancy_weight <command> [options]\n";
        return 1;
    }
    PregnancyWeight app;
    string cmd = argv[1];

    if (cmd == "height") {
        if (argc < 3) { cerr << "height <cm>\n"; return 1; }
        app.setHeight(stod(argv[2]));
    } else if (cmd == "preweight") {
        if (argc < 3) { cerr << "preweight <kg>\n"; return 1; }
        app.setPreWeight(stod(argv[2]));
    } else if (cmd == "add") {
        if (argc < 3) { cerr << "add <weight> [--date DATE] [--note TEXT]\n"; return 1; }
        double weight = stod(argv[2]);
        string dateStr, note;
        for (int i=3; i<argc; i++) {
            if (string(argv[i]) == "--date" && i+1 < argc) dateStr = argv[++i];
            if (string(argv[i]) == "--note" && i+1 < argc) note = argv[++i];
        }
        app.add(weight, dateStr, note);
    } else if (cmd == "list") {
        app.list();
    } else if (cmd == "stats") {
        app.stats();
    } else if (cmd == "bmi") {
        app.showBMI();
    } else if (cmd == "gain") {
        app.showGain();
    } else if (cmd == "trend") {
        int days = 5;
        for (int i=2; i<argc; i++) {
            if (string(argv[i]) == "--days" && i+1 < argc) {
                days = stoi(argv[++i]);
            }
        }
        app.trend(days);
    } else if (cmd == "export") {
        string filename = argc > 2 ? argv[2] : "weights.csv";
        app.exportCSV(filename);
    } else {
        cerr << "Unknown command.\n";
        return 1;
    }
    return 0;
}
