export default {
    meta: {
        languageCode: "DE",
        languageName: "Deutsch",
    },
    app: {
        title: "GenDev",
        subtitle: "CHECK Dir Deinen Neuen Internet-Provider",
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
        month: "Monat",
        months: "Monate",
        years: "Jahre",
        over: "über",
        contractDurationMonth: "Monate Vertragslaufzeit",
        share: "Angebote Teilen",
        copyToClipboard: "Link in die Zwischenablage kopiert",
        offeredOn: "Angebot vom",
        showMore: "Mehr anzeigen",
        showLess: "Weniger anzeigen",
        backToSearch: "Alle Angebote",
        discount: "Rabatt",
        bonus: "Bonus",
        unthrottledCapacity: "Ungedrosselte Geschwindigkeit",
        minAge: "Mindestalter",
        maxAge: "Höchstalter",
        minContractDuration: "Mindestvertragslaufzeit",
        subsequentCosts: "Preis nach",
        minOrderValue: "Mindestbestellwert",
        installation: {
            title: "Installation",
            included: "inklusive",
            notIncluded: "nicht inklusive",
        },
        tv: "TV-Paket",
        noResults: "Keine Angebote gefunden",
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
            description: "Die Postleitzahl Deiner Stadt",
            errors: {
                too_small: "Die Postleitzahl muss mindestens {minimum} Zeichen haben",
                too_big: "Die Postleitzahl darf maximal {maximum} Zeichen haben",
                invalid_type: "Die Postleitzahl enthält ungültige Zeichen",
            }
        },
        city: {
            title: "Stadt/Region",
            placeholder: "Stadt oder Region",
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
    },
      notifications: {
        newOffers: "Neue Angebote verfügbar",
    },
    actions: {
        refresh: "Aktualisieren",
        retry: "Erneut versuchen",
        compare: "Jetzt vergleichen",
    },  
    errors: {
        unknown: "Ein unbekannter Fehler ist aufgetreten.",
        network: "Ein Netzwerkfehler ist aufgetreten :(",
    },
    productFilter: {
        title: "Angebote filtern",
        description: "Angebote sortieren und filtern",
        sortBy: "Sortieren nach",
        selectField: "Feld auswählen",
        price: "Preis",
        speed: "Geschwindigkeit",
        ascending: "Aufsteigend",
        descending: "Absteigend",
        reset: "Zurücksetzen",
        filters: "Filter",
        age: "Alter",
        ageLabel: "Ihr Alter (für Jugendangebote)",
        agePlaceholder: "Alter eingeben",
        providers: "Anbieter",
        allProviders: "Alle Anbieter",
        selectProviders: "Anbieter auswählen",
        priceRange: "Preisbereich",
        speedRange: "Geschwindigkeitsbereich",
        maxPrice: "Max €{max}/Monat",
        minSpeed: "Min {min} Mbps",
        showFilters: "Filter anzeigen",
        hideFilters: "Filter ausblenden",
        clearFilters: "Alle Filter löschen",
        applyFilters: "Filter anwenden",
        includeTV: "TV-Pakete einschließen",
        tvLabel: "TV enthalten",
    },
} as const