function validateEmail(email) {
    const res = /^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$/;
    return res.test(String(email).toLowerCase());
}

function login() {
    let regform = document.getElementById('login-form');

    if(!validateEmail(regform['email'].value)){
        return
    }

    fetch('/api', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            query: `
query {
  login(
    email: \"${regform['email'].value}\",
    password: \"${regform['password'].value}\"
  )
}`
        })
    }).then(r => r.json()).then(data => {
        window.location.reload();
    });
}

function register() {
    let regform = document.getElementById('register-form');
    if (regform['new-password'].value !== regform['confirm-password'].value){
        alert("Passwords aren't same!")
        return;
    }

    if(!validateEmail(regform['new-email'].value)){
        return
    }

    fetch('/api', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            query: `
query {
  register(
    email: \"${regform['new-email'].value}\",
    name: \"${regform['new-name'].value}\",
    surname: \"${regform['new-surname'].value}\",
    gender: \"${regform['new-gender'].value}\",
    password: \"${regform['new-password'].value}\"
  )
}`
        })
    }).then(r => r.json()).then(data => {
        window.location.reload();
    });
}