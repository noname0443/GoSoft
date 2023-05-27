// When the button is clicked, toggle the menu
document.getElementById("menu-toggle").addEventListener("click", function() {
    var menu = document.getElementById("mobile_menu");
    if (menu.style.display === "block") {
        menu.style.display = "none";
    } else {
        menu.style.display = "block";
    }
});