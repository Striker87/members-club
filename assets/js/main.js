$(document).ready(function() {
    const emailPattern = /^([a-z0-9_\.-])+@[a-z0-9-]+\.([a-z]{2,4}\.)?[a-z]{2,4}$/i;
    const namePattern = /^[a-zA-Z\.\s]+$/

    $('#clear-btn').on('click', function () {
        $('#name, #email').val('');
    });

    $('#add-btn').on('click', function () {
       const name = $.trim($('#name').val());
       const email = $.trim($('#email').val());

        if(!name || !namePattern.test(name)) {
            $('#name').focus();
            alert('Enter correct name');
            return false;
        }

        if(!emailPattern.test(email))
        {
            $('#email').focus();
            alert('Enter correct email');
            return false;
        }

        $.ajax({
            type: 'POST',
            dataType: 'json',
            url: '/add_member',
            data: JSON.stringify({
                'name': name,
                'email': email
            }),
            success: function (data, textStatus, xhr) {
                if(xhr.status === 200) {
                    location.reload();
                }
            }, error: function(xhr, ajaxOptions, thrownError) {
                alert(JSON.parse(xhr.responseText).error);
            }
        });
    });
});