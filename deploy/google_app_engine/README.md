# go-lazycache-appengine

This is the package for deploying [go-lazycache-app](https://github.com/amarburg/go-lazycache-app) through [Google App Engine.](https://cloud.google.com/appengine/)   We use the custom runtime functinoality where the deployment package is basically a Dockerfile plus the vendor-specific `app.yaml`  configuration file.

As it's all Docker-based, there actually no Go code in this repo.  
