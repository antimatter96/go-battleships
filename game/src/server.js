const express = require('express');
const favicon = require('serve-favicon');
const nunjucks = require('nunjucks');

function getExpressApp(config) {
  const app = express();

  app.use(favicon(config.faviconDirectory));
  app.use(config.staticPath, express.static(config.staticDirectory));
  app.set('views', config.viewsDirectory);

  nunjucks.configure(app.get('views'), {
    ...config.nunjucks,
    express: app
  });

  app.use(function (_req, res, next) {
    res.setHeader("X-Frame-Options", "DENY");
    res.setHeader("Content-Security-Policy", "frame-ancestors 'none'");
    res.setHeader("X-powered-by", "none");
    next();
  });

  app.get('/', function (_req, res) {
    res.render('index.njk');
  });

  return app;
}

module.exports = {
  getExpressApp: getExpressApp
};
