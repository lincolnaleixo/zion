document.addEventListener('DOMContentLoaded', () => {
    const serverList = document.getElementById('server-list');
    const logsTableBody = document.querySelector('#logs-table tbody');
    const levelFilter = document.getElementById('level-filter');
    const sortOrder = document.getElementById('sort-order');

    let selectedServer = '';

    // Fetch and display servers
    fetch('/api/servers')
        .then(response => response.json())
        .then(servers => {
            servers.forEach(server => {
                const li = document.createElement('li');
                li.textContent = server;
                li.addEventListener('click', () => {
                    document.querySelectorAll('#server-list li').forEach(item => item.classList.remove('active'));
                    li.classList.add('active');
                    selectedServer = server;
                    fetchLogs();
                });
                serverList.appendChild(li);
            });
        })
        .catch(error => console.error('Error fetching servers:', error));

    // Fetch and display logs
    function fetchLogs() {
        let url = '/api/logs?';

        if (levelFilter.value) {
            url += `level=${encodeURIComponent(levelFilter.value)}&`;
        }
        if (selectedServer) {
            url += `server=${encodeURIComponent(selectedServer)}&`;
        }
        if (sortOrder.value) {
            url += `sort=${encodeURIComponent(sortOrder.value)}&`;
        }

        fetch(url)
            .then(response => response.json())
            .then(logs => {
                logsTableBody.innerHTML = '';
                logs.forEach(log => {
                    const tr = document.createElement('tr');

                    const timestampTd = document.createElement('td');
                    timestampTd.textContent = new Date(log.timestamp).toLocaleString();
                    tr.appendChild(timestampTd);

                    const levelTd = document.createElement('td');
                    levelTd.textContent = log.level;
                    tr.appendChild(levelTd);

                    const serverTd = document.createElement('td');
                    serverTd.textContent = log.server_name;
                    tr.appendChild(serverTd);

                    const appTd = document.createElement('td');
                    appTd.textContent = log.application;
                    tr.appendChild(appTd);

                    const messageTd = document.createElement('td');
                    messageTd.textContent = log.message;
                    tr.appendChild(messageTd);

                    const errorCodeTd = document.createElement('td');
                    errorCodeTd.textContent = log.error_code || '-';
                    tr.appendChild(errorCodeTd);

                    logsTableBody.appendChild(tr);
                });
            })
            .catch(error => console.error('Error fetching logs:', error));
    }

    // Event listeners for filters
    levelFilter.addEventListener('change', fetchLogs);
    sortOrder.addEventListener('change', fetchLogs);
});