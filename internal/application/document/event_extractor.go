package documentapp

import (
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	domaincalendar "lexbox/internal/domain/calendar"
)

const (
	extractedDateKindAbsolute = "absolute"
	extractedDateKindRelative = "relative"

	anchorSourceInline               = "inline"
	anchorSourcePreviousLine         = "previous_line"
	anchorSourceNotificationLine     = "notification_line"
	anchorSourceProceduralAnchorLine = "procedural_anchor_line"
)

type extractedEventCandidate struct {
	EventType      string
	EventDate      string
	SourceText     string
	AnchorDate     string
	DateKind       string
	AnchorSource   string
	RelativeDays   int
	IsBusinessDays bool
	AddExtraDay    bool
	TriggerText    string
}

type relativeDateMatch struct {
	RawText        string
	Days           int
	IsBusinessDays bool
	AddExtraDay    bool
}

var (
	numericDateRegex = regexp.MustCompile(`\b(\d{1,2})[\/\-.](\d{1,2})[\/\-.](\d{4})\b`)
	textualDateRegex = regexp.MustCompile(`\b(\d{1,2})\s+de\s+([a-záéíóú]+)\s+de\s+(\d{4})\b`)

	relativeDaysRegex = regexp.MustCompile(
		`((?:en\s+el\s+plazo\s+de|en\s+plazo\s+de|plazo\s+de|dentro\s+de|en)\s+` +
			`(\d{1,2}|uno|una|dos|tres|cuatro|cinco|seis|siete|ocho|nueve|diez|once|doce|trece|catorce|quince|dieciseis|dieciséis|diecisiete|dieciocho|diecinueve|veinte)` +
			`\s+dias?` +
			`(?:\s+(habil(?:es)?|natural(?:es)?))?)`,
	)

	nextDayNotificationRegex = regexp.MustCompile(
		`(al\s+dia\s+siguiente\s+de\s+la\s+notificacion)`,
	)

	nextDayChainRegex = regexp.MustCompile(
		`(a\s+contar\s+desde\s+el\s+siguiente\s+dia\s+habil|` +
			`desde\s+el\s+siguiente\s+dia\s+habil|` +
			`a\s+partir\s+del\s+siguiente\s+dia\s+habil|` +
			`a\s+contar\s+desde\s+el\s+dia\s+habil\s+siguiente|` +
			`desde\s+el\s+dia\s+habil\s+siguiente|` +
			`a\s+partir\s+del\s+dia\s+habil\s+siguiente|` +
			`a\s+contar\s+desde\s+el\s+dia\s+siguiente|` +
			`desde\s+el\s+dia\s+siguiente|` +
			`a\s+partir\s+del\s+dia\s+siguiente|` +
			`a\s+contar\s+desde\s+el\s+siguiente|` +
			`desde\s+el\s+siguiente|` +
			`a\s+partir\s+del\s+siguiente)`,
	)
)

