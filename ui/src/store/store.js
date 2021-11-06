//store.js
import Vue from 'vue'

export const state = Vue.observable({privacy: false, portfolioID:'', portfolioKeyValue: []})

export const mutations = {
    setPrivacy(newVal) {
        console.log('Setting "privacy": ', newVal)
        state.privacy = newVal
        localStorage.setItem("privateMode", newVal)
    },
    setPortfolioID(newVal) {
        console.log('Setting "portfolioID": ', newVal)
        state.portfolioID = newVal
        localStorage.setItem("portfolioID", newVal)
    },
    setPortfolioKeyValue(newVal) {
        console.log('Setting "portfolioKeyValue": ', newVal)
        state.portfolioKeyValue = newVal
        localStorage.setItem("portfolioKeyValue", newVal)
    }
}