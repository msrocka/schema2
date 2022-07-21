# schema2

This is an experimental repository for generating the artifacts for olca-schema
version 2.

## The `osch` tool

### Building

```bash
cd oschgo ; go build -o ../osch ; cd ..
```

or on Windows:

```batch
cd oschgo && go build -o ..\osch.exe && cd ..
```

### Usage

```
usage:

$ osch [command] [options]

commands:

  check  - checks the schema
  doc    - generates the schema documentation
  help   - prints this help
  proto  - converts the schema to ProtocolBuffers
  python - generates a Python class model for the schema

```

