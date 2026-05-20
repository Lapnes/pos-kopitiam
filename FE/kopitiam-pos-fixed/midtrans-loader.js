(function() {
  // Fetch Midtrans config from API
  fetch('/api/v1/config/midtrans')
    .then(response => {
      if (!response.ok) {
        throw new Error('Failed to fetch Midtrans config');
      }
      return response.json();
    })
    .then(config => {
      if (!config.client_key) {
        console.warn('Midtrans client key is empty in config');
        return;
      }
      // Determine script URL based on environment
      const scriptUrl = (config.environment === 'production' || config.environment === 'live')
        ? 'https://app.midtrans.com/snap/snap.js'
        : 'https://app.sandbox.midtrans.com/snap/snap.js';
      
      console.log(`Loading Midtrans Snap SDK from ${scriptUrl} with client key: ${config.client_key}`);
      
      const script = document.createElement('script');
      script.type = 'text/javascript';
      script.src = scriptUrl;
      script.setAttribute('data-client-key', config.client_key);
      document.head.appendChild(script);
    })
    .catch(error => {
      console.error('Error loading Midtrans configuration:', error);
    });
})();
