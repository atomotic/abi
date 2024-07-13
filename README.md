# ABI

[Anagrafe Biblioteche Italiane](https://anagrafe.iccu.sbn.it/it/)

parquet dump: https://atomotic.github.io/abi/abi.parquet

[example queries](https://shell.duckdb.org/#queries=v0,CREATE-TABLE-libraries-AS-FROM-'https%3A%2F%2Fatomotic.github.io%2Fabi%2Fabi.parquet'~,DESCRIBE-libraries~,SELECT-id%2C-name-FROM-libraries-WHERE-location.comune.nome%3D'Bologna'-LIMIT-10~)
