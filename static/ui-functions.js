
// Makes the script run only after html content is loaded
document.addEventListener("DOMContentLoaded", function () {
    const theme = localStorage.getItem("theme");
    const icon = document.getElementById('themeIcon');
    if (theme === "dark") {
        document.body.classList.add("dark-mode");
        icon.textContent = 'light_mode';
    } else {
        icon.textContent = 'dark_mode';
    }
});

// Change between light and dark mode
function toggleDarkMode() {
    const body = document.body;
    const icon = document.getElementById('themeIcon');
    body.classList.toggle("dark-mode");
    const currentTheme = body.classList.contains("dark-mode") ? "dark" : "light";
    localStorage.setItem("theme", currentTheme);
    icon.textContent = currentTheme === "dark" ? "light_mode" : "dark_mode";
}

// Scroll to the top of the page when the button is clicked
function scrollToTop() {
    window.scrollTo({
        top: 0,
        behavior: "smooth"
    });
}
