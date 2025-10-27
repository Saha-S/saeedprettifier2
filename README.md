<style>
.title { 
  font-size: 2.5em; 
  font-weight: bold; 
  text-align: center; 
  background: linear-gradient(90deg, #B18EFF, #E6FF00, #FFDB58); 
  background-clip: text; 
  -webkit-background-clip: text; 
  color: transparent; 
  animation: gradientFlow 5s infinite linear; 
}
@keyframes gradientFlow { 0% { background-position: 0% 50%; } 100% { background-position: 100% 50%; } }
.subtitle { 
  text-align: center; 
  color: #E6FF00; 
  font-style: italic; 
}
</style>

<div class="title">✈ KOOD AIRLINES PRETTIFIER ✈</div>
<div class="subtitle">Make your flight itineraries customer-friendly!</div>

---

## 🚀 Features
- Convert ISO dates `D(...)` → `DD MMM YYYY` (e.g., `01 Nov 2025`)
- Convert times `T12(...)` → `09:30PM (-02:00)` or `T24(...)` → `21:30 (+02:00)`
- Swap airport codes `#LAX` / `##EGLL` → airport **names** or `*#LHR` → **city names**
- Trim vertical whitespace (`\v`, `\f`, `\r`) → single `\n`
- Collapse multiple blank lines → **one blank line max**
- Optional terminal highlights with `-color` 💜
- Optional city display with `-city` 🏙️
- Random travel inspiration with `-getlucky` 🎲

---

## ⚡ Usage
```bash
# Basic
go run . ./input.txt ./output.txt ./airport-lookup.csv

# Show city names
go run . ./input.txt ./output.txt ./airport-lookup.csv -city

# Add KOOD violet highlights
go run . ./input.txt ./output.txt ./airport-lookup.csv -color

# Both city + color
go run . ./input.txt ./output.txt ./airport-lookup.csv -city -color

# Random travel inspiration
go run . -getlucky