func extractDocumentEvents(content string, cfg EventComputationConfig) []extractedEventCandidate {
	text := normalizeEventText(content)
	if text == "" {
		return []extractedEventCandidate{}
	}

	lines := splitEventLines(content)
	candidates := make([]extractedEventCandidate, 0)

	var lastAbsoluteDate string
	var lastNotificationDate string
	var lastStrongProceduralDate string

	for _, line := range lines {
		lineNormalized := normalizeEventText(line)
		if lineNormalized == "" {
			continue
		}

		eventType := classifyEventType(lineNormalized)
		absoluteDates := extractDatesFromLine(line)

		if eventType == "unknown" {
			if len(absoluteDates) > 0 {
				lastAbsoluteDate = absoluteDates[len(absoluteDates)-1]
				if isStrongProceduralAnchorLine(lineNormalized) {
					lastStrongProceduralDate = absoluteDates[len(absoluteDates)-1]
				}
			}
			continue
		}

		for _, date := range absoluteDates {
			candidates = append(candidates, extractedEventCandidate{
				EventType:      eventType,
				EventDate:      date,
				SourceText:     strings.TrimSpace(line),
				AnchorDate:     date,
				DateKind:       extractedDateKindAbsolute,
				AnchorSource:   anchorSourceInline,
				RelativeDays:   0,
				IsBusinessDays: false,
				AddExtraDay:    false,
				TriggerText:    "",
			})
		}

		anchorDate, anchorSource := selectAnchorDate(
			eventType,
			absoluteDates,
			lastNotificationDate,
			lastStrongProceduralDate,
			lastAbsoluteDate,
		)

		relativeMatches := extractRelativeDateMatchesFromLine(line)
		for _, match := range relativeMatches {
			if strings.TrimSpace(anchorDate) == "" {
				continue
			}

			targetDate, ok := computeRelativeTargetDate(
				anchorDate,
				match.Days,
				match.IsBusinessDays,
				match.AddExtraDay,
				cfg,
			)
			if !ok {
				continue
			}

			candidates = append(candidates, extractedEventCandidate{
				EventType:      eventType,
				EventDate:      targetDate,
				SourceText:     strings.TrimSpace(line),
				AnchorDate:     anchorDate,
				DateKind:       extractedDateKindRelative,
				AnchorSource:   anchorSource,
				RelativeDays:   match.Days,
				IsBusinessDays: match.IsBusinessDays,
				AddExtraDay:    match.AddExtraDay,
				TriggerText:    match.RawText,
			})
		}

		if len(absoluteDates) > 0 {
			lastAbsoluteDate = absoluteDates[len(absoluteDates)-1]

			if eventType == "notification" {
				lastNotificationDate = absoluteDates[len(absoluteDates)-1]
			}

			if isStrongProceduralAnchorLine(lineNormalized) {
				lastStrongProceduralDate = absoluteDates[len(absoluteDates)-1]
			}
		}
	}

	return deduplicateEventCandidates(candidates)
}

func selectAnchorDate(
	eventType string,
	absoluteDates []string,
	lastNotificationDate string,
	lastStrongProceduralDate string,
	lastAbsoluteDate string,
) (string, string) {
	if len(absoluteDates) > 0 {
		return absoluteDates[len(absoluteDates)-1], anchorSourceInline
	}

	if usesNotificationAnchorPhrase(eventType) && strings.TrimSpace(lastNotificationDate) != "" {
		return lastNotificationDate, anchorSourceNotificationLine
	}

	if usesNotificationPreferredAnchor(eventType) && strings.TrimSpace(lastNotificationDate) != "" {
		return lastNotificationDate, anchorSourceNotificationLine
	}

	if strings.TrimSpace(lastStrongProceduralDate) != "" {
		return lastStrongProceduralDate, anchorSourceProceduralAnchorLine
	}

	if strings.TrimSpace(lastAbsoluteDate) != "" {
		return lastAbsoluteDate, anchorSourcePreviousLine
	}

	return "", ""
}

func usesNotificationPreferredAnchor(eventType string) bool {
	return eventType == "deadline" || eventType == "requirement"
}

func usesNotificationAnchorPhrase(eventType string) bool {
	return eventType == "deadline" || eventType == "requirement"
}

func isStrongProceduralAnchorLine(line string) bool {
	return containsAny(line,
		"resolucion de fecha",
		"resolución de fecha",
		"auto de fecha",
		"decreto de fecha",
		"providencia de fecha",
		"diligencia de ordenacion de fecha",
		"diligencia de ordenación de fecha",
	)
}

