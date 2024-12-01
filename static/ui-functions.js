    // Makes the script run only after html content is loaded
    document.addEventListener("DOMContentLoaded", function () {  
        const theme = localStorage.getItem("theme");
        const button = document.getElementById('buttonText');
        if (theme === "dark") {
            document.body.classList.add("dark-mode");
            button.textContent = 'Light mode';
        } else {
            button.textContent = 'Dark mode';
        }
    });

    /* Toggle dark mode and save the preference */
    function toggleDarkMode() {
        const body = document.body;
        const button = document.getElementById('buttonText');
        body.classList.toggle("dark-mode");

        // Save the current theme in localStorage
        const currentTheme = body.classList.contains("dark-mode") ? "dark" : "light";
        localStorage.setItem("theme", currentTheme);
        button.textContent = currentTheme === "dark" ? "Light mode" : "Dark mode";
    }

    /* Make all checkboxes, number inputs and dropdowns submit the form */
    document.addEventListener("DOMContentLoaded", function () {   // Makes the script run only after html content is loaded
        const inputs = document.querySelectorAll("input[type='checkbox'], input[type='number'], select");
        inputs.forEach(input => {
            input.addEventListener("change", function () {
                this.form.submit();
            });
        });
    });

    /* "checkAll" checks all countries  */
    document.getElementById('checkAll').addEventListener('click', () => {
        const countryChecks = document.querySelectorAll('.countryCB');
        countryChecks.forEach(checkbox => {
            checkbox.checked = true;
        });
    });

    /* "uncheckAll" unchecks all countries  */
    document.getElementById('uncheckAll').addEventListener('click', () => {
        const countryChecks = document.querySelectorAll('.countryCB');
        countryChecks.forEach(checkbox => {
            checkbox.checked = false;
        });
    });

    /* expand the collapsible and reveal the hidden info box, and the reverse */
    var coll = document.getElementsByClassName("collapsible");
    var i;
    for (i = 0; i < coll.length; i++) {
        coll[i].addEventListener("click", function () {
            this.classList.toggle("active");
            var content = this.nextElementSibling;
            if (content.style.maxHeight) {
                content.style.maxHeight = null;
                /* this.style.height = "320px"; */
                this.style.borderRadius = "10px";
                this.children[0].style.width = "240px";
            } else {
                content.style.maxHeight = content.scrollHeight + "px";
                /* this.style.height = "320px"; */
                this.style.borderRadius = "10px 10px 0px 0px";
                this.children[0].style.width = "240px";
            }
        });
    }