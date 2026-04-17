package calendar

type ProceduralRules struct {
	AugustNonBusiness bool
}

func DefaultProceduralRules() ProceduralRules {
	return ProceduralRules{
		AugustNonBusiness: false,
	}
}
