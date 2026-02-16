
---

````markdown
# 🚀 Itinerary-Prettifier – Live Demo Cheat Sheet

## Setup
```bash
cd ~/kood_sisu/itinerary-prettifier
ls
````

Expected → airports/, cli/, config/, fileio/, formatter/, types/, input.txt, airport-lookup.csv, main.go

---

## 1️⃣ Help flag

```bash
go run . -h
```

Expected:

```
itinerary usage:
go run . ./input.txt ./output.txt ./airport-lookup.csv
```

✅ Say: Shows correct CLI usage.

---

## 2️⃣ Error Handling

**Missing input**

```bash
go run . ./missing.txt ./output.txt ./airport-lookup.csv
```

→ `Input not found`

**Missing lookup**

```bash
go run . ./input.txt ./output.txt ./missing.csv
```

→ `Airport lookup not found`

**Malformed lookup**
(delete a CSV column)

```bash
go run . ./input.txt ./output.txt ./airport-lookup.csv
```

→ `Airport lookup malformed`

---

## 3️⃣ Output Safety

```bash
echo "Do not overwrite" > output.txt
go run . ./missing.txt ./output.txt ./airport-lookup.csv
cat output.txt
```

✅ Output file unchanged

---

## 4️⃣ Successful Run

```bash
go run . ./input.txt ./output.txt ./airport-lookup.csv
cat output.txt
```

Expected:

```
Flight from Los Angeles International Airport to London Heathrow Airport
Departure 05 Jun 2024
Arrival 09:00AM (+01:00)
```

---

## 5️⃣ Date & Time Formatting

Input contains ISO timestamps
→ Output shows friendly dates/times
Example:
`D(2022-05-09T08:07Z)` → `09 May 2022`

---

## 6️⃣ Whitespace Handling

Control chars `\v \f \r` → newline `\n`
Consecutive blank lines ≤ 1

---

## 7️⃣ Airport & City Lookup

```bash
Your flight departs from #HAJ, destination ##EDDW
Your city of departure is *#HAJ
```

Output:

```
Hannover Airport … Bremen Airport
Your city of departure is Hannover
```

---

## 8️⃣ Run Tests

```bash
go test ./... -v
```

Expected:

```
=== RUN   TestParseAirportCSV
--- PASS: TestParseAirportCSV
=== RUN   TestFormatterPrettify
--- PASS: TestFormatterPrettify
PASS
```

---

## 🧠 Quick Script

> “Tool validates inputs, formats data, expands airport codes, normalizes whitespace, and never overwrites output on error. Fully tested and reliable.”

```
---

