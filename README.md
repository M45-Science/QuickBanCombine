# QuickBanCombine
[![License: Unlicense](https://img.shields.io/badge/license-Unlicense-blue.svg)](http://unlicense.org/)
[![Go](https://github.com/Distortions81/QuickBanCombine/actions/workflows/go.yml/badge.svg)](https://github.com/Distortions81/QuickBanCombine/actions/workflows/go.yml)
[![ReportCard](https://github.com/Distortions81/QuickBanCombine/actions/workflows/report.yml/badge.svg)](https://github.com/Distortions81/QuickBanCombine/actions/workflows/report.yml)
[![CodeQL](https://github.com/Distortions81/QuickBanCombine/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/Distortions81/QuickBanCombine/actions/workflows/codeql-analysis.yml)
[![BinaryBuild](https://github.com/Distortions81/QuickBanCombine/actions/workflows/build-linux64.yml/badge.svg)](https://github.com/Distortions81/QuickBanCombine/actions/workflows/build-linux64.yml)<br>
<br>Quickly combine multiple Factorio server-banlist.json files, removing duplicates.<br>

Run `QuickBanCombine file1.json file2.json ...` to write `composite.json`.
Inputs can contain username strings or objects with `username` and `reason`.
Usernames are combined without regard to case across all files. Duplicate bans
are retained once, with distinct nonempty reasons joined under a `[dup]` prefix.
Invalid inputs stop processing before the existing output is replaced.
