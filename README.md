# form2mail

[![Build Status](https://github.com/vearutop/form2mail/workflows/test-unit/badge.svg)](https://github.com/vearutop/form2mail/actions?query=branch%3Amaster+workflow%3Atest-unit)
[![Coverage Status](https://codecov.io/gh/vearutop/form2mail/branch/master/graph/badge.svg)](https://codecov.io/gh/vearutop/form2mail)
[![GoDevDoc](https://img.shields.io/badge/dev-doc-00ADD8?logo=go)](https://pkg.go.dev/github.com/vearutop/form2mail)
[![Time Tracker](https://wakatime.com/badge/github/vearutop/form2mail.svg)](https://wakatime.com/badge/github/vearutop/form2mail)
![Code lines](https://sloc.xyz/github/vearutop/form2mail/?category=code)
![Comments](https://sloc.xyz/github/vearutop/form2mail/?category=comments)

This microservice sends emails from a form.

## Usage

Create configuration file in `.env` or set environment variables.

```
LOG_LEVEL=warn
HTTP_LISTEN_ADDR=127.0.0.1:8008

# Create recaptcha secret keys at https://www.google.com/recaptcha/admin/create.
# Optional, if RECAPTCHA_SECRET_KEY is empty, recaptcha will be disabled.
RECAPTCHA_SECRET_KEY=6Lddl<redacted>
RECAPTCHA_SITE_KEY=6Lddlic<redacted>
RECAPTCHA_V3=true

# SMTP server configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_FROM=myservice@gmail.com
SMTP_PASSWORD=monkey
SMTP_TO=customers@myservice.com

# Optional customizations.
## Subject may be a template, containing placeholders for form fields.
SMTP_SUBJECT=New form received from {{.name}}
## Reply to is the name of the form field that contains the email address.
SMTP_REPLY_TO=email
## Body template files, should contain placeholders for form fields.
SMTP_BODY_HTML_TEMPLATE_FILE=email.template.html
SMTP_BODY_TEMPLATE_FILE=email.template.txt
## Optional path to your langing page contents, `index.html` will be served at /, other resources at /static/*. 
STATIC_DIR=./testpage
```

Start the service:

```
form2mail
{"level":"info","@timestamp":"2022-04-02T01:11:45.554+0200","message":"starting server, Swagger UI at http://127.0.0.1:8008/docs"}
```

Service will serve `POST` `/receive` endpoint and will send an email with POST form parameters (including file
uploads) and URL query parameters.

Optional `success_url` and `fail_url` query/form parameters can be provided to redirect the user.