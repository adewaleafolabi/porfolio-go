export const DARK_THEME = "dark";
export const LIGHT_THEME = "light";
const THEMES = {
    "dark": 'https://cdn.jsdelivr.net/npm/bootstrap-dark-5@1.0.0/dist/css/bootstrap-night.min.css',
    "light": 'https://cdn.jsdelivr.net/npm/bootstrap-dark-5@1.0.0/dist/css/bootstrap.min.css'
};

export const setClassAndThemes = (theme) => {
    // if (!(theme === DARK_THEME && theme === LIGHT_THEME)) {
    //     theme = DARK_THEME
    // }
    const body = document.body;
    body.classList.remove("bg-dark", "bg-light");
    body.classList.add(`bg-${theme}`);
    const head = document.body.parentElement.firstElementChild;
    try {
        head.removeChild(createStylesheet(DARK_THEME, THEMES))
    } catch (e) {
        console.log(e)
    }
    try {
        head.removeChild(createStylesheet(LIGHT_THEME, THEMES))
    } catch (e) {
        console.log(e)
    }

    head.appendChild(createStylesheet(theme, THEMES))

}

function createStylesheet(theme,themeList) {
    const link = document.createElement('link');
    link.setAttribute('href', themeList[theme]);
    link.setAttribute('id', theme);
    link.setAttribute('rel', 'stylesheet');
    return link
}