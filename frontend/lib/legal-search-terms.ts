export type LegalSearchTerm = {
  trigger: string;
  value: string;
  aliases?: string[];
  keywords?: string[];
};

export const LEGAL_SEARCH_TERMS: LegalSearchTerm[] = [
  // =========================
  // Procesal general
  // =========================

  {
    trigger: "sent",
    value: "sentencia",
    aliases: ["sen", "fallo"],
    keywords: [
      "resolucion judicial",
      "juez",
      "tribunal",
    ],
  },

  {
    trigger: "art",
    value: "artículo",
    aliases: ["articulo"],
    keywords: [
      "ley",
      "norma",
      "precepto",
    ],
  },

  {
    trigger: "dem",
    value: "demanda",
    aliases: [
      "demand",
      "deman",
    ],
    keywords: [
      "escrito rector",
      "procedimiento",
      "actor",
    ],
  },

  {
    trigger: "cont",
    value: "contestación",
    aliases: ["contestacion"],
    keywords: [
      "oposicion",
      "demandado",
      "alegaciones",
    ],
  },

  {
    trigger: "recu",
    value: "recurso",
    aliases: ["rec"],
    keywords: ["impugnacion"],
  },

  {
    trigger: "apel",
    value: "apelación",
    aliases: ["apelacion"],
    keywords: [
      "segunda instancia",
      "audiencia provincial",
    ],
  },

  {
    trigger: "suplic",
    value: "suplicación",
    aliases: ["suplicacion"],
    keywords: [
      "social",
      "tsj",
      "recurso social",
    ],
  },

  {
    trigger: "cas",
    value: "casación",
    aliases: ["casacion"],
    keywords: [
      "tribunal supremo",
      "recurso extraordinario",
    ],
  },

  {
    trigger: "ejec",
    value: "ejecución",
    aliases: [
      "ejecu",
      "ejecucion",
    ],
    keywords: [
      "titulo judicial",
      "embargo",
      "apremio",
    ],
  },

  {
    trigger: "notif",
    value: "notificación",
    aliases: ["notificacion"],
    keywords: [
      "lexnet",
      "comunicacion",
    ],
  },

  {
    trigger: "tras",
    value: "traslado",
    keywords: [
      "plazo",
      "alegaciones",
    ],
  },

  {
    trigger: "subsa",
    value: "subsanación",
    aliases: ["subsanacion"],
    keywords: [
      "requerimiento",
      "defecto",
      "plazo",
    ],
  },

  {
    trigger: "requ",
    value: "requerimiento",
    aliases: [
      "req",
      "requer",
    ],
    keywords: [
      "subsanacion",
      "aportacion",
      "documentacion",
    ],
  },

  {
    trigger: "decr",
    value: "decreto",
    keywords: [
      "laj",
      "letrado administracion justicia",
    ],
  },

  {
    trigger: "auto",
    value: "auto",
    keywords: [
      "resolucion judicial",
      "juez",
    ],
  },

  {
    trigger: "provid",
    value: "providencia",
    aliases: ["provi"],
    keywords: [
      "resolucion judicial",
      "tramite",
    ],
  },

  {
    trigger: "dilig",
    value: "diligencia",
    aliases: [
      "diligencia ordenacion",
      "ordenacion",
    ],
    keywords: [
      "laj",
      "tramite",
    ],
  },

  {
    trigger: "vista",
    value: "vista",
    keywords: [
      "juicio",
      "señalamiento",
      "comparecencia",
    ],
  },

  {
    trigger: "juic",
    value: "juicio",
    aliases: ["juicio verbal"],
    keywords: [
      "vista",
      "señalamiento",
      "acto",
    ],
  },

  {
    trigger: "comp",
    value: "comparecencia",
    keywords: [
      "vista",
      "acto judicial",
      "señalamiento",
    ],
  },

  {
    trigger: "exp",
    value: "expediente",
    keywords: [
      "procedimiento",
      "autos",
      "referencia",
    ],
  },

  {
    trigger: "proc",
    value: "procedimiento",
    keywords: [
      "autos",
      "expediente",
      "juicio",
    ],
  },

  // =========================
  // Civil / Familia
  // =========================

  {
    trigger: "monit",
    value: "monitorio",
    keywords: [
      "deuda",
      "reclamacion cantidad",
      "comunidad propietarios",
    ],
  },

  {
    trigger: "desa",
    value: "desahucio",
    keywords: [
      "arrendamiento",
      "precario",
      "lanzamiento",
    ],
  },

  {
    trigger: "lanz",
    value: "lanzamiento",
    keywords: [
      "desahucio",
      "posesion",
      "vivienda",
    ],
  },

  {
    trigger: "prec",
    value: "precario",
    keywords: [
      "ocupacion",
      "desahucio",
      "posesion",
    ],
  },

  {
    trigger: "div",
    value: "divorcio",
    keywords: [
      "familia",
      "matrimonio",
      "conyuge",
    ],
  },

  {
    trigger: "cust",
    value: "custodia",
    keywords: [
      "menor",
      "familia",
      "guarda",
    ],
  },

  {
    trigger: "alim",
    value: "alimentos",
    keywords: [
      "pension",
      "menor",
      "hijos",
    ],
  },

  {
    trigger: "vis",
    value: "visitas",
    keywords: [
      "menor",
      "regimen",
      "familia",
    ],
  },

  {
    trigger: "ganan",
    value: "gananciales",
    keywords: [
      "liquidacion",
      "sociedad",
      "matrimonio",
    ],
  },

  {
    trigger: "usuca",
    value: "usucapión",
    aliases: ["usucapion"],
    keywords: [
      "dominio",
      "propiedad",
      "posesion",
    ],
  },

  // =========================
  // Laboral / Seguridad Social
  // =========================

  {
    trigger: "desp",
    value: "despido",
    keywords: [
      "laboral",
      "improcedente",
      "disciplinario",
    ],
  },

  {
    trigger: "imv",
    value: "ingreso mínimo vital",
    aliases: ["minimo vital"],
    keywords: [
      "prestacion",
      "seguridad social",
    ],
  },

  {
    trigger: "inca",
    value: "incapacidad permanente",
    aliases: [
      "perm",
      "permanente",
    ],
    keywords: [
      "prestacion",
      "enfermedad",
    ],
  },

  {
    trigger: "smi",
    value: "salario mínimo interprofesional",
    aliases: [
      "salario minimo",
      "minimo interprofesional",
    ],
    keywords: [
      "laboral",
      "salario",
      "retribucion",
    ],
  },

  {
    trigger: "it",
    value: "incapacidad temporal",
    aliases: ["baja"],
    keywords: ["seguridad social"],
  },

  {
    trigger: "ss",
    value: "seguridad social",
    aliases: [
      "inss",
      "tgss",
      "seg social",
    ],
    keywords: [
      "prestacion",
      "cotizacion",
      "laboral",
    ],
  },

  {
    trigger: "fog",
    value: "fogasa",
    aliases: ["fondo garantia salarial"],
    keywords: [
      "salarios",
      "insolvencia",
      "laboral",
    ],
  },

  {
    trigger: "finiq",
    value: "finiquito",
    keywords: [
      "liquidacion",
      "vacaciones",
      "pagas",
    ],
  },

  {
    trigger: "extra",
    value: "horas extraordinarias",
    aliases: ["horas extra"],
    keywords: [
      "jornada",
      "laboral",
      "salario",
    ],
  },

  // =========================
  // Extranjería
  // =========================

  {
    trigger: "arra",
    value: "arraigo",
    keywords: [
      "extranjeria",
      "residencia",
      "social",
      "formacion",
    ],
  },

  {
    trigger: "ext",
    value: "extranjería",
    aliases: ["extranjeria"],
    keywords: [
      "nie",
      "tie",
      "residencia",
      "autorizacion",
    ],
  },

  {
    trigger: "nie",
    value: "NIE",
    keywords: [
      "extranjero",
      "identificacion",
    ],
  },

  {
    trigger: "tie",
    value: "TIE",
    keywords: [
      "tarjeta",
      "residencia",
      "extranjero",
    ],
  },

  {
    trigger: "ue",
    value: "ciudadano de la Unión Europea",
    aliases: [
      "union europea",
      "comunitario",
    ],
    keywords: [
      "familiar",
      "residencia",
    ],
  },

  {
    trigger: "fam",
    value: "familiar comunitario",
    aliases: ["comunitario"],
    keywords: [
      "ue",
      "ciudadano union",
      "residencia",
    ],
  },

  {
    trigger: "asilo",
    value: "protección internacional",
    aliases: [
      "proteccion internacional",
      "refugio",
    ],
    keywords: ["extranjeria"],
  },

  // =========================
  // Administrativo
  // =========================

  {
    trigger: "contad",
    value: "contencioso-administrativo",
    aliases: ["contencioso"],
    keywords: [
      "administracion",
      "jurisdiccion",
    ],
  },

  {
    trigger: "alza",
    value: "recurso de alzada",
    keywords: [
      "administrativo",
      "impugnacion",
      "resolucion",
    ],
  },

  {
    trigger: "repo",
    value: "reposición",
    aliases: ["reposicion"],
    keywords: [
      "administrativo",
      "recurso",
      "impugnacion",
    ],
  },

  {
    trigger: "sil",
    value: "silencio administrativo",
    keywords: [
      "administracion",
      "plazo",
      "resolucion",
    ],
  },
];