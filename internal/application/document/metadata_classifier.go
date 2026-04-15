package documentapp

import "strings"

type metadataClassification struct {
	DocumentType string
	LegalArea    string
}

func classifyDocumentMetadata(content string) metadataClassification {
	text := normalizeMetadataText(content)

	if !hasEnoughLegalSignal(text) {
		return metadataClassification{
			DocumentType: "unknown",
			LegalArea:    "unknown",
		}
	}

	return metadataClassification{
		DocumentType: classifyDocumentType(text),
		LegalArea:    classifyLegalArea(text),
	}
}

func normalizeMetadataText(content string) string {
	return strings.ToLower(strings.Join(strings.Fields(content), " "))
}

func hasEnoughLegalSignal(text string) bool {
	strongSignals := []string{
		"demanda",
		"sentencia",
		"auto",
		"decreto",
		"providencia",
		"diligencia de ordenación",
		"diligencia de ordenacion",
		"juzgado",
		"tribunal",
		"fundamentos de derecho",
		"suplico al juzgado",
		"despido",
		"finiquito",
		"nómina",
		"nomina",
		"relación laboral",
		"relacion laboral",
		"custodia",
		"divorcio",
		"pensión alimenticia",
		"pension alimenticia",
		"tarjeta de residencia",
		"autorización de residencia",
		"autorizacion de residencia",
		"nie",
		"arraigo",
		"reagrupación familiar",
		"reagrupacion familiar",
		"sociedad",
		"administrador",
		"mercantil",
		"arrendamiento",
		"compraventa",
		"incumplimiento",
	}

	score := countKeywordMatches(text, strongSignals)
	return score >= 2
}

func classifyDocumentType(text string) string {
	switch {
	case containsAny(text,
		"carta de despido",
		"despido disciplinario",
		"despido objetivo",
	):
		return "dismissal_letter"

	case containsAny(text,
		"demanda de divorcio",
		"divorcio contencioso",
		"medidas paternofiliales",
		"guarda y custodia",
	):
		return "divorce_petition"

	case containsAny(text,
		"diligencia de ordenación",
		"diligencia de ordenacion",
		"providencia",
		"decreto",
	):
		return "order"

	case containsAny(text,
		"debo condenar y condeno",
		"fallo",
		"sentencia",
	):
		return "judgment"

	case containsAny(text,
		"autorización de residencia",
		"autorizacion de residencia",
		"tarjeta de residencia",
	) && containsAny(text,
		"resolución",
		"resolucion",
		"denegación",
		"denegacion",
	):
		return "residence_decision"
	}

	type docRule struct {
		documentType string
		keywords     []string
		minScore     int
	}

	rules := []docRule{
		{
			documentType: "payroll",
			keywords: []string{
				"nómina", "nomina", "devengos", "líquido a percibir", "liquido a percibir",
				"base de cotización", "base de cotizacion",
			},
			minScore: 2,
		},
		{
			documentType: "settlement",
			keywords: []string{
				"finiquito", "saldo y finiquito", "liquidación", "liquidacion",
			},
			minScore: 2,
		},
		{
			documentType: "claim",
			keywords: []string{
				"demanda", "hechos", "fundamentos de derecho", "suplico al juzgado",
			},
			minScore: 2,
		},
		{
			documentType: "contract",
			keywords: []string{
				"contrato", "cláusula", "clausula", "las partes acuerdan",
			},
			minScore: 2,
		},
	}

	bestType := "unknown"
	bestScore := 0

	for _, rule := range rules {
		score := countKeywordMatches(text, rule.keywords)
		if score >= rule.minScore && score > bestScore {
			bestScore = score
			bestType = rule.documentType
		}
	}

	return bestType
}

func classifyLegalArea(text string) string {
	type areaRule struct {
		area     string
		keywords []string
		minScore int
	}

	rules := []areaRule{
		{
			area: "labor",
			keywords: []string{
				"despido",
				"nómina", "nomina",
				"finiquito",
				"salario",
				"trabajador",
				"empresa",
				"relación laboral", "relacion laboral",
				"horas extraordinarias",
				"contrato de trabajo",
				"extinción de la relación laboral",
				"extincion de la relacion laboral",
			},
			minScore: 2,
		},
		{
			area: "immigration",
			keywords: []string{
				"nie",
				"residencia",
				"autorización de residencia", "autorizacion de residencia",
				"tarjeta de residencia",
				"reagrupación familiar", "reagrupacion familiar",
				"extranjero",
				"ciudadano de la unión", "ciudadano de la union",
				"arraigo",
			},
			minScore: 2,
		},
		{
			area: "family",
			keywords: []string{
				"divorcio",
				"custodia",
				"pensión alimenticia", "pension alimenticia",
				"medidas paterno filiales",
				"medidas paternofiliales",
				"guarda y custodia",
				"hijo menor",
				"régimen de visitas", "regimen de visitas",
			},
			minScore: 2,
		},
		{
			area: "commercial",
			keywords: []string{
				"sociedad",
				"administrador",
				"mercantil",
				"acciones",
				"participaciones",
				"junta general",
				"responsabilidad limitada",
				"s.l.",
				"sl",
			},
			minScore: 2,
		},
		{
			area: "civil",
			keywords: []string{
				"contrato",
				"reclamación de cantidad", "reclamacion de cantidad",
				"incumplimiento",
				"arrendamiento",
				"alquiler",
				"compraventa",
			},
			minScore: 2,
		},
	}

	bestArea := "unknown"
	bestScore := 0

	for _, rule := range rules {
		score := countKeywordMatches(text, rule.keywords)
		if score >= rule.minScore && score > bestScore {
			bestScore = score
			bestArea = rule.area
		}
	}

	return bestArea
}

func containsAny(text string, values ...string) bool {
	for _, value := range values {
		if strings.Contains(text, value) {
			return true
		}
	}
	return false
}

func countKeywordMatches(text string, keywords []string) int {
	score := 0
	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			score++
		}
	}
	return score
}
