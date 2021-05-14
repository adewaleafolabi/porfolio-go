import Vue from 'vue'
import App from './App.vue'
import router from './router'
import {DARK_THEME, LIGHT_THEME, setClassAndThemes} from "./service/theme";

Vue.config.productionTip = false

new Vue({
    router,
    render: h => h(App)
}).$mount('#app');


(function () {
    try {
        // Checks for the color scheme of the device.
        // In this case it checks for anything that is not light theme.
        const media = window.matchMedia("not all and (prefers-color-scheme: light)");
        if (localStorage.getItem("theme") == "dark") {
            setClassAndThemes(DARK_THEME);
        } else if (localStorage.getItem("theme") == "light") {
            setClassAndThemes(LIGHT_THEME);
        } else if (media.matches) {
            setClassAndThemes(DARK_THEME);
            localStorage.setItem("theme", "dark")
        }
        media.addListener(function () {
            if (media.matches) {
                setClassAndThemes(DARK_THEME);
            } else {
                setClassAndThemes(LIGHT_THEME);
            }
        });
    } catch (err) {
        console.log(err)
    }
})();
