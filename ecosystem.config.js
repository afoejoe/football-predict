module.exports = {
    apps: [{
      name: "predict.dev",
      "script": "make run",
      "exec_interpreter": "none",
      "exec_mode": "",
      "env": {
          "NODE_ENV": "dev",
          "PORT": 4000,
          "HOST": "127.0.0.1",
          "ORIGIN": "http://api.naijaanswers.com"
      }
    }],
    deploy: {
      // "dev" is the environment name
      dev: {
        // SSH key path, default to $HOME/.ssh
        key: "deploy.key",
        // SSH user
        user: "deploy",
        // SSH host
        host: "156.227.0.205",
        // SSH options with no command-line flag, see 'man ssh'
        // can be either a single string or an array of strings
        ssh_options: "StrictHostKeyChecking=no",
        // GIT remote/branch
        ref: "origin/dev",
        // GIT remote
        repo: "git@github.com:afoejoe/football-predict.git",
        // path in the server
        path: "/home/deploy/web/predict-test/apps/development",
        // Pre-setup command or path to a script on your local machine
        'pre-setup': "ls -la && echo 'pre set up'",
        // Post-setup commands or path to a script on the host machine
        // eg: placing configurations in the shared dir etc
        'post-setup': "ls -la",
        // pre-deploy action
        'pre-deploy-local': "echo 'This is a local executed command'",
        // post-deploy action
        "post-deploy": "pwd && ls -la && touch .envrc && > .envrc && echo DB_DSN=$DEV_DB_DSN >> .envrc && echo HTTP_PORT=$DEV_HTTP_PORT >> .envrc && echo BASIC_AUTH_HASHED_PASSWORD=$BASIC_AUTH_HASHED_PASSWORD >> .envrc  && echo NOTIFICATIONS_EMAIL=$NOTIFICATIONS_EMAIL >> .envrc && make build && pm2 reload ecosystem.config.js --env dev && pm2 save",
        "env": {
            "DEV_DB_DSN": process.env.DEV_DB_DSN,
            "DEV_HTTP_PORT": process.env.DEV_HTTP_PORT,
            "BASIC_AUTH_HASHED_PASSWORD": process.env.BASIC_AUTH_HASHED_PASSWORD,
            "NOTIFICATIONS_EMAIL": process.env.NOTIFICATIONS_EMAIL,
        }
      },
    }
  }