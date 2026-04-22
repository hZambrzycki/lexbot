package documentapp

import (
	"strings"

	domaincalendar "lexbox/internal/domain/calendar"
	domaincasefile "lexbox/internal/domain/casefile"
)

func EventComputationConfigFromCaseFile(cf domaincasefile.CaseFile) EventComputationConfig {
	scope := strings.TrimSpace(cf.CalendarScope)
	if scope == "" {
		scope = domaincalendar.ScopeMadrid
	}

	return EventComputationConfig{
		CalendarScope: scope,
		ProceduralRules: domaincalendar.ProceduralRules{
			AugustNonBusiness: cf.AugustNonBusiness,
		},
	}
}
