const express = require('express');
const httpProxy = require('http-proxy');
const path = require('path');
const log4js = require('log4js');
const https = require('https');
const fs = require('fs');

// Configure the logger
log4js.configure({
  appenders: { myLogger: { type: 'console' } },
  categories: { default: { appenders: ['myLogger'], level: 'info' } }
});

// Load the TLS certificate and private key
const options = {
  key: fs.readFileSync(`/cert/${process.env.TLS_CERT_REL}/privkey.pem`),
  cert: fs.readFileSync(`/cert/${process.env.TLS_CERT_REL}/cert.pem`)
};

// Create a logger
const logger = log4js.getLogger();

// Log a message
logger.info('Starting server...');

const app = express();
const apiProxy = httpProxy.createProxyServer();

// Serve static files from the public directory
app.use(express.static(path.join(__dirname, 'public')));

// Forward requests starting with /json/ to the lastfm-srv service
const target = process.env.BACKEND_HOST || 'localhost:3001';
app.use('/json', (req, res) => {
  logger.info(`Forwarding request to http://${target}/json${req.url}`);
  apiProxy.web(req, res, { target: `http://${target}/json` });
});

// // Redirect HTTP traffic to HTTPS
// app.use((req, res, next) => {
//   if (req.secure) {
//     // Request is already secure
//     next();
//   } else {
//     // Redirect to HTTPS
//     res.redirect(`https://${req.hostname}${req.url}`);
//   }
// });

// Serve the static files from the React app over HTTPS
app.use(express.static(path.join(__dirname, '../build')));

// Handle requests that fall through to index.html over HTTPS
app.get('*', (req, res) => {
  res.sendFile(path.join(__dirname, '../build', 'index.html'));
});

// Start the HTTPS server on port 4000
https.createServer(options, app).listen(4000, () => {
  logger.info('HTTPS server started on port 4000');
});

// Start the HTTP server on port 3000
app.listen(3000, () => {
  logger.info('HTTP server started on port 3000');
});
