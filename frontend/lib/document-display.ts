// frontend/lib/document-display.ts

export function displayMimeType(value: string): string {
  switch (value) {
    case "text/plain":
      return "Texto plano";
    case "application/pdf":
      return "PDF";
    case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
      return "Word";
    default:
      return value || "Sin tipo";
  }
}

export function displayDocumentType(value?: string): string {
  switch (value) {
    case "dismissal_letter":
      return "Carta de despido";
    case "divorce_petition":
      return "Demanda de divorcio";
    case "residence_decision":
      return "Resolución de extranjería";
    case "order":
      return "Resolución / acto procesal";
    case "appeal_motion":
      return "Recurso";
    case "appeal_brief":
      return "Escrito de recurso";
    case "conciliation_filing":
      return "Papeleta de conciliación";
    case "order_decision":
      return "Resolución judicial";
    case "claim":
      return "Demanda / reclamación";
    case "judgment":
      return "Sentencia";
    case "payroll":
      return "Nómina";
    case "settlement":
      return "Finiquito";
    case "answer":
      return "Contestación";
    case "contract":
      return "Contrato";
    case "administrative_resolution":
      return "Resolución administrativa";
    case "enforcement_filing":
      return "Demanda de ejecución";
    case "monitorio_filing":
      return "Petición inicial de monitorio";
    case "notice":
      return "Notificación";
    case "evidence":
      return "Documento probatorio";
    case "unknown":
    case "":
    case undefined:
      return "Sin clasificar";
    default:
      return value;
  }
}

export function displayLegalArea(value?: string): string {
  switch (value) {
    case "procedural":
      return "Procesal";
    case "civil":
      return "Civil";
    case "family":
      return "Familia";
    case "labor":
    case "laboral":
      return "Laboral";
    case "administrative":
    case "administrativo":
      return "Administrativo";
    case "immigration":
    case "extranjeria":
      return "Extranjería";
    case "commercial":
    case "mercantil":
      return "Mercantil";
    case "unknown":
    case "":
    case undefined:
      return "Sin clasificar";
    default:
      return value;
  }
}

export function documentExtractionLabel(
  hasExtractedText?: boolean,
): string {
  return hasExtractedText ? "Texto extraído" : "Sin texto extraíble";
}

export function documentEventsLabel(eventsCount?: number): string {
  const count = eventsCount ?? 0;

  if (count === 0) return "0 hitos";
  if (count === 1) return "1 hito";

  return `${count} hitos`;
}

export function displayCaseType(value?: string) {
  switch (value) {
    case "civil":
      return "Civil";

    case "labor":
    case "laboral":
      return "Laboral";

    case "extranjeria":
      return "Extranjería";

    case "mercantil":
      return "Mercantil";

    case "administrativo":
    case "administrative":
      return "Administrativo";

    case "otros":
      return "Otros";

    case "":
    case undefined:
      return "Sin clasificar";

    default:
      return value;
  }
}

export function displayStatus(value?: string) {
  switch (value) {
    case "open":
      return "Abierto";

    case "closed":
      return "Cerrado";

    case "archived":
      return "Archivado";

    case "":
    case undefined:
      return "Sin estado";

    default:
      return value;
  }
}