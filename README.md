# push

Once a webhook has been created, you can:
- Configure a JQ filter to be run against all request payloads before persisting
  - The pre transform payload can be persisted as well by setting the `preserve_payload` field to `true` when registering the webhook.
  - [ ] Conditional transforms
- Inspect all payloads the webhook has received
  - [ ] Filtering, pagination
- Configure a the webhook to forward requests to a defined URL.
  - The data that gets forwarded can be either pre or post transform
  - [ ] conditional forwarding
