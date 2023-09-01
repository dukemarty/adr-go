# adr-go
Simple tool for handling ADRs (Architecture Decision Records), inspired by the original adr-tools.


## Developer information

### Build release version

Based on [this stackoverflow article](https://stackoverflow.com/questions/29599209/how-to-build-a-release-version-binary-in-go), use this command line to strip debug and symbol information during build:

    go build -ldflags "-s -w"

