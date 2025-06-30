Ring ding ding, here is your package!

# WhatDoesTheFoxSay

Small Foxpost tracking CLI. Wanted something more, but for me this will be enoguth to quickly check the status of my packages without additional registration.

## Usage

Without building it, just call:
```
go run cmd/foxpost/main.go <package number>
```

Example output:
```
Checking package: xxxxxxxxxxx
Found 2 updates, results:

Status: Úton
Date: 2025.06.30 13:35
Description: Csomagodat futárunk kivette a csomagautomatából.
-----------------------------
Status: Csomagod elkészült
Date: 2025.06.25 14:58
Description: Csomagod létrejött a rendszerünkben, a feladó még nem adta át azt a FOXPOST részére.
-----------------------------
```
