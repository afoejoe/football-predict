module.exports = {
    apps: [{
      name: "predict.dev",
      "script": "make run",
    }],
    deploy: {
      // "dev" is the environment name
      dev: {
        // SSH key path, default to $HOME/.ssh
        key: "/Users/josh/.ssh/football_predict",
        // SSH user
        user: "deploy",
        // SSH host
        host: ["156.227.0.205"],
        // SSH options with no command-line flag, see 'man ssh'
        // can be either a single string or an array of strings
        ssh_options: "StrictHostKeyChecking=no",
        // GIT remote/branch
        ref: "origin/dev",
        // GIT remote
        repo: "git@github.com:afoejoe/football-predict.git",
        // path in the server
        path: "~/web/predict-test/apps/development/",
        // Pre-setup command or path to a script on your local machine
        'pre-setup': "ls -la && echo 'pre set up'",
        // Post-setup commands or path to a script on the host machine
        // eg: placing configurations in the shared dir etc
        'post-setup': "ls -la",
        // pre-deploy action
        'pre-deploy-local': "echo 'This is a local executed command'",
        // post-deploy action
        "post-deploy": "touch .envrc && > .envrc && echo DB_DSN=$DEV_DB_DSN >> .envrc && echo HTTP_PORT=$DEV_HTTP_PORT >> .envrc && echo basic-auth-hashed-password=$basic-auth-hashed-password >> .envrc  && echo notifications-email=$notifications-email >> .envrc && make build && pm2 reload ecosystem.config.js --env dev && pm2 save",
        "env": {
            "DEV_DB_DSN": process.env.DEV_DB_DSN,
            "DEV_HTTP_PORT": process.env.DEV_HTTP_PORT,
            "basic-auth-hashed-password": process.env.basic-auth-hashed-password,
            "notifications-email": process.env.notifications-email,
        }
      },
    }
  }