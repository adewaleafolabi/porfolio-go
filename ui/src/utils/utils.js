import {DateTime} from "luxon";

export const formatCurrency =(amount, currency='CAD') =>{
    if (!amount) {
        return "--";
    }

    const formatter = new Intl.NumberFormat("en-CA", {
        style: "currency",
        currency: currency
    });
    return formatter.format(amount );
}

export const formatDecimal = (amount, dp=2) => {
    if (!amount) {
        return "--";
    }

    const formatter = new Intl.NumberFormat("en-CA", {
        maximumFractionDigits: dp
    });
    return formatter.format(amount);
}

export const formatPercentage = (amount, dp=2) => {
    if (!amount) {
        return "--";
    }

    const formatter = new Intl.NumberFormat("en-CA", {
        style: 'percent',
        maximumFractionDigits: dp
    });
    return formatter.format(amount);
}

export const formatDate = (input) =>{
    if (!input) {
        return "--";
    }
    if (String(input).startsWith("0")) {
        return "--";
    }

    return DateTime.fromISO(input).setLocale('en-CA').toISODate();
}