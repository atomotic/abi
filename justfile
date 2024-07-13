default:
    just -l 

to-jsonl:
    jq -c ".biblioteche[]" < sources/biblioteche.json > sources/biblioteche.jsonl


