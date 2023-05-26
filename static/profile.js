const userId = 5/* get user ID from URL or local storage */;

fetch(`/users/${userId}`)
    .then(response => response.json())
    .then(user => {
        document.getElementById('user-id').textContent = user.UserID;
        document.getElementById('user-email').textContent = user.Email;
        document.getElementById('user-name').textContent = user.Name;
        document.getElementById('user-surname').textContent = user.Surname;
        document.getElementById('user-password').textContent = user.Password;
        document.getElementById('user-gender').textContent = user.Gender;
        document.getElementById('user-token').textContent = user.Token;
        document.getElementById('user-token-date').textContent = user.TokenDate;
        document.getElementById('user-registration-date').textContent = user.RegistrationDate;
        document.getElementById('user-role').textContent = user.Role;

        const purchaseTable = document.getElementById('purchase-history');
        user.Purchases.forEach(purchase => {
            const row = purchaseTable.insertRow();
            row.insertCell().textContent = purchase.ProductName;
            row.insertCell().textContent = purchase.ProductPrice;
            row.insertCell().textContent = purchase.PurchaseDate;
        });
    })
    .catch(error => console.error(error));