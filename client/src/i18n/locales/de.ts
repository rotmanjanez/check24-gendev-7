export default {
    meta: {
        languageCode: "DE",
        languageName: "Deutsch",
    },
    settings: {
        toggleTheme: "Stil wechseln",
        themes: {
            light: "Hell",
            dark: "Dunkel",
            system: "System"
        }
    },
    product: {
        shareFailed: "Teilen nicht möglich",
        tryCopyManually: "Bitte Link manuell kopieren.",
        copyFailed: "Link konnte nicht kopiert werden",
        clipboardDenied: "Zugriff auf Zwischenablage verweigert.",
        shareNotSupported: "Teilen nicht unterstützt",
        useCopyLinkInstead: "Bitte Link kopieren.",
    },
    addressForm: {
        countries: {
            DE: { name: "Deutschland" },
            AT: { name: "Österreich" },
            CH: { name: "Schweiz" },
        },
        edit: "Adresse ändern",
        address: {
            title: "Adresse",
            placeholder: "Gib Deine Adresse ein",
            description: "Bitte gib Deine vollständige Adresse ein, um die besten Internet-Optionen zu finden.",
        },
        street: {
            title: "Straße",
            placeholder: "Straßenname",
            description: "Der Name der Straße, in der Du Internet nutzen möchtest",
            errors: {
                too_small: "Die Straße muss mindestens {minimum} Zeichen haben",
                too_big: "Die Straße darf maximal {maximum} Zeichen haben",
                invalid_type: "Die Straße enthält ungültige Zeichen",
            }
        },
        houseNumber: {
            title: "Hausnummer",
            placeholder: "Hausnummer",
            description: "Die Nummer des Hauses oder Deiner Wohnung",
            errors: {
                too_small: "Die Hausnummer muss mindestens {minimum} Zeichen haben",
                too_big: "Die Hausnummer darf maximal {maximum} Zeichen haben",
                invalid_type: "Die Hausnummer enthält ungültige Zeichen",
            }
        },
        postalCode: {
            title: "Postleitzahl",
            placeholder: "Postleitzahl",
            description: "Die Postleitzahl Deiner Region",
            errors: {
                too_small: "Die Postleitzahl muss mindestens {minimum} Zeichen haben",
                too_big: "Die Postleitzahl darf maximal {maximum} Zeichen haben",
                invalid_type: "Die Postleitzahl enthält ungültige Zeichen",
            }
        },
        region: {
            title: "Region",
            placeholder: "Region",
            description: "Die Stadt oder Region, in der Du lebst",
            errors: {
                too_small: "Die Stadt muss mindestens {minimum} Zeichen haben",
                too_big: "Die Stadt darf maximal {maximum} Zeichen haben",
                invalid_type: "Die Stadt enthält ungültige Zeichen",
            }
        },
        country: {
            title: "Land",
            placeholder: "Land",
            description: "Das Land, in dem Du lebst",
        },
        countryCode: {
            title: "Ländercode",
            placeholder: "Ländercode",
            description: "Der Ländercode, in dem Du lebst",
            errors: {
                invalid_enum_value: "Der Ländercode ist ungültig",
                invalid_type: "Der Ländercode ist ungültig",
            }
        },
        submit: "Jetzt vergleichen",
    },
    errors: {
        unknown: "Ein unbekannter Fehler ist aufgetreten.",
    },
} as const