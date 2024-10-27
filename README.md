# push

A webhook server allowing users to create a new webhook that will:
- dump the request body in a Minio bucket
  - this can be done before and/or after it has been transformed by the user defined JQ filter
- can run jq filters on the data
  - if no filter is provided then the original request body will be stored 
- can forward request's to a defined url
  - data forwarded can be either pre- or post-transform
  - if no url is provided, no attempt to forward the data is made
-----
## Makefile

-----
## Wishlist
 - [ ] clients to interact with the server (web, cli)