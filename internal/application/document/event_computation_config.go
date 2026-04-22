package documentapp

import domaincalendar "lexbox/internal/domain/calendar"

type EventComputationConfig struct {
	CalendarScope   string
	ProceduralRules domaincalendar.ProceduralRules
}

func DefaultEventComputationConfig() EventComputationConfig {
	return EventComputationConfig{
		CalendarScope: domaincalendar.ScopeMadrid,
		ProceduralRules: domaincalendar.ProceduralRules{
			AugustNonBusiness: true,
		},
	}
}
