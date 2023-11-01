---
title: "Waiting for an HTTP Service in GitLab CI/CD"
date: 2023-11-01T10:00:00+06:00
tags:
  - GitLab
  - 100DaysToOffload
---

Last weekend, I was setting up a Cypress test pipeline in GitLab for Toph for the 5th time.

I have no idea why this pipeline keeps breaking over time. It's like bread left in the open. Cypress is such a fantastic end-to-end testing tool. But it seems to need a lot of _extra love_.

Something that I needed to do was, in the CI/CD script, wait for Toph to start handling HTTP requests before starting Cypress.

Take the following as an example:

- Start `./server` in the background
- Wait for `./server` to start handling HTTP requests
- Run Cypress tests

In the case of Toph, I could have added a `sleep 1` after starting Toph's web server and called it a day.

But it didn't sit right with me. If the start-up time increases in the future or due to external reasons, the test will fail.

Instead, I added a one-liner in the GitLab CI/CD script that uses `until` and `curl` to wait for Toph to start handling HTTP requests.

For the example above, it would look like this:

``` yaml
test:cypress:
  script:
  - "./server &"
  - until $(curl --output /dev/null --silent --head --fail http://localhost:5000/api/ping); do printf '.'; sleep 1; done
  - npx cypress run --browser firefox
```

The second line in the script will block until `curl` receives an HTTP response from `http://localhost:5000/api/ping`.

