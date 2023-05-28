function removeProduct(productid) {
    fetch('/api', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            query: `query { StoreRemove(productid: ${productid})}`
        })
    })
        .then(r => r.json())
        .then(data => {
            window.location.reload()
        })
}

function updateProduct(productid){
    fetch('/api', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            query: `query { StoreUpdate(
            productid: ${productid},
                product: {
                    name: \"${document.getElementsByName('product-name')[0].value}\",
                    description: \"${document.getElementsByName('product-description')[0].value}\",
                    photo: \"${document.getElementsByName('product-photo')[0].value}\",
                    file: \"${document.getElementsByName('product-file')[0].value}\",
                    price: ${parseFloat(document.getElementsByName('product-price')[0].value)},
                    subscriptiontype: \"${document.getElementsByName('product-subscriptiontype')[0].value}\"
                }
            )}`
        })
    })
        .then(r => r.json())
        .then(data => {
            window.location.reload()
        })
}

function createProduct(){
    fetch('/api', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            query: `query { StoreAdd(
                product: {
                    name: \"${document.getElementsByName('name')[0].value}\",
                    description: \"${document.getElementsByName('description')[0].value}\",
                    photo: \"${document.getElementsByName('photo')[0].value}\",
                    file: \"${document.getElementsByName('file')[0].value}\",
                    price: ${parseFloat(document.getElementsByName('price')[0].value)},
                    subscriptiontype: \"${document.getElementsByName('subscriptiontype')[0].value}\"
                }
            )}`
        })
    })
        .then(r => r.json())
        .then(data => {
            window.location.reload()
        })
}

const overlay = document.getElementById('overlay');

function displayProduct(productid) {
    document.getElementById('accept-edit').setAttribute('editid', productid)
    overlay.style.display = 'block';
    let rootElem = document.getElementById(productid.toString())
    let tdArray = rootElem.getElementsByTagName('td')

    //document.getElementsByName('product-id')[0].value = tdArray[0].innerText
    document.getElementsByName('product-name')[0].value = tdArray[1].innerText
    document.getElementsByName('product-description')[0].value = tdArray[2].innerText
    document.getElementsByName('product-photo')[0].value = tdArray[3].innerText
    document.getElementsByName('product-file')[0].value = tdArray[4].innerText
    document.getElementsByName('product-price')[0].value = tdArray[5].innerText
    document.getElementsByName('product-subscriptiontype')[0].value = tdArray[6].innerText
}