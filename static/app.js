document.addEventListener('DOMContentLoaded', () => {
  const app = {
      selectedServer: null,
      servers: [],
      logs: [],
      filters: {
          level: '',
          date: ''
      },

      init() {
          this.serverList = document.getElementById('server-list');
          this.logsContainer = document.getElementById('logs-container');
          this.filtersContainer = document.getElementById('filters');
          this.levelFilter = document.getElementById('levelFilter');
          this.dateFilter = document.getElementById('dateFilter');
          this.clearFiltersButton = document.getElementById('clearFilters');

          this.fetchServers();
          this.setupEventListeners();
      },

      setupEventListeners() {
          this.levelFilter.addEventListener('change', () => this.applyFilters());
          this.dateFilter.addEventListener('change', () => this.applyFilters());
          this.clearFiltersButton.addEventListener('click', () => this.clearFilters());
      },

      fetchServers() {
          fetch('/api/servers')
              .then(response => response.json())
              .then(servers => {
                  this.servers = servers.map((name, index) => ({ id: index + 1, name: name }));
                  this.renderServers();
              })
              .catch(error => console.error('Error fetching servers:', error));
      },

      renderServers() {
          this.serverList.innerHTML = '';
          this.servers.forEach(server => {
              const button = document.createElement('button');
              button.textContent = server.name;
              button.className = `w-full mb-2 px-4 py-2 text-sm font-medium rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-all duration-200 ease-in-out ${this.getButtonClasses(server.id)}`;
              button.addEventListener('click', () => this.selectServer(server.id));
              this.serverList.appendChild(button);
          });
      },

      selectServer(id) {
          this.selectedServer = id;
          this.fetchLogs();
          this.updateServerButtons();
          this.filtersContainer.style.display = 'block';
      },

      updateServerButtons() {
          const buttons = this.serverList.querySelectorAll('button');
          buttons.forEach(button => {
              const serverId = this.servers.find(server => server.name === button.textContent).id;
              button.className = `w-full mb-2 px-4 py-2 text-sm font-medium rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-all duration-200 ease-in-out ${this.getButtonClasses(serverId)}`;
          });
      },

      fetchLogs() {
          if (!this.selectedServer) return;

          let url = `/api/logs?server=${encodeURIComponent(this.servers.find(s => s.id === this.selectedServer).name)}`;

          fetch(url)
              .then(response => response.json())
              .then(logs => {
                  this.logs = logs;
                  this.renderLogs();
              })
              .catch(error => console.error('Error fetching logs:', error));
      },

      renderLogs() {
          this.logsContainer.innerHTML = '';
          this.getFilteredLogs().forEach(log => {
              const logElement = document.createElement('div');
              logElement.className = 'bg-white shadow-md rounded-lg p-4 mb-4';
              logElement.innerHTML = `
                  <div class="flex justify-between items-center mb-2">
                      <span class="text-sm text-gray-600">${this.formatDate(log.timestamp)}</span>
                      <span class="px-2 py-1 rounded-full text-xs font-semibold text-white ${this.getLevelClass(log.level)}">${log.level}</span>
                  </div>
                  <p class="text-gray-800 mb-2">${log.message}</p>
                  <div class="flex justify-between text-xs text-gray-500">
                      <span>Source: ${log.source}</span>
                      <span>User: ${log.user}</span>
                  </div>
              `;
              this.logsContainer.appendChild(logElement);
          });
      },

      getButtonClasses(id) {
          return this.selectedServer === id
              ? 'bg-indigo-600 text-white border-indigo-600'
              : 'bg-white text-gray-700 border-gray-300 hover:bg-gray-50';
      },

      getLevelClass(level) {
          switch (level) {
              case 'INFO': return 'bg-blue-500';
              case 'WARNING': return 'bg-yellow-500';
              case 'ERROR': return 'bg-red-500';
              default: return 'bg-gray-500';
          }
      },

      formatDate(dateString) {
          const options = { year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' };
          return new Date(dateString).toLocaleDateString(undefined, options);
      },

      applyFilters() {
          this.filters.level = this.levelFilter.value;
          this.filters.date = this.dateFilter.value;
          this.renderLogs();
      },

      clearFilters() {
          this.filters.level = '';
          this.filters.date = '';
          this.levelFilter.value = '';
          this.dateFilter.value = '';
          this.renderLogs();
      },

      getFilteredLogs() {
          return this.logs.filter(log => {
              const levelMatch = !this.filters.level || log.level === this.filters.level;
              const dateMatch = !this.filters.date || log.timestamp.startsWith(this.filters.date);
              return levelMatch && dateMatch;
          });
      }
  };

  app.init();
});
