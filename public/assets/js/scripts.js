$(function () {

    /**
     * Append containers to table.
     *
     * @param {array} containers
     */
    function appendContainers(containers) {
        var $sites = $('.sites table tbody');
        var keys = ['Id', 'Image', 'State', 'Status'];

        for (var i = 0, l = containers.length; i < l; i++) {
            var container = containers[i];
            var url = container.Names[0].substr(1);

            var $tr = $sites.find('tr[data-container-id="' + container.Id + '"]');

            if ($tr.size()) {
                for (var key in container) {
                    if (keys.indexOf(key) === -1) {
                        continue;
                    }

                    if (key == 'Id') {
                        container[key] = container[key].substr(0, 12);
                    }

                    $tr.find('td.container-' + key.toLowerCase()).text(container[key]);
                }
            } else {
                var html = [
                    '<tr data-container-id="' + container.Id + '">',
                        '<td class="container-url"><a href="//' + url + '" target="_blank">' + url + '</a></td>'
                ];

                for (var key in container) {
                    if (keys.indexOf(key) === -1) {
                        continue;
                    }

                    if (key == 'Id') {
                        container[key] = container[key].substr(0, 12);
                    }

                    html.push('<td class="container-' + key.toLowerCase() + '">' + container[key] + '</td>');
                }

                html.push('<td class="container-actions"><a href="#" class="restart">restart</a><a href="#" class="delete">delete</a></td>');
                html.push('</tr>');

                $sites.append(html.join(''))
            }
        }

        $('.loader').hide();
    }

    /**
     * Update containers.
     */
    function updateContainers() {
        $('.loader').show();
        $.getJSON('/api/containers', appendContainers);
    }
    updateContainers();
    setInterval(updateContainers, 300000);

    // Fetch configuration.
    $.getJSON('/api/config', function (res) {
        $('.domain').html('<a href="' + res.Domain + '">' + res.Domain + '</a>');
    });

    // Delete a container.
    $(document.body).on('click', '.container-actions .delete', function(e) {
        e.preventDefault();

        var $this = $(this);
        var domain = $this.closest('tr').find('.container-url').text();
        var result = prompt('Type "delete" to confirm delete of container');

        if (result !== 'delete') {
            return;
        }

        $.getJSON('/api/containers?action=delete&domain=' + domain, function(res) {
            $('.loader').hide();

            if (res.success) {
                $this.closest('tr').remove();
            }
        });
    });

    // Restart a container.
    $(document.body).on('click', '.container-actions .restart', function(e) {
        e.preventDefault();

        var $this = $(this);
        var domain = $this.closest('tr').find('.container-url').text();

        $('.loader').show();

        $.getJSON('/api/containers?action=restart&domain=' + domain, appendContainers);
    });
});