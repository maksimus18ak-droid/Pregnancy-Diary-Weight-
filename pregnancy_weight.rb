# pregnancy_weight.rb
#!/usr/bin/env ruby
require 'json'
require 'date'
require 'optparse'
require 'csv'

DATA_FILE = 'pregnancy_weight.json'

class PregnancyWeight
  attr_reader :height, :pre_weight, :entries

  def initialize
    @height = nil
    @pre_weight = nil
    @entries = []
    load
  end

  def load
    if File.exist?(DATA_FILE)
      data = JSON.parse(File.read(DATA_FILE))
      @height = data['height']
      @pre_weight = data['pre_weight']
      @entries = data['entries'] || []
    end
  end

  def save
    File.write(DATA_FILE, JSON.pretty_generate({
      height: @height,
      pre_weight: @pre_weight,
      entries: @entries
    }))
  end

  def set_height(cm)
    @height = cm
    save
    puts "✅ Height set to #{cm} cm"
  end

  def set_pre_weight(kg)
    @pre_weight = kg
    save
    puts "✅ Pre‑pregnancy weight set to #{kg} kg"
  end

  def bmi(weight)
    return nil unless @height
    h = @height / 100.0
    weight / (h * h)
  end

  def bmi_category(bmi_val)
    if bmi_val < 18.5
      "Underweight"
    elsif bmi_val < 25.0
      "Normal"
    elsif bmi_val < 30.0
      "Overweight"
    else
      "Obese"
    end
  end

  def add(weight, date_str = nil, note = '')
    unless @height
      puts "❌ Please set your height first: height <cm>"
      return
    end
    date_str ||= Date.today.to_s
    bmi_val = bmi(weight)
    entry = { 'date' => date_str, 'weight' => weight, 'note' => note, 'bmi' => bmi_val }
    @entries << entry
    @entries.sort_by! { |e| e['date'] }
    save
    puts "✅ Logged #{weight} kg on #{date_str}"
  end

  def list
    if @entries.empty?
      puts "No entries."
      return
    end
    puts "\n📋 Weight Entries:"
    @entries.each do |e|
      note = e['note'].empty? ? '' : " – #{e['note']}"
      bmi_str = e['bmi'] ? " (BMI: #{e['bmi'].round(1)})" : ''
      puts "  #{e['date']}: #{e['weight'].round(1)} kg#{bmi_str}#{note}"
    end
  end

  def stats
    if @entries.empty?
      puts "No entries."
      return
    end
    weights = @entries.map { |e| e['weight'] }
    total = weights.size
    current = weights.last
    avg = weights.sum / total
    mn = weights.min
    mx = weights.max
    gain = @pre_weight ? current - @pre_weight : nil

    puts "\n📊 Statistics:"
    puts "  Entries: #{total}"
    puts "  Current: #{current.round(1)} kg"
    puts "  Average: #{avg.round(1)} kg"
    puts "  Min: #{mn.round(1)} kg"
    puts "  Max: #{mx.round(1)} kg"
    puts "  Total gain: #{gain.round(1)} kg" if gain
  end

  def show_bmi
    if @entries.empty?
      puts "No entries."
      return
    end
    unless @height
      puts "❌ Height not set."
      return
    end
    current = @entries.last['weight']
    bmi_val = bmi(current)
    cat = bmi_category(bmi_val)
    puts "\n📐 Current BMI: #{bmi_val.round(1)} (#{cat})"
    puts "BMI Range: 18.5 - 24.9"
  end

  def show_gain
    if @entries.empty?
      puts "No entries."
      return
    end
    unless @pre_weight
      puts "❌ Pre‑pregnancy weight not set. Use: preweight <kg>"
      return
    end
    current = @entries.last['weight']
    first = @entries.first['weight']
    gain = current - @pre_weight
    puts "\n📈 Weight Gain:"
    puts "  Pre‑pregnancy: #{@pre_weight.round(1)} kg"
    puts "  First logged: #{first.round(1)} kg"
    puts "  Current: #{current.round(1)} kg"
    puts "  Total gain: #{gain.round(1)} kg"
  end

  def trend(days = 5)
    if @entries.size < days
      puts "Need at least #{days} entries for trend."
      return
    end
    recent = @entries.last(days)
    avg = recent.sum { |e| e['weight'] } / recent.size
    puts "\n📈 Trend (last #{days} days):"
    puts "  Average weight: #{avg.round(1)} kg"
    puts "  Recent entries:"
    recent.last(5).each do |e|
      puts "    #{e['date']}: #{e['weight'].round(1)} kg"
    end
  end

  def export(filename = 'weights.csv')
    if @entries.empty?
      puts "No entries."
      return
    end
    CSV.open(filename, 'w') do |csv|
      csv << ['Date', 'Weight (kg)', 'BMI', 'Note']
      @entries.each do |e|
        bmi_str = e['bmi'] ? e['bmi'].round(1).to_s : ''
        csv << [e['date'], e['weight'].round(1), bmi_str, e['note']]
      end
    end
    puts "✅ Exported #{@entries.size} entries to #{filename}"
  end
end

options = {}
OptionParser.new do |opts|
  opts.banner = "Usage: pregnancy_weight.rb <command> [options]"
  opts.on("height CM", Float, "Set height") { |v| options[:height] = v }
  opts.on("preweight KG", Float, "Set pre‑pregnancy weight") { |v| options[:preweight] = v }
  opts.on("add WEIGHT", Float, "Log weight") { |v| options[:add] = v }
  opts.on("--date DATE", "YYYY-MM-DD") { |v| options[:date] = v }
  opts.on("--note TEXT", "Optional note") { |v| options[:note] = v }
  opts.on("list", "List entries") { options[:list] = true }
  opts.on("stats", "Show statistics") { options[:stats] = true }
  opts.on("bmi", "Show BMI") { options[:bmi] = true }
  opts.on("gain", "Show weight gain") { options[:gain] = true }
  opts.on("trend", "Show trend") { options[:trend] = true }
  opts.on("--days N", Integer, "Days for trend") { |v| options[:trend_days] = v }
  opts.on("export", "Export to CSV") { options[:export] = true }
  opts.on("--filename FILE", "Output filename") { |v| options[:filename] = v }
end.parse!

app = PregnancyWeight.new

if options[:height]
  app.set_height(options[:height])
elsif options[:preweight]
  app.set_pre_weight(options[:preweight])
elsif options[:add]
  app.add(options[:add], options[:date], options[:note] || '')
elsif options[:list]
  app.list
elsif options[:stats]
  app.stats
elsif options[:bmi]
  app.show_bmi
elsif options[:gain]
  app.show_gain
elsif options[:trend]
  app.trend(options[:trend_days] || 5)
elsif options[:export]
  app.export(options[:filename] || 'weights.csv')
else
  puts "Usage: pregnancy_weight.rb <command> [options]"
end
