function removeFromCart(productid, count){
    fetch('/api', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            query: `query { CartRemove(
                productid: ${productid},
                count: ${count} ) }`
        })
    }).then(r => r.json()).then(data => {
        window.location.reload();
    });
}