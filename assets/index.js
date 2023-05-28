let slideIndex = 0;
showSlides();

// Next/previous controls
function plusSlides(n) {
    slideIndex += n
    let slides = document.getElementsByClassName("mySlides");
    for (i = 0; i < slides.length; i++) {
        slides[i].style.display = "none";
    }
    if (slideIndex > slides.length) {slideIndex = 1}
    if (slideIndex <= 0) {slideIndex = slides.length }
    slides[slideIndex-1].style.display = "block";
}

// Thumbnail image controls
function currentSlide(n) {
    showSlides(slideIndex = n);
}

function showSlides() {
    let i;
    let slides = document.getElementsByClassName("mySlides");
    for (i = 0; i < slides.length; i++) {
        slides[i].style.display = "none";
    }
    slideIndex++;
    if (slideIndex > slides.length) {slideIndex = 1}
    if (slideIndex <= 0) {slideIndex = slides.length }
    slides[slideIndex-1].style.display = "block";
    setTimeout(showSlides, 10000);
}