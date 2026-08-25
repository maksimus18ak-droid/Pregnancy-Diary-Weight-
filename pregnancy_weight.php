# pregnancy_weight.php
#!/usr/bin/env php
<?php

define('DATA_FILE', 'pregnancy_weight.json');

class PregnancyWeight {
    private $height;
    private $preWeight;
    private $entries = [];

    public function __construct() {
        $this->load();
    }

    private function load() {
        if (file_exists(DATA_FILE)) {
            $data = json_decode(file_get_contents(DATA_FILE), true);
            $this->height = $data['height'] ?? null;
            $this->preWeight = $data['pre_weight'] ?? null;
            $this->entries = $data['entries'] ?? [];
        }
    }

    private function save() {
        file_put_contents(DATA_FILE, json_encode([
            'height' => $this->height,
            'pre_weight' => $this->preWeight,
            'entries' => $this->entries
        ], JSON_PRETTY_PRINT));
    }

    public function setHeight($cm) {
        $this->height = $cm;
        $this->save();
        echo "✅ Height set to {$cm} cm\n";
    }

    public function setPreWeight($kg) {
        $this->preWeight = $kg;
        $this->save();
        echo "✅ Pre‑pregnancy weight set to {$kg} kg\n";
    }

    private function bmi($weight) {
        if (!$this->height) return null;
        $h = $this->height / 100.0;
        return $weight / ($h * $h);
    }

    private function bmiCategory($bmiVal) {
        if ($bmiVal < 18.5) return "Underweight";
        if ($bmiVal < 25.0) return "Normal";
        if ($bmiVal < 30.0) return "Overweight";
        return "Obese";
    }

    public function add($weight, $dateStr = null, $note = '') {
        if (!$this->height) {
            echo "❌ Please set your height first: height <cm>\n";
            return;
        }
        $dateStr = $dateStr ?: date('Y-m-d');
        $bmiVal = $this->bmi($weight);
        $entry = [
            'date' => $dateStr,
            'weight' => $weight,
            'note' => $note,
            'bmi' => $bmiVal
        ];
        $this->entries[] = $entry;
        usort($this->entries, function($a, $b) {
            return strcmp($a['date'], $b['date']);
        });
        $this->save();
        echo "✅ Logged {$weight} kg on {$dateStr}\n";
    }

    public function list() {
        if (empty($this->entries)) {
            echo "No entries.\n";
            return;
        }
        echo "\n📋 Weight Entries:\n";
        foreach ($this->entries as $e) {
            $note = $e['note'] ? " – {$e['note']}" : '';
            $bmiStr = $e['bmi'] ? " (BMI: " . round($e['bmi'], 1) . ")" : '';
            echo "  {$e['date']}: " . round($e['weight'], 1) . " kg{$bmiStr}{$note}\n";
        }
    }

    public function stats() {
        if (empty($this->entries)) {
            echo "No entries.\n";
            return;
        }
        $weights = array_column($this->entries, 'weight');
        $total = count($weights);
        $current = $weights[$total - 1];
        $avg = array_sum($weights) / $total;
        $mn = min($weights);
        $mx = max($weights);
        $gain = $this->preWeight ? $current - $this->preWeight : null;

        echo "\n📊 Statistics:\n";
        echo "  Entries: {$total}\n";
        echo "  Current: " . round($current, 1) . " kg\n";
        echo "  Average: " . round($avg, 1) . " kg\n";
        echo "  Min: " . round($mn, 1) . " kg\n";
        echo "  Max: " . round($mx, 1) . " kg\n";
        if ($gain !== null) {
            echo "  Total gain: " . round($gain, 1) . " kg\n";
        }
    }

