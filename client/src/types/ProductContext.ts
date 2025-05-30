import { InternetProduct, Address } from '@/api'

export interface ProductContext {
    address: Address
    product: InternetProduct
    version: string
}

interface AverageMonthlyCost {
    monthlyCostInCent: number;
    durationInMonths: number;
}

export function getAverageMonthlyCost(
    product: InternetProduct | undefined
): AverageMonthlyCost {
    if (!product) {
        return { monthlyCostInCent: 0, durationInMonths: 0 };
    }
    const pricing = product.pricing;
    if (!pricing) {
        console.warn('No pricing information available.');
        return { monthlyCostInCent: 0, durationInMonths: 0 };
    }

    const {
        monthlyCostInCent,
        contractDurationInMonths,
        minContractDurationInMonths,
        subsequentCosts,
        absoluteDiscount,
        percentageDiscount,
    } = pricing;

    if (monthlyCostInCent <= 0) {
        console.warn('No valid monthly cost provided.');
        return { monthlyCostInCent: 0, durationInMonths: 0 };
    }

    let duration = Math.max(minContractDurationInMonths ?? 0, 24);
    if (contractDurationInMonths) {
        duration = Math.min(duration, contractDurationInMonths);
    }
    const startMonth = subsequentCosts?.startMonth ?? duration;

    // Calculate initial total cost
    const initialCost =
        monthlyCostInCent * startMonth +
        (subsequentCosts
            ? subsequentCosts.monthlyCostInCent * (duration - startMonth)
            : 0);

    let totalCost = initialCost;
    // Apply percentage discount first (with max cap if present)
    if (percentageDiscount && percentageDiscount.percentage) {
        let pctAmount = (totalCost * percentageDiscount.percentage) / 100;
        if (percentageDiscount.maxDiscountInCent !== undefined && percentageDiscount.maxDiscountInCent !== null) {
            pctAmount = Math.min(pctAmount, percentageDiscount.maxDiscountInCent);
        }
        totalCost -= pctAmount;
    }
    // Apply absolute discount
    if (absoluteDiscount && absoluteDiscount.valueInCent) {
        totalCost -= absoluteDiscount.valueInCent;
    }

    return {
        monthlyCostInCent: Math.round(totalCost / duration),
        durationInMonths: duration
    };
}

