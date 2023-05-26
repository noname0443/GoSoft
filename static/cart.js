$(document).ready(function() {
    // Retrieve cart data from server
    $.ajax({
        url: "/api/cart",
        type: "GET",
        dataType: "json",
        success: function(data) {
            // Update table with cart data
            var tbody = $("#cart-table tbody");
            tbody.empty();
            var total = 0;
            $.each(data.items, function(index, item) {
                var row = $("<tr>")
                    .append($("<td>").text(item.name))
                    .append($("<td>").text("$" + item.price.toFixed(2)))
                    .append($("<td>").text(item.quantity))
                    .append($("<td>").text("$" + (item.price * item.quantity).toFixed(2)));
                tbody.append(row);
                total += item.price * item.quantity;
            });
            $("#cart-total").text("$" + total.toFixed(2));
        },
        error: function(jqXHR, textStatus, errorThrown) {
            console.error("Error retrieving cart data:", errorThrown);
        }
    });
});