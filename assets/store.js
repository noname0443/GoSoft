const addToCartButtons = document.querySelectorAll('.add-to-cart');

addToCartButtons.forEach(button => {
    button.addEventListener('click', () => {
        alert('Item added to cart!');
    });
});

function getProducts() {
    document.getElementsByClassName('product-list')[0].innerHTML = "";
    fetch('/api', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            query: `query { search(
                name: \"${document.getElementsByName('search-input')[0].value}\",
                lower_price: ${parseFloat(document.getElementsByName('lower-price')[0].value)},
                highest_price: ${parseFloat(document.getElementsByName('highest-price')[0].value)},
                categories: \"${document.getElementsByName('category')[0].value}\"
            )
            { id, name, photo, price, description, subscriptiontype, company } }`
        })
    })
        .then(r => r.json())
        .then(data =>
        {
            console.log(data)
            let product_list = document.getElementsByClassName('product-list')[0];
            let array = data['data']['search'];
            for(let i = 0; i < array.length; i++){
                product_list.innerHTML += `
        <li class="product-item" itemid="${array[i]['id']}">
             <div class="box" >
                 <img src="${array[i]['photo']}">
                 <div id="text">
                    <h3>${array[i]['name']}</h3>
                    <p class="product-price">Price: $${array[i]['price']}/${array[i]['subscriptiontype']}</p>
                    <p>${array[i]['company']}</p>
                 </div>
             </div>
             <a class="add-to-cart" href="/store/${array[i]['id']}">Check Software</a>
        </li>
    `
                /*let textBox = product_list.querySelectorAll('[itemid]')[i].querySelector('div').querySelector('[id]').querySelector('p');
                var textBoxWidth = textBox.clientWidth;
                var textBoxPadding = parseInt(getComputedStyle(textBox).paddingLeft) + parseInt(getComputedStyle(textBox).paddingRight);
                var fontSize = parseInt(getComputedStyle(textBox).fontSize);
                var lineHeight = parseInt(getComputedStyle(textBox).lineHeight);
                var textBoxHeight = textBox.clientHeight;
                var charCount = Math.floor(((textBoxWidth - textBoxPadding) / fontSize) * (textBoxHeight / lineHeight));
                textBox.innerText = textBox.innerText.substr(0, charCount)*/
            }
        });
}