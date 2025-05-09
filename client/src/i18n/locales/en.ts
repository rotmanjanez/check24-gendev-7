export default {
    meta: {
        languageCode: "EN",
        languageName: "English",
    },
    app: {
        title: "GenDev",
        subtitle: "CHECK Out Your Next Internet Provider",
    },
    settings: {
        toggleTheme: "Toggle Theme",
        themes: {
            light: "Light",
            dark: "Dark",
            system: "System"
        }
    },
    product: {
        month: "month",
        months: "months",
        years: "years",
        over: "over",
        contractDurationMonth: "months contract",
        share: "Share this offer",
        copyToClipboard: "Share link copied to clipboard",
        offeredOn: "Offer from",
        showMore: "Show more",
        showLess: "Show less",
        backToSearch: "Search alternatives",
        discount: "Discount",
        bonus: "Bonus",
        unthrottledCapacity: "Unthrottled Speed",
        minAge: "Minimum age",
        maxAge: "Maximum age",
        minContractDuration: "Minimum contract",
        subsequentCosts: "Price after",
        minOrderValue: "Minimum order value",
        installation: {
            title: "Installation",
            included: "included",
            notIncluded: "not included",
        },
        tv: "TV Package",
        noResults: "No offers found",
        shareFailed: "Could not share",
        tryCopyManually: "Try copying the link manually.",
        copyFailed: "Failed to copy link",
        clipboardDenied: "Clipboard access denied.",
        shareNotSupported: "Share not supported",
        useCopyLinkInstead: "Use copy link instead.",
    },
    addressForm: {
        countries: {
            DE: { name: "Germany" },
            AT: { name: "Austria" },
            CH: { name: "Switzerland" },
        },
        edit: "Edit Adress",
        address: {
            title: "Address",
            placeholder: "Enter your address",
            description: "Please enter your full address to find the best internet options available.",
        },
        street: {
            title: "Street",
            placeholder: "Street Name",
            description: "The name of the street where you want to use the internet",
            errors: {
                too_small: "Street name must be at least {minimum} characters",
                too_big: "Street name cannot exceed {maximum} characters",
                invalid_type: "Street name contains invalid characters",
            }
        },
        houseNumber: {
            title: "House Number",
            placeholder: "House Number",
            description: "The number of your house or apartment",
            errors: {
                too_small: "House number must be at least {minimum} characters",
                too_big: "House number cannot exceed {maximum} characters",
                invalid_type: "House number contains invalid characters",
            }
        },
        postalCode: {
            title: "Postal Code",
            placeholder: "Postal Code",
            description: "The postal code of your area",
            errors: {
                too_small: "Postal code must be at least {minimum} characters",
                too_big: "Postal code cannot exceed {maximum} characters",
                invalid_type: "Postal code contains invalid characters",
            }
        },
        region: {
            title: "Region",
            placeholder: "Region",
            description: "The citiy or region where you live",
            errors: {
                too_small: "City name must be at least {minimum} characters",
                too_big: "City name cannot exceed {maximum} characters",
                invalid_type: "City name contains invalid characters",
            }
        },
        country: {
            title: "Country",
            placeholder: "Country",
            description: "The country where you live",
        },
        countryCode: {
            title: "Country Code",
            placeholder: "Country Code",
            description: "The country code of your area",
            errors: {
                invalid_enum_value: "Country code is invalid",
                invalid_type: "Country code is invalid",
            }
        },
        submit: "Compare Now",
    },
    errors: {
        unknown: "An unknown error occurred.",
        network: "Network error occurred.",
    },
    notifications: {
        newOffers: "New offers available",
    },
    actions: {
        refresh: "Refresh",
        retry: "Retry",
        compare: "Compare Now",
    },
} as const