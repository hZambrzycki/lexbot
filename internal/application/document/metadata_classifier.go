package documentapp

import "strings"

type metadataClassification struct {
	DocumentType string
	LegalArea    string
}

func classifyDocumentMetadata(content string) metadataClassification {
	text := normalizeMetadataText(content)

	if isLikelyCVOrProfile(text) {
		return metadataClassification{
			DocumentType: "non_legal",
			LegalArea:    "non_legal",
		}
	}

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
	return normalizeASCIIText(content)
}

func isLikelyCVOrProfile(text string) bool {
	cvSignals := []string{
		"curriculum",
		"curriculum vitae",
		"cv",
		"perfil profesional",
		"sobre mi",
		"experiencia laboral",
		"educacion",
		"certificaciones",
		"stack tecnico",
		"lenguajes",
		"backend",
		"frontend",
		"github",
		"linkedin",
		"desarrollador",
		"programador",
		"java",
		"python",
		"javascript",
		"react",
		"next.js",
		"docker",
		"kubernetes",
		"linux",
		"sql",
	}

	legalSignals := []string{
		"juzgado",
		"tribunal",
		"procedimiento",
		"autos",
		"demanda",
		"contestacion a la demanda",
		"fundamentos de derecho",
		"suplico",
		"fallo",
		"sentencia",
		"auto",
		"providencia",
		"diligencia de ordenacion",
		"decreto",
		"notifiquese",
		"recurso",
		"plazo",
		"dias habiles",
		"dias naturales",
	}

	cvScore := countKeywordMatches(text, cvSignals)
	legalScore := countKeywordMatches(text, legalSignals)

	return cvScore >= 3 && legalScore == 0
}

