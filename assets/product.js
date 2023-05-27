function submitComment() {
    fetch('/api', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            query: `query { CommentAdd(
                productid: ${parseInt(window.location['pathname'].substr(window.location['pathname'].lastIndexOf('/')+1))},
                content: \"${document.getElementById('comment').value}\") }`
        })
    }).then(r => r.json()).then(data => {
        console.log(data)
        window.location.reload();
    });
}

function addToCart(){
    fetch('/api', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            query: `query { CartAdd(
                productid: ${parseInt(window.location['pathname'].substr(window.location['pathname'].lastIndexOf('/')+1))},
                count: ${parseInt(document.getElementById('soft-duration').value)}) }`
        })
    }).then(r => r.json()).then(data => {
        console.log(data);
        alert("Successfully added to cart!")
    });
}