### Create Your Databases

    $ buffalo db create -a

### Run Migration

    $ buffalo db migrate

## Set Application Service Keys

    $ cp .env.example .env

Replace the values for each key in the `.env` file with the applicable values

## Starting the Application

    $ PORT=4000 buffalo dev

If you point your browser to [http://127.0.0.1:4000](http://127.0.0.1:4000) you should see a "Welcome to Text Campaign!" page.

## Heroku

### Set Environment Variables

    $ heroku config:set GO_ENV=production
    $ heroku config:set TWILIO_USER=abcdef
    $ heroku config:set TWILIO_PW=abcdef
    $ heroku config:set TWILIO_AC_SID=abcdef
    $ heroku config:set TWILIO_AUTH_TOKEN=abcdef
    $ heroku config:set BASE_URL=f7da2e67.ngrok.io
    $ heroku config:set MASTER_USERNAME=master@test.com
    $ heroku config:set MASTER_PASSWORD=ffffff

### Deployment

    $ heroku container:login
    $ heroku container:push web
    $ heroku run /bin/app migrate
    $ heroku run /bin/app task db:seed
    $ heroku run /bin/app task db:reset

[Powered by Buffalo](http://gobuffalo.io)
