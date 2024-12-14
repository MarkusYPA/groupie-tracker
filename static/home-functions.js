
/* expand the collapsible and reveal the hidden info box, and the reverse */
var coll = document.getElementsByClassName("collapsible");
var i;
for (i = 0; i < coll.length; i++) {
    coll[i].addEventListener("click", function () {
        this.classList.toggle("active");
        var content = this.nextElementSibling;
        if (content.style.maxHeight) {
            content.style.maxHeight = null;
            this.style.borderRadius = "10px";
        } else {
            content.style.maxHeight = content.scrollHeight + "px";
            this.style.borderRadius = "10px 10px 0px 0px";
        }
    });
}

/* Make all checkboxes, number inputs and dropdowns submit the form */
document.addEventListener("DOMContentLoaded", function () {   // Makes the script run only after html content is loaded
    const inputs = document.querySelectorAll("input[type='checkbox'], input[type='number'], input[type='range'], select");
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

/* "checkAllLoc" checks all locales  */
document.getElementById('checkAllLoc').addEventListener('click', () => {
    const localeChecks = document.querySelectorAll('.localeCB');
    localeChecks.forEach(checkbox => {
        checkbox.checked = true;
    });
});

/* "uncheckAllLoc" unchecks all locales  */
document.getElementById('uncheckAllLoc').addEventListener('click', () => {
    const localeChecks = document.querySelectorAll('.localeCB');
    localeChecks.forEach(checkbox => {
        checkbox.checked = false;
    });
});

// Restore scroll positions for all scrollable areas
function restoreScrollPositions() {
    document.querySelectorAll("[data-scrollable]").forEach((element) => {
        let key = element.getAttribute("data-scrollable");
        let storedScrollPosition = localStorage.getItem(`scrollPos-${key}`);
        if (storedScrollPosition !== null) {
            element.scrollTop = Number(storedScrollPosition);
        }
    });
}

/* Check all locales when a country is clicked  */
document.addEventListener("DOMContentLoaded", function () {   // Run only after html content is loaded
    const countryCBs = document.querySelectorAll('.countryCB');
    countryCBs.forEach(checkbox => {
        checkbox.addEventListener('click', () => {
            const localeChecks = document.querySelectorAll('.localeCB');
            localeChecks.forEach(checkbox => {
                checkbox.checked = true;
            });
        });
    });
    document.getElementById('checkAll').addEventListener('click', () => {
        const localeChecks = document.querySelectorAll('.localeCB');
        localeChecks.forEach(checkbox => {
            checkbox.checked = true;
        });
    });
});

// Function to save scroll positions for all scrollable areas
function saveScrollPositions() {
    document.querySelectorAll("[data-scrollable]").forEach((element) => {
        let key = element.getAttribute("data-scrollable");
        localStorage.setItem(`scrollPos-${key}`, element.scrollTop);
    });
}

// Restore scroll positions when the page loads
document.addEventListener("DOMContentLoaded", restoreScrollPositions);

// Save scroll positions before the page is unloaded
window.addEventListener("beforeunload", saveScrollPositions);