func hasEnoughLegalSignal(text string) bool {
	if containsAny(text,
		"diligencia de ordenacion",
		"providencia",
		"decreto",
		"auto",
		"sentencia",
		"nomina",
		"carta de despido",
		"saldo y finiquito",
		"demanda de divorcio",
		"medidas paternofiliales",
		"autorizacion de residencia",
		"tarjeta de residencia",
		"recurso de reposicion",
		"recurso de apelacion",
		"notifiquese",
		"se concede plazo",
		"dias habiles",
		"dias naturales",
		"formular alegaciones",
		"aportar la documentacion",
		"presentar escrito",
	) {
		return true
	}

	strongSignals := []string{
		"demanda",
		"contestacion a la demanda",
		"juzgado",
		"tribunal",
		"fundamentos de derecho",
		"suplico al juzgado",
		"hechos",
		"fallo",
		"despido",
		"finiquito",
		"nomina",
		"relacion laboral",
		"custodia",
		"divorcio",
		"pension alimenticia",
		"tarjeta de residencia",
		"autorizacion de residencia",
		"nie",
		"arraigo",
		"reagrupacion familiar",
		"sociedad",
		"administrador",
		"mercantil",
		"arrendamiento",
		"compraventa",
		"incumplimiento",
		"recurso de reposicion",
		"recurso de apelacion",
		"demanda ejecutiva",
		"ejecucion",
		"monitorio",
		"smac",
		"devengos",
		"liquido a percibir",
		"base de cotizacion",
		"irpf",
		"notifiquese",
		"plazo",
		"alegaciones",
		"dias habiles",
		"dias naturales",
		"aportar la documentacion",
		"presentar escrito",
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
		"medidas paterno filiales",
		"guarda y custodia",
	):
		return "divorce_petition"

	case containsAny(text,
		"diligencia de ordenacion",
		"providencia",
		"decreto",
	):
		return "order"

	case containsAny(text,
		"recurso de reposicion",
	):
		return "appeal_motion"

	case containsAny(text,
		"recurso de apelacion",
	):
		return "appeal_brief"

	case containsAny(text,
		"papeleta de conciliacion",
		"smac",
		"acto de conciliacion",
	):
		return "conciliation_filing"

	case containsAny(text,
		"autorizacion de residencia",
		"tarjeta de residencia",
		"nie",
	) && containsAny(text,
		"resolucion",
		"denegacion",
		"concede",
		"se deniega",
	):
		return "residence_decision"

	case containsAny(text,
		"auto",
	) && containsAny(text,
		"parte dispositiva",
		"fundamentos de derecho",
		"razonamientos juridicos",
	):
		return "order_decision"

	case containsAny(text,
		"notifiquese",
		"se concede plazo",
		"plazo",
		"dias habiles",
		"dias naturales",
		"formular alegaciones",
		"aportar la documentacion",
		"presentar escrito",
	):
		return "order"
	}

	type docRule struct {
		documentType string
		keywords     []string
		minScore     int
	}

	rules := []docRule{
		{
			documentType: "claim",
			keywords: []string{
				"demanda",
				"hechos",
				"fundamentos de derecho",
				"suplico al juzgado",
			},
			minScore: 3,
		},
		{
			documentType: "judgment",
			keywords: []string{
				"sentencia",
				"sentencia firme",
				"debo condenar y condeno",
				"fallo",
				"parte dispositiva",
				"fundamentos de derecho",
			},
			minScore: 2,
		},
		{
			documentType: "payroll",
			keywords: []string{
				"nomina",
				"devengos",
				"liquido a percibir",
				"base de cotizacion",
				"contingencias comunes",
				"irpf",
			},
			minScore: 2,
		},
		{
			documentType: "settlement",
			keywords: []string{
				"finiquito",
				"saldo y finiquito",
				"liquidacion final",
				"indemnizacion",
			},
			minScore: 2,
		},
		{
			documentType: "answer",
			keywords: []string{
				"contestacion a la demanda",
				"esta parte se opone",
				"se impugna",
				"hechos",
				"fundamentos de derecho",
			},
			minScore: 2,
		},
		{
			documentType: "contract",
			keywords: []string{
				"contrato",
				"clausula",
				"las partes acuerdan",
				"objeto del contrato",
				"duracion",
			},
			minScore: 2,
		},
		{
			documentType: "administrative_resolution",
			keywords: []string{
				"resolucion",
				"antecedentes de hecho",
				"fundamentos de derecho",
				"parte dispositiva",
			},
			minScore: 3,
		},
		{
			documentType: "enforcement_filing",
			keywords: []string{
				"demanda ejecutiva",
				"ejecucion",
				"despachar ejecucion",
				"titulo ejecutivo",
			},
			minScore: 2,
		},
		{
			documentType: "monitorio_filing",
			keywords: []string{
				"procedimiento monitorio",
				"peticion inicial de procedimiento monitorio",
				"reclamacion de cantidad",
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
				"nomina",
				"finiquito",
				"salario",
				"trabajador",
				"empresa",
				"relacion laboral",
				"horas extraordinarias",
				"contrato de trabajo",
				"extincion de la relacion laboral",
				"smac",
				"conciliacion laboral",
				"devengos",
				"liquido a percibir",
				"base de cotizacion",
				"irpf",
				"contingencias comunes",
			},
			minScore: 2,
		},
		{
			area: "immigration",
			keywords: []string{
				"nie",
				"residencia",
				"autorizacion de residencia",
				"tarjeta de residencia",
				"reagrupacion familiar",
				"extranjero",
				"ciudadano de la union",
				"arraigo",
				"oficina de extranjeria",
			},
			minScore: 2,
		},
		{
			area: "family",
			keywords: []string{
				"divorcio",
				"custodia",
				"pension alimenticia",
				"medidas paterno filiales",
				"medidas paternofiliales",
				"guarda y custodia",
				"hijo menor",
				"regimen de visitas",
				"patria potestad",
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
				"sociedad limitada",
			},
			minScore: 2,
		},
		{
			area: "civil",
			keywords: []string{
				"contrato",
				"reclamacion de cantidad",
				"incumplimiento",
				"arrendamiento",
				"alquiler",
				"compraventa",
				"resolucion contractual",
			},
			minScore: 2,
		},
		{
			area: "procedural",
			keywords: []string{
				"juzgado",
				"tribunal",
				"providencia",
				"diligencia de ordenacion",
				"decreto",
				"auto",
				"sentencia",
				"recurso de reposicion",
				"recurso de apelacion",
				"ejecucion",
				"notifiquese",
				"plazo",
				"alegaciones",
				"dias habiles",
				"dias naturales",
				"aportar la documentacion",
				"presentar escrito",
				"se concede plazo",
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
	return containsAnyPhrase(text, values...)
}

func containsAnyPhrase(text string, values ...string) bool {
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
		keyword = strings.TrimSpace(keyword)
		if keyword == "" {
			continue
		}

		if len(keyword) <= 3 {
			if hasStandaloneWord(text, keyword) {
				score++
			}
			continue
		}

		if strings.Contains(text, keyword) {
			score++
		}
	}

	return score
}

func hasStandaloneWord(text string, word string) bool {
	fields := strings.FieldsFunc(text, func(r rune) bool {
		return !(r >= 'a' && r <= 'z') && !(r >= '0' && r <= '9')
	})

	for _, field := range fields {
		if field == word {
			return true
		}
	}

	return false
}