func splitEventLines(content string) []string {
	raw := strings.Split(content, "\n")
	lines := make([]string, 0, len(raw))
	for _, line := range raw {
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

func normalizeEventText(value string) string {
	return strings.ToLower(strings.Join(strings.Fields(value), " "))
}

func classifyEventType(line string) string {
	if looksLikeRelativeDeadline(line) {
		return "deadline"
	}

	switch {
	case containsAny(line,
		"juicio",
		"vista",
		"señalado para el día", "senalado para el dia",
		"se señala juicio", "se senala juicio",
		"se celebra el día", "se celebra el dia",
	):
		return "hearing"

	case containsAny(line,
		"comparecencia",
		"citación", "citacion",
		"cítese", "citese",
	):
		return "appearance"

	case containsAny(line,
		"plazo",
		"dentro del plazo",
		"en el plazo de",
		"se concede plazo",
		"hasta el",
		"improrrogable",
	):
		return "deadline"

	case containsAny(line,
		"presentación",
		"presentacion",
		"presentar escrito",
		"presentación de demanda", "presentacion de demanda",
		"presentación del recurso", "presentacion del recurso",
		"fecha de presentación", "fecha de presentacion",
	):
		return "filing"

	case containsAny(line,
		"requerir",
		"requerimiento",
		"requiérase", "requierase",
	):
		return "requirement"

	case containsAny(line,
		"notifíquese", "notifiquese",
		"notificación", "notificacion",
		"queda notificado",
	):
		return "notification"
	}

	return "unknown"
}

func looksLikeRelativeDeadline(line string) bool {
	normalized := normalizeASCIIText(line)

	if nextDayNotificationRegex.MatchString(normalized) {
		return true
	}

	if !relativeDaysRegex.MatchString(normalized) {
		return false
	}

	return containsAny(normalized,
		"debera",
		"debera aportar",
		"debera presentar",
		"debera formular",
		"aportar",
		"presentar",
		"formular alegaciones",
		"a contar desde la notificacion",
		"desde la notificacion",
		"a contar desde su notificacion",
		"desde su notificacion",
		"desde la notificacion de la presente resolucion",
		"a contar desde el dia siguiente",
		"desde el dia siguiente",
		"a contar desde el siguiente",
		"desde el siguiente",
		"a partir del dia siguiente",
		"a partir del siguiente",
		"a contar desde el siguiente dia habil",
		"desde el siguiente dia habil",
		"a partir del siguiente dia habil",
		"a contar desde el dia habil siguiente",
		"desde el dia habil siguiente",
		"a partir del dia habil siguiente",
	)
}

func extractDatesFromLine(line string) []string {
	dates := make([]string, 0)

	numericMatches := numericDateRegex.FindAllStringSubmatch(line, -1)
	for _, match := range numericMatches {
		if len(match) != 4 {
			continue
		}
		isoDate, ok := normalizeNumericDate(match[1], match[2], match[3])
		if ok {
			dates = append(dates, isoDate)
		}
	}

	textualMatches := textualDateRegex.FindAllStringSubmatch(strings.ToLower(line), -1)
	for _, match := range textualMatches {
		if len(match) != 4 {
			continue
		}
		isoDate, ok := normalizeTextualDate(match[1], match[2], match[3])
		if ok {
			dates = append(dates, isoDate)
		}
	}

	return uniqueDates(dates)
}

func extractRelativeDateMatchesFromLine(line string) []relativeDateMatch {
	normalized := normalizeASCIIText(line)

	result := make([]relativeDateMatch, 0)

	if nextDayMatches := nextDayNotificationRegex.FindAllStringSubmatch(normalized, -1); len(nextDayMatches) > 0 {
		for _, match := range nextDayMatches {
			if len(match) < 2 {
				continue
			}
			result = append(result, relativeDateMatch{
				RawText:        strings.TrimSpace(match[1]),
				Days:           1,
				IsBusinessDays: false,
				AddExtraDay:    false,
			})
		}
	}

	matches := relativeDaysRegex.FindAllStringSubmatch(normalized, -1)
	addExtraDay := nextDayChainRegex.MatchString(normalized)

	for _, match := range matches {
		if len(match) < 3 {
			continue
		}

		days, ok := parseSpanishCardinal(normalizeASCIIText(match[2]))
		if !ok || days <= 0 {
			continue
		}

		modifier := ""
		if len(match) >= 4 {
			modifier = strings.TrimSpace(normalizeASCIIText(match[3]))
		}

		isBusinessDays := modifier == "habil" || modifier == "habiles"

		rawText := buildRelativeTriggerText(normalized, strings.TrimSpace(match[1]), addExtraDay)

		result = append(result, relativeDateMatch{
			RawText:        rawText,
			Days:           days,
			IsBusinessDays: isBusinessDays,
			AddExtraDay:    addExtraDay,
		})
	}

	return deduplicateRelativeMatches(result)
}

func buildRelativeTriggerText(normalizedLine, baseTrigger string, addExtraDay bool) string {
	baseTrigger = strings.TrimSpace(baseTrigger)
	if baseTrigger == "" {
		return ""
	}

	if !addExtraDay {
		return baseTrigger
	}

	nextDayText := detectNextDayTriggerText(normalizedLine)
	if nextDayText == "" {
		return baseTrigger
	}

	return strings.TrimSpace(baseTrigger + " " + nextDayText)
}

func detectNextDayTriggerText(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}

	phrases := []string{
		"a contar desde el siguiente dia habil",
		"desde el siguiente dia habil",
		"a partir del siguiente dia habil",
		"a contar desde el dia habil siguiente",
		"desde el dia habil siguiente",
		"a partir del dia habil siguiente",
		"a contar desde el dia siguiente",
		"desde el dia siguiente",
		"a partir del dia siguiente",
		"a contar desde el siguiente",
		"desde el siguiente",
		"a partir del siguiente",
	}

	for _, phrase := range phrases {
		if strings.Contains(value, phrase) {
			return phrase
		}
	}

	return ""
}

func deduplicateRelativeMatches(items []relativeDateMatch) []relativeDateMatch {
	type key struct {
		rawText        string
		days           int
		isBusinessDays bool
		addExtraDay    bool
	}

	seen := make(map[key]struct{}, len(items))
	result := make([]relativeDateMatch, 0, len(items))

	for _, item := range items {
		k := key{
			rawText:        item.RawText,
			days:           item.Days,
			isBusinessDays: item.IsBusinessDays,
			addExtraDay:    item.AddExtraDay,
		}
		if _, exists := seen[k]; exists {
			continue
		}
		seen[k] = struct{}{}
		result = append(result, item)
	}

	return result
}

func calendarForRange(scope string, startYear, endYear int) domaincalendar.Calendar {
	if endYear < startYear {
		endYear = startYear
	}

	scope = strings.TrimSpace(scope)
	if scope == "" {
		scope = domaincalendar.ScopeMadrid
	}

	holidays := make([]string, 0)
	for year := startYear; year <= endYear; year++ {
		holidays = append(holidays, domaincalendar.Holidays(scope, year)...)
	}

	return domaincalendar.NewCalendar(holidays)
}

func computeRelativeTargetDate(
	anchorDate string,
	days int,
	isBusinessDays bool,
	addExtraDay bool,
	cfg EventComputationConfig,
) (string, bool) {
	if strings.TrimSpace(anchorDate) == "" || days <= 0 {
		return "", false
	}

	baseDate, err := time.Parse("2006-01-02", anchorDate)
	if err != nil {
		return "", false
	}

	totalDays := days
	if addExtraDay {
		totalDays++
	}

	var target time.Time
	if isBusinessDays {
		endYear := baseDate.Year()
		if totalDays > 0 {
			endYear = baseDate.AddDate(0, 0, totalDays+370).Year()
		}

		cal := calendarForRange(cfg.CalendarScope, baseDate.Year(), endYear)
		target = cal.AddBusinessDaysWithRules(baseDate, totalDays, cfg.ProceduralRules)
	} else {
		target = baseDate.AddDate(0, 0, totalDays)
	}

	return target.Format("2006-01-02"), true
}

func normalizeNumericDate(day, month, year string) (string, bool) {
	value := strings.TrimSpace(day) + "/" + strings.TrimSpace(month) + "/" + strings.TrimSpace(year)
	t, err := time.Parse("2/1/2006", value)
	if err != nil {
		return "", false
	}
	return t.Format("2006-01-02"), true
}

func normalizeTextualDate(day, monthName, year string) (string, bool) {
	month, ok := spanishMonthToNumber(monthName)
	if !ok {
		return "", false
	}

	dayInt, err := strconv.Atoi(strings.TrimSpace(day))
	if err != nil {
		return "", false
	}

	yearInt, err := strconv.Atoi(strings.TrimSpace(year))
	if err != nil {
		return "", false
	}

	t := time.Date(yearInt, time.Month(month), dayInt, 0, 0, 0, 0, time.UTC)
	if t.Day() != dayInt || int(t.Month()) != month || t.Year() != yearInt {
		return "", false
	}

	return t.Format("2006-01-02"), true
}

func spanishMonthToNumber(month string) (int, bool) {
	month = strings.ToLower(strings.TrimSpace(month))

	months := map[string]int{
		"enero":      1,
		"febrero":    2,
		"marzo":      3,
		"abril":      4,
		"mayo":       5,
		"junio":      6,
		"julio":      7,
		"agosto":     8,
		"septiembre": 9,
		"setiembre":  9,
		"octubre":    10,
		"noviembre":  11,
		"diciembre":  12,
	}

	value, ok := months[month]
	return value, ok
}

func parseSpanishCardinal(value string) (int, bool) {
	value = strings.ToLower(strings.TrimSpace(value))

	if n, err := strconv.Atoi(value); err == nil {
		return n, true
	}

	numbers := map[string]int{
		"uno":        1,
		"una":        1,
		"dos":        2,
		"tres":       3,
		"cuatro":     4,
		"cinco":      5,
		"seis":       6,
		"siete":      7,
		"ocho":       8,
		"nueve":      9,
		"diez":       10,
		"once":       11,
		"doce":       12,
		"trece":      13,
		"catorce":    14,
		"quince":     15,
		"dieciseis":  16,
		"dieciséis":  16,
		"diecisiete": 17,
		"dieciocho":  18,
		"diecinueve": 19,
		"veinte":     20,
	}

	n, ok := numbers[value]
	return n, ok
}

func deduplicateEventCandidates(items []extractedEventCandidate) []extractedEventCandidate {
	type key struct {
		EventType  string
		EventDate  string
		SourceText string
	}

	seen := make(map[key]struct{}, len(items))
	result := make([]extractedEventCandidate, 0, len(items))

	for _, item := range items {
		k := key{
			EventType:  item.EventType,
			EventDate:  item.EventDate,
			SourceText: item.SourceText,
		}

		if _, exists := seen[k]; exists {
			continue
		}

		seen[k] = struct{}{}
		result = append(result, item)
	}

	return result
}

func uniqueDates(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))

	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}

	return result
}

func normalizeASCIIText(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))

	var b strings.Builder
	b.Grow(len(value))

	for _, r := range value {
		switch r {
		case 'á', 'à', 'ä', 'â':
			b.WriteRune('a')
		case 'é', 'è', 'ë', 'ê':
			b.WriteRune('e')
		case 'í', 'ì', 'ï', 'î':
			b.WriteRune('i')
		case 'ó', 'ò', 'ö', 'ô':
			b.WriteRune('o')
		case 'ú', 'ù', 'ü', 'û':
			b.WriteRune('u')
		case 'ñ':
			b.WriteRune('n')
		default:
			if unicode.IsSpace(r) {
				b.WriteRune(' ')
			} else {
				b.WriteRune(r)
			}
		}
	}

	return strings.Join(strings.Fields(b.String()), " ")
}
