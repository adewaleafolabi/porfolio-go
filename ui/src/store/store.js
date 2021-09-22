//store.js
import Vue from 'vue'

export const state = Vue.observable({privacy: false})

export const mutations = {
    setPrivacy(newVal) {
        console.log('Setting "privacy": ', newVal)
        state.privacy = newVal
        localStorage.setItem("privateMode", newVal)
    }
}