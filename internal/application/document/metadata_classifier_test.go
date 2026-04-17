package documentapp

import "testing"

func TestClassifyDocumentMetadata_DismissalLetter(t *testing.T) {
	content := `
	CARTA DE DESPIDO

	La empresa comunica al trabajador la extinción de la relación laboral
	por causas disciplinarias. Se adjunta finiquito y liquidación.
	`

	got := classifyDocumentMetadata(content)

	if got.DocumentType != "dismissal_letter" {
		t.Fatalf("expected document type %q, got %q", "dismissal_letter", got.DocumentType)
	}

	if got.LegalArea != "labor" {
		t.Fatalf("expected legal area %q, got %q", "labor", got.LegalArea)
	}
}

func TestClassifyDocumentMetadata_DivorcePetition(t *testing.T) {
	content := `
	DEMANDA DE DIVORCIO CONTENCIOSO

	Se interesan medidas paternofiliales, guarda y custodia,
	régimen de visitas y pensión alimenticia respecto del hijo menor.
	`

	got := classifyDocumentMetadata(content)

	if got.DocumentType != "divorce_petition" {
		t.Fatalf("expected document type %q, got %q", "divorce_petition", got.DocumentType)
	}

	if got.LegalArea != "family" {
		t.Fatalf("expected legal area %q, got %q", "family", got.LegalArea)
	}
}

func TestClassifyDocumentMetadata_ResidenceDecision(t *testing.T) {
	content := `
	RESOLUCIÓN DE DENEGACIÓN

	Se deniega la autorización de residencia temporal y la expedición
	de la tarjeta de residencia. Consta NIE del interesado.
	`

	got := classifyDocumentMetadata(content)

	if got.DocumentType != "residence_decision" {
		t.Fatalf("expected document type %q, got %q", "residence_decision", got.DocumentType)
	}

	if got.LegalArea != "immigration" {
		t.Fatalf("expected legal area %q, got %q", "immigration", got.LegalArea)
	}
}

func TestClassifyDocumentMetadata_ProceduralOrder_BriefOrder(t *testing.T) {
	content := `
	DILIGENCIA DE ORDENACIÓN

	Notifíquese la resolución. Se concede plazo de cinco días
	para formular alegaciones.
	`

	got := classifyDocumentMetadata(content)

	if got.DocumentType != "order" {
		t.Fatalf("expected document type %q, got %q", "order", got.DocumentType)
	}

	if got.LegalArea != "procedural" {
		t.Fatalf("expected legal area %q, got %q", "procedural", got.LegalArea)
	}
}

func TestClassifyDocumentMetadata_LaborClaim(t *testing.T) {
	content := `
	DEMANDA

	HECHOS

	El trabajador fue despedido por la empresa y reclama salario,
	finiquito y horas extraordinarias.

	FUNDAMENTOS DE DERECHO

	SUPLICO AL JUZGADO que dicte sentencia estimatoria.
	`

	got := classifyDocumentMetadata(content)

	if got.DocumentType != "claim" {
		t.Fatalf("expected document type %q, got %q", "claim", got.DocumentType)
	}

	if got.LegalArea != "labor" {
		t.Fatalf("expected legal area %q, got %q", "labor", got.LegalArea)
	}
}

func TestClassifyDocumentMetadata_Payroll(t *testing.T) {
	content := `
	NÓMINA

	Devengos
	Base de cotización
	Líquido a percibir
	IRPF
	`

	got := classifyDocumentMetadata(content)

	if got.DocumentType != "payroll" {
		t.Fatalf("expected document type %q, got %q", "payroll", got.DocumentType)
	}

	if got.LegalArea != "labor" {
		t.Fatalf("expected legal area %q, got %q", "labor", got.LegalArea)
	}
}

func TestClassifyDocumentMetadata_Settlement(t *testing.T) {
	content := `
	SALDO Y FINIQUITO

	Liquidación final e indemnización derivadas de la extinción
	de la relación laboral.
	`

	got := classifyDocumentMetadata(content)

	if got.DocumentType != "settlement" {
		t.Fatalf("expected document type %q, got %q", "settlement", got.DocumentType)
	}

	if got.LegalArea != "labor" {
		t.Fatalf("expected legal area %q, got %q", "labor", got.LegalArea)
	}
}

func TestClassifyDocumentMetadata_NonLegalDocument(t *testing.T) {
	content := `
	CURRICULUM VITAE

	Experiencia en desarrollo web, Linux, soporte técnico y atención al cliente.
	Conocimientos de Go, JavaScript, SAP y administración de sistemas.
	`

	got := classifyDocumentMetadata(content)

	if got.DocumentType != "unknown" {
		t.Fatalf("expected document type %q, got %q", "unknown", got.DocumentType)
	}

	if got.LegalArea != "unknown" {
		t.Fatalf("expected legal area %q, got %q", "unknown", got.LegalArea)
	}
}
