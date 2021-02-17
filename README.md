# KAFUTIL
----

Utility helper for kafka

Currently used for [WIP]
- Serialize protobuf message and write using descriptor set
- Read & deserialize same protobuf message

### Instructions

- Build using `make`
- Configure the utility using `.kafutil.yml`
- Modify protos/app.proto to update the proto used for reading and writing
- Make process also generates a descriptor set

- To write a random message in topic provided in .kafutil.yaml
```
./kafutil write
```

- To read messages from topic provided in .kafutil.yaml
```
./kafutil read
```

### WARNING

This is in very raw state and only made for testing