const links = document.querySelectorAll(".header__link");

links.forEach(link => {
    const url = new URL(link.href);

    if (url.pathname === document.body.dataset.path) {
        link.classList.add("header__link--current");
    }
});

if (document.forms["delete"] !== undefined) {
    const email = document.forms["delete"].email;
    const buttons = document.querySelectorAll("[data-email]");

    buttons.forEach(button => {
        button.addEventListener("click", () => {
            email.value = button.dataset.email;
        });
    });
}

if (document.forms["config"] !== undefined) {
    const domain = document.forms["config"].domain;

    domain.setSelectionRange(domain.value.length, domain.value.length);
}