    public function showBMI() {
        if (empty($this->entries)) {
            echo "No entries.\n";
            return;
        }
        if (!$this->height) {
            echo "❌ Height not set.\n";
            return;
        }
        $current = $this->entries[count($this->entries)-1]['weight'];
        $bmiVal = $this->bmi($current);
        $cat = $this->bmiCategory($bmiVal);
        echo "\n📐 Current BMI: " . round($bmiVal, 1) . " ({$cat})\n";
        echo "BMI Range: 18.5 - 24.9\n";
    }

    public function showGain() {
        if (empty($this->entries)) {
            echo "No entries.\n";
            return;
        }
        if (!$this->preWeight) {
            echo "❌ Pre‑pregnancy weight not set. Use: preweight <kg>\n";
            return;
        }
        $current = $this->entries[count($this->entries)-1]['weight'];
        $first = $this->entries[0]['weight'];
        $gain = $current - $this->preWeight;
        echo "\n📈 Weight Gain:\n";
        echo "  Pre‑pregnancy: " . round($this->preWeight, 1) . " kg\n";
        echo "  First logged: " . round($first, 1) . " kg\n";
        echo "  Current: " . round($current, 1) . " kg\n";
        echo "  Total gain: " . round($gain, 1) . " kg\n";
    }

    public function trend($days = 5) {
        if (count($this->entries) < $days) {
            echo "Need at least {$days} entries for trend.\n";
            return;
        }
        $recent = array_slice($this->entries, -$days);
        $avg = array_sum(array_column($recent, 'weight')) / count($recent);
        echo "\n📈 Trend (last {$days} days):\n";
        echo "  Average weight: " . round($avg, 1) . " kg\n";
        echo "  Recent entries:\n";
        foreach (array_slice($recent, -5) as $e) {
            echo "    {$e['date']}: " . round($e['weight'], 1) . " kg\n";
        }
    }

    public function export($filename = 'weights.csv') {
        if (empty($this->entries)) {
            echo "No entries.\n";
            return;
        }
        $fp = fopen($filename, 'w');
        fputcsv($fp, ['Date', 'Weight (kg)', 'BMI', 'Note']);
        foreach ($this->entries as $e) {
            $bmiStr = $e['bmi'] ? round($e['bmi'], 1) : '';
            fputcsv($fp, [$e['date'], round($e['weight'], 1), $bmiStr, $e['note']]);
        }
        fclose($fp);
        echo "✅ Exported " . count($this->entries) . " entries to {$filename}\n";
    }
}

if ($argc < 2) {
    die("Usage: php pregnancy_weight.php <command> [options]\n");
}
$app = new PregnancyWeight();
$cmd = $argv[1];

switch ($cmd) {
    case 'height':
        if ($argc < 3) die("height <cm>\n");
        $app->setHeight((float)$argv[2]);
        break;
    case 'preweight':
        if ($argc < 3) die("preweight <kg>\n");
        $app->setPreWeight((float)$argv[2]);
        break;
    case 'add':
        if ($argc < 3) die("add <weight> [--date DATE] [--note TEXT]\n");
        $weight = (float)$argv[2];
        $date = null;
        $note = '';
        for ($i=3; $i<$argc; $i++) {
            if ($argv[$i] == '--date' && isset($argv[$i+1])) {
                $date = $argv[++$i];
            }
            if ($argv[$i] == '--note' && isset($argv[$i+1])) {
                $note = $argv[++$i];
            }
        }
        $app->add($weight, $date, $note);
        break;
    case 'list':
        $app->list();
        break;
    case 'stats':
        $app->stats();
        break;
    case 'bmi':
        $app->showBMI();
        break;
    case 'gain':
        $app->showGain();
        break;
    case 'trend':
        $days = 5;
        for ($i=2; $i<$argc; $i++) {
            if ($argv[$i] == '--days' && isset($argv[$i+1])) {
                $days = (int)$argv[++$i];
            }
        }
        $app->trend($days);
        break;
    case 'export':
        $filename = $argv[2] ?? 'weights.csv';
        $app->export($filename);
        break;
    default:
        echo "Unknown command.\n";
}
?>
