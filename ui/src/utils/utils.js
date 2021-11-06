import {DateTime} from "luxon";

export const formatCurrency =(amount, currency='CAD', compact=false) =>{
    if (!amount) {
        return "--";
    }

    const option = {
        style: "currency",
        currency: currency
    }
    if (compact){
        option.notation = 'compact'
    }
    const formatter = new Intl.NumberFormat("en-CA", option);

    return formatter.format(amount );
}
export const formatNumber =(amount) =>{
    if (!amount) {
        return "--";
    }

    const option = {
        maximumFractionDigits: 2,
        notation: 'compact'
    }

    const formatter = new Intl.NumberFormat("en-CA", option);

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

export  const hashCode = (input)=> {
    let hash = 0, i, chr;
    if (input.length === 0) return hash;
    for (i = 0; i < input.length; i++) {
        chr   = input.charCodeAt(i);
        hash  = ((hash << 5) - hash) + chr;
        hash |= 0; // Convert to 32bit integer
    }
    return hash;
};
