package formatter

import (
    "fmt"
    "regexp"
    "strconv"
    "strings"
    "time"
)

// DateFormatter handles date and time formatting
type DateFormatter interface {
    ReplaceTimesThenDates(text string) string
    FormatTimeToken(token string, format string) string
    FormatDateToken(token string) string
}

type DateTimeProcessor struct{}

func NewDateFormatter() *DateTimeProcessor {
    return &DateTimeProcessor{}
}

func (f *DateTimeProcessor) ReplaceTimesThenDates(text string) string {
    // Process T12 and T24 first
    t12Re := regexp.MustCompile(`T12\(([^)]+)\)`)
    text = t12Re.ReplaceAllStringFunc(text, func(match string) string {
        return f.FormatTimeToken(match, "12h")
    })
    
    t24Re := regexp.MustCompile(`T24\(([^)]+)\)`)
    text = t24Re.ReplaceAllStringFunc(text, func(match string) string {
        return f.FormatTimeToken(match, "24h")
    })
    
    // Process D(...) last
    dRe := regexp.MustCompile(`D\(([^)]+)\)`)
    text = dRe.ReplaceAllStringFunc(text, func(match string) string {
        return f.FormatDateToken(match)
    })
    
    return text
}

func (f *DateTimeProcessor) FormatTimeToken(token string, format string) string {
    // Extract the ISO string inside parentheses
    var isoStr string
    if format == "12h" {
        re := regexp.MustCompile(`T12\(([^)]+)\)`)
        match := re.FindStringSubmatch(token)
        if len(match) < 2 {
            return token
        }
        isoStr = strings.TrimSpace(match[1])
    } else {
        re := regexp.MustCompile(`T24\(([^)]+)\)`)
        match := re.FindStringSubmatch(token)
        if len(match) < 2 {
            return token
        }
        isoStr = strings.TrimSpace(match[1])
    }
    
    // Validate the format first - offset must be in format ±HH:MM
    if !f.isValidTimeFormat(isoStr) {
        return token
    }
    
    // Parse the ISO time
    t, err := f.parseTime(isoStr)
    if err != nil {
        return token
    }
    
    // Format the offset
    offsetStr := f.formatOffset(isoStr, t)
    
    // Format the time according to specification
    if format == "12h" {
        return f.format12HourTime(t, offsetStr)
    } else {
        return f.format24HourTime(t, offsetStr)
    }
}

func (f *DateTimeProcessor) FormatDateToken(token string) string {
    re := regexp.MustCompile(`D\(([^)]+)\)`)
    match := re.FindStringSubmatch(token)
    if len(match) < 2 {
        return token
    }
    
    isoStr := strings.TrimSpace(match[1])
    t, err := f.parseDate(isoStr)
    if err != nil {
        return token
    }
    
    // Format as DD-Mmm-YYYY
    months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", 
        "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
    
    day := t.Day()
    month := months[t.Month()-1]
    year := t.Year()
    
    return fmt.Sprintf("%02d %s %d", day, month, year)
}

func (f *DateTimeProcessor) parseTime(isoStr string) (time.Time, error) {
    // Replace any non-standard minus characters
    isoStr = strings.ReplaceAll(isoStr, "−", "-")
    
    // Try different time formats
    formats := []string{
        "2006-01-02T15:04:05-07:00",
        "2006-01-02T15:04-07:00",
        "2006-01-02T15:04:05Z07:00",
        "2006-01-02T15:04Z07:00",
        "2006-01-02T15:04:05Z",
        "2006-01-02T15:04Z",
    }
    
    for _, format := range formats {
        t, err := time.Parse(format, isoStr)
        if err == nil {
            return t, nil
        }
    }
    
    return time.Time{}, fmt.Errorf("unable to parse time")
}

func (f *DateTimeProcessor) parseDate(isoStr string) (time.Time, error) {
    // Replace any non-standard minus characters
    isoStr = strings.ReplaceAll(isoStr, "−", "-")
    
    // Try different date formats
    formats := []string{
        "2006-01-02T15:04:05-07:00",
        "2006-01-02T15:04-07:00", 
        "2006-01-02T15:04:05Z",
        "2006-01-02T15:04Z",
        "2006-01-02",
    }
    
    for _, format := range formats {
        t, err := time.Parse(format, isoStr)
        if err == nil {
            return t, nil
        }
    }
    
    return time.Time{}, fmt.Errorf("unable to parse date")
}

func (f *DateTimeProcessor) isValidTimeFormat(isoStr string) bool {
    // Check for Zulu time
    if strings.HasSuffix(strings.ToUpper(isoStr), "Z") {
        return true
    }
    
    // Find the offset part (should be at the end)
    offsetIndex := strings.LastIndex(isoStr, "+")
    if offsetIndex == -1 {
        offsetIndex = strings.LastIndex(isoStr, "-")
        // Make sure it's not the first character (which would be part of the date)
        if offsetIndex <= 10 { // Date starts with YYYY-MM-DD
            return false
        }
    }
    
    if offsetIndex == -1 {
        return false
    }
    
    offsetStr := isoStr[offsetIndex:]
    
    // Offset should be in format ±HH:MM
    if len(offsetStr) != 6 {
        return false
    }
    
    // Check colon position
    if offsetStr[3] != ':' {
        return false
    }
    
    // Check that hours and minutes are digits
    hours := offsetStr[1:3]
    minutes := offsetStr[4:6]
    
    if _, err := strconv.Atoi(hours); err != nil {
        return false
    }
    if _, err := strconv.Atoi(minutes); err != nil {
        return false
    }
    
    // Validate hour range
    hourNum, _ := strconv.Atoi(hours)
    if hourNum < 0 || hourNum > 23 {
        return false
    }
    
    // Validate minute range  
    minuteNum, _ := strconv.Atoi(minutes)
    if minuteNum < 0 || minuteNum > 59 {
        return false
    }
    
    return true
}

func (f *DateTimeProcessor) formatOffset(isoStr string, t time.Time) string {
    _, offset := t.Zone()
    offsetHours := offset / 3600
    offsetMinutes := (offset % 3600) / 60
    
    // Handle Zulu time specifically
    if strings.HasSuffix(strings.ToUpper(isoStr), "Z") {
        return "(+00:00)"
    }
    
    // Format with proper sign
    sign := "+"
    if offsetHours < 0 {
        sign = "-"
        offsetHours = -offsetHours
    }
    
    return fmt.Sprintf("(%s%02d:%02d)", sign, offsetHours, offsetMinutes)
}

func (f *DateTimeProcessor) format12HourTime(t time.Time, offsetStr string) string {
    // 12-hour format with AM/PM
    hour := t.Hour()
    ampm := "AM"
    if hour >= 12 {
        ampm = "PM"
        if hour > 12 {
            hour -= 12
        }
    }
    if hour == 0 {
        hour = 12
    }
    return fmt.Sprintf("%02d:%02d%s %s", hour, t.Minute(), ampm, offsetStr)
}

func (f *DateTimeProcessor) format24HourTime(t time.Time, offsetStr string) string {
    // 24-hour format
    return fmt.Sprintf("%02d:%02d %s", t.Hour(), t.Minute(), offsetStr)
}