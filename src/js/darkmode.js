document.addEventListener('DOMContentLoaded', function () {
    const htmlElement = document.documentElement;
    const themeButton = document.getElementById('themeButton');
    const themeIcon = document.getElementById('themeIcon');

    // Setează tema implicită pe "dark" dacă nu există nici o valoare în localStorage
    const currentTheme = localStorage.getItem('bsTheme') || 'dark';
    htmlElement.setAttribute('data-bs-theme', currentTheme);

    // Schimbă iconița în funcție de tema activă
    function updateIcon() {
        if (htmlElement.getAttribute('data-bs-theme') === 'dark') {
            themeIcon.classList.remove('bi-sun');
            themeIcon.classList.add('bi-moon');
        } else {
            themeIcon.classList.remove('bi-moon');
            themeIcon.classList.add('bi-sun');
        }
    }

    // Schimbă tema când se apasă butonul
    themeButton.addEventListener('click', function () {
        const currentTheme = htmlElement.getAttribute('data-bs-theme');
        if (currentTheme === 'dark') {
            htmlElement.setAttribute('data-bs-theme', 'light');
            localStorage.setItem('bsTheme', 'light');
        } else {
            htmlElement.setAttribute('data-bs-theme', 'dark');
            localStorage.setItem('bsTheme', 'dark');
        }
        updateIcon(); // Actualizează iconița
    });

    // Inițializare iconiță și tema
    updateIcon(); 
});
