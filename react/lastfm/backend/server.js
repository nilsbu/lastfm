const express = require('express');
const httpProxy = require('http-proxy');
const path = require('path');
const log4js = require('log4js');

// Configure the logger
log4js.configure({
  appenders: { myLogger: { type: 'console' } },
  categories: { default: { appenders: ['myLogger'], level: 'info' } }
});

// Create a logger
const logger = log4js.getLogger();

// Log a message
logger.info('Starting server...');

const app = express();
const apiProxy = httpProxy.createProxyServer();

// Serve static files from the public directory
app.use(express.static(path.join(__dirname, 'public')));

// Forward requests starting with /json/ to the lastfm-srv service
const target = process.env.BACKEND || 'localhost:3001';
app.use('/json', (req, res) => {
    logger.info(`Forwarding request to http://${target}/json${req.url}`);
    apiProxy.web(req, res, { target: `http://${target}/json` });
});


// Serve the static files from the React app
app.use(express.static(path.join(__dirname, '../build')));

// Handle requests that fall through to index.html
app.get('*', (req, res) => {
    res.sendFile(path.join(__dirname, '../build', 'index.html'));
});

// Start the server
app.listen(3000, () => {
    logger.info('Server started on port 3000');
});
