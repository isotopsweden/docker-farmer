$(function () {

    /**
     * Append containers to table.
     *
     * @param {array} containers
     */
    function appendContainers(containers) {
        var $sites = $('.sites table tbody');
        var keys = ['Id', 'Image', 'State', 'Status'];
        var keep = [];

        for (var i = 0, l = containers.length; i < l; i++) {
            var container = containers[i];
            var url = container.Names[0].substr(1);

            var $tr = $sites.find('tr[data-container-id="' + container.Id + '"]');

            if ($tr.size()) {
                keep.push(container.Id);

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

                keep.push(container.Id);

                for (var key in container) {
                    if (keys.indexOf(key) === -1) {
                        continue;
                    }

                    if (key == 'Id') {
                        container[key] = container[key].substr(0, 12);
                    }

                    html.push('<td class="container-' + key.toLowerCase() + '">' + container[key] + '</td>');
                }

                var links = [];

                if (typeof farmer !== 'undefined') {
                    // replace {id} if id is found in url.
                    for (var key in farmer.links) {
                        var id = /(\w+\-\d+)/.exec(url);
                        var link = farmer.links[key];

                        if (!id || !id.length) {
                            continue;
                        }

                        links.push('<a class="btn" href="' + link.replace('{id}', id[0].toUpperCase()) + '" target="_blank">' + key + '</a>');
                    }

                    // replace url for all links that contains {url}.
                    for (var key in farmer.links) {
                        var link = farmer.links[key];

                        if (link.indexOf('{url}') === -1) {
                            continue;
                        }

                        links.push('<a class="btn" href="' + link.replace('{url}', url) + '" target="_blank">' + key + '</a>');
                    }
                }

                html.push('<td class="container-actions"><a href="#" class="restart">restart</a><a href="#" class="delete">delete</a>' + links.join('') + '</td>');
                html.push('</tr>');

                $sites.append(html.join(''))
            }
        }

        if (keep.length > 1 || !window.all) {
            $sites.find('tr').each(function () {
                var $this = $(this);

                if (keep.indexOf($this.data('container-id')) === -1) {
                    $this.remove();
                }
            });
        }

        $('.loader').hide();
    }

    /**
     * Update containers.
     */
    function updateContainers() {
        window.all = typeof window.all === 'undefined' ? false : window.all;

        $('.loader').show();

        $.getJSON('/api/containers?all=' + window.all, appendContainers);
    }
    updateContainers();
    setInterval(updateContainers, 300000);

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

    // Show all/less.
    $('.show-all').on('click', function(e) {
        e.preventDefault();

        var $this = $(this);
        window.all = $this.text() === 'Show all';
        updateContainers();

        var text = $this.attr('data-text');
        $this.attr('data-text', $this.text()).text(text);
    });
});