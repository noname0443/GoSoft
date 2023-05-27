const commentForm = document.getElementById('comment-form');
const commentList = document.getElementById('comment-list');

// Function to submit comment
function submitComment() {
    fetch('/api', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            query: `query { CommentAdd(
                productid: 1,
                content: \"${document.getElementById('comment').value}\") }`
        })
    }).then(r => r.json()).then(data => {
        console.log(data)
        window.location.reload();
    });
}

// Add event listener to comment form
commentForm.addEventListener('submit', submitComment);