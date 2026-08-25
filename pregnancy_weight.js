// pregnancy_weight.js
#!/usr/bin/env node
const fs = require('fs');
const { program } = require('commander');

const DATA_FILE = 'pregnancy_weight.json';

class PregnancyWeight {
    constructor() {
        this.data = this.load();
        this.height = this.data.height || null;
        this.preWeight = this.data.preWeight || null;
        this.entries = this.data.entries || [];
    }

    load() {
        if (fs.existsSync(DATA_FILE)) {
            try {
                return JSON.parse(fs.readFileSync(DATA_FILE));
            } catch (e) {}
        }
        return { height: null, preWeight: null, entries: [] };
    }

    save() {
        fs.writeFileSync(DATA_FILE, JSON.stringify({
            height: this.height,
            preWeight: this.preWeight,
            entries: this.entries
        }, null, 2));
    }

    setHeight(cm) {
        this.height = cm;
        this.save();
        console.log(`✅ Height set to ${cm} cm`);
    }

    setPreWeight(kg) {
        this.preWeight = kg;
        this.save();
        console.log(`✅ Pre‑pregnancy weight set to ${kg} kg`);
    }

    bmi(weight) {
        if (!this.height) return null;
        const h = this.height / 100;
        return weight / (h * h);
    }

    bmiCategory(bmiVal) {
        if (bmiVal < 18.5) return 'Underweight';
        if (bmiVal < 25.0) return 'Normal';
        if (bmiVal < 30.0) return 'Overweight';
        return 'Obese';
    }

    add(weight, dateStr, note) {
        if (!this.height) {
            console.log('❌ Please set your height first: height <cm>');
            return;
        }
        if (!dateStr) dateStr = new Date().toISOString().slice(0,10);
        const bmiVal = this.bmi(weight);
        const entry = { date: dateStr, weight, note: note || '', bmi: bmiVal };
        this.entries.push(entry);
        this.entries.sort((a, b) => a.date.localeCompare(b.date));
        this.save();
        console.log(`✅ Logged ${weight} kg on ${dateStr}`);
    }

    list() {
        if (this.entries.length === 0) {
            console.log('No entries.');
            return;
        }
        console.log('\n📋 Weight Entries:');
        for (const e of this.entries) {
            const note = e.note ? ` – ${e.note}` : '';
            const bmiStr = e.bmi ? ` (BMI: ${e.bmi.toFixed(1)})` : '';
            console.log(`  ${e.date}: ${e.weight.toFixed(1)} kg${bmiStr}${note}`);
        }
    }

    stats() {
        if (this.entries.length === 0) {
            console.log('No entries.');
            return;
        }
        const weights = this.entries.map(e => e.weight);
        const total = weights.length;
        const current = weights[weights.length - 1];
        const avg = weights.reduce((a, b) => a + b, 0) / total;
        const mn = Math.min(...weights);
        const mx = Math.max(...weights);
        const gain = this.preWeight ? current - this.preWeight : null;

        console.log('\n📊 Statistics:');
        console.log(`  Entries: ${total}`);
        console.log(`  Current: ${current.toFixed(1)} kg`);
        console.log(`  Average: ${avg.toFixed(1)} kg`);
        console.log(`  Min: ${mn.toFixed(1)} kg`);
        console.log(`  Max: ${mx.toFixed(1)} kg`);
        if (gain !== null) {
            console.log(`  Total gain: ${gain.toFixed(1)} kg`);
        }
    }

    showBMI() {
        if (this.entries.length === 0) {
            console.log('No entries.');
            return;
        }
        if (!this.height) {
            console.log('❌ Height not set.');
            return;
        }
        const current = this.entries[this.entries.length - 1].weight;
        const bmiVal = this.bmi(current);
        const cat = this.bmiCategory(bmiVal);
        console.log(`\n📐 Current BMI: ${bmiVal.toFixed(1)} (${cat})`);
        console.log('BMI Range: 18.5 - 24.9');
    }

    showGain() {
        if (this.entries.length === 0) {
            console.log('No entries.');
            return;
        }
        if (!this.preWeight) {
            console.log('❌ Pre‑pregnancy weight not set. Use: preweight <kg>');
            return;
        }
        const current = this.entries[this.entries.length - 1].weight;
        const first = this.entries[0].weight;
        const gain = current - this.preWeight;
        console.log('\n📈 Weight Gain:');
        console.log(`  Pre‑pregnancy: ${this.preWeight.toFixed(1)} kg`);
        console.log(`  First logged: ${first.toFixed(1)} kg`);
        console.log(`  Current: ${current.toFixed(1)} kg`);
        console.log(`  Total gain: ${gain.toFixed(1)} kg`);
    }

    trend(days = 5) {
        if (this.entries.length < days) {
            console.log(`Need at least ${days} entries for trend.`);
            return;
        }
        const recent = this.entries.slice(-days);
        const avg = recent.reduce((a, b) => a + b.weight, 0) / recent.length;
        console.log(`\n📈 Trend (last ${days} days):`);
        console.log(`  Average weight: ${avg.toFixed(1)} kg`);
        console.log('  Recent entries:');
        for (const e of recent.slice(-5)) {
            console.log(`    ${e.date}: ${e.weight.toFixed(1)} kg`);
        }
    }

    export(filename = 'weights.csv') {
        if (this.entries.length === 0) {
            console.log('No entries.');
            return;
        }
        let csv = 'Date,Weight (kg),BMI,Note\n';
        for (const e of this.entries) {
            const bmiStr = e.bmi ? e.bmi.toFixed(1) : '';
            csv += `${e.date},${e.weight.toFixed(1)},${bmiStr},${e.note}\n`;
        }
        fs.writeFileSync(filename, csv);
        console.log(`✅ Exported ${this.entries.length} entries to ${filename}`);
    }
}

program
    .command('height <cm>')
    .action((cm) => {
        const app = new PregnancyWeight();
        app.setHeight(parseFloat(cm));
    });

program
    .command('preweight <kg>')
    .action((kg) => {
        const app = new PregnancyWeight();
        app.setPreWeight(parseFloat(kg));
    });

program
    .command('add <weight>')
    .option('--date <date>', 'YYYY-MM-DD')
    .option('--note <note>', 'Optional note')
    .action((weight, options) => {
        const app = new PregnancyWeight();
        app.add(parseFloat(weight), options.date, options.note);
    });

program
    .command('list')
    .action(() => {
        const app = new PregnancyWeight();
        app.list();
    });

program
    .command('stats')
    .action(() => {
        const app = new PregnancyWeight();
        app.stats();
    });

program
    .command('bmi')
    .action(() => {
        const app = new PregnancyWeight();
        app.showBMI();
    });

program
    .command('gain')
    .action(() => {
        const app = new PregnancyWeight();
        app.showGain();
    });

program
    .command('trend')
    .option('--days <n>', 'Number of days', parseInt, 5)
    .action((options) => {
        const app = new PregnancyWeight();
        app.trend(options.days);
    });

program
    .command('export')
    .argument('[filename]', 'Output filename', 'weights.csv')
    .action((filename) => {
        const app = new PregnancyWeight();
        app.export(filename);
    });

program.parse(process.argv);
