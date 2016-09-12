$(function () {
    // Fetch configuration.
    $.getJSON('/api/config', function (res) {
        $('.domain').html('<a href="' + res.Domain + '">' + res.Domain + '</a>');
    });

    // Fetch all containers.
    $.getJSON('/api/containers', function (res) {
        for (var i = 0, l = res.length; i < l; i++) {
            var container = res[i];
            var url = container.Names[0].substr(1);

            var html = [
                '<tr>',
                '<td><a href="//' + url + '" target="_blank">' + url + '</a></td>',
                '<td>' + container.Id.substr(0, 12) + '</td>',
                '<td>' + container.Image + '</td>',
                '<td>' + container.State + '</td>',
                '<td>' + container.Status + '</td>',
                '</tr>'
            ];

            $('.sites table tbody').append(html.join(''))
        }
    });
});