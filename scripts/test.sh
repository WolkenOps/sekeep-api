#!/usr/bin/env bash

token="set-me-please"
header="Authorization: bearer $token"
url="set-me-please"

declare -a data=(
    '{"name":"/test/variable", "value": "hey"}'
    '{"name":"/test/variable", "value": "bye"}'
    '{"name":"/test/variable"}'
    '{"name":"/test/variable"}'
)

i=0
name="/work/google/myemail3"
value="mysecretpassword"
valuen="myotherpassissafer"


### Create a new password
r=$(curl -s -XPOST -d "{\"name\":\"$name\", \"value\":\"$value\"}" -H "Authorization: bearer $token" $url)
if [ "$r" = "$name" ] ; then echo "create new password succeeded" ; else echo "create new password failed" ; exit 1 ; fi

### Get the created password
r=$(curl -s -XGET -H "Authorization: bearer $token" $url?name=$name)
if [ "$r" = "$value" ] ; then echo "get password succeeded" ; else echo "get password failed" ; exit 1 ; fi

### Validate that password already exists
e="Password already exists"
r=$(curl -s -XPOST -d "{\"name\":\"$name\", \"value\":\"$value\"}" -H "Authorization: bearer $token" $url)
if [ "$r" = "$e" ] ; then echo "validate password exists succeeded" ; else echo "validate password exists failed" ; exit 1 ; fi

### Update the password
r=$(curl -s -XPATCH -d "{\"name\":\"$name\", \"value\":\"$valuen\"}" -H "Authorization: bearer $token" $url)
if [ "$r" = "$name" ] ; then echo "update password succeeded" ; else echo "update password failed" ; exit 1 ; fi

### Validate new value
r=$(curl -s -XGET -H "Authorization: bearer $token" $url?name=$name)
if [ "$r" = "$valuen" ] ; then echo "validate new password succeeded" ; else echo "validate new password failed" ; exit 1 ; fi

### Validate password deletion
r=$(curl -s -XDELETE -d "{\"name\":\"$name\"}" -H "Authorization: bearer $token" $url)
if [ "$r" = "$name" ] ; then echo "delete password succeeded" ; else echo "delete password failed" ; exit 1 ; fi

