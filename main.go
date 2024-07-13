package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/cheggaaa/pb/v3"
	_ "github.com/marcboeker/go-duckdb"
)

var migrations = `SET autoinstall_known_extensions=true;
SET autoload_known_extensions=true;
CREATE TABLE libraries AS SELECT * FROM read_json_auto('sources/biblioteche.jsonl');
ALTER TABLE libraries ADD COLUMN id text; UPDATE libraries SET id = "codici-identificativi".isil;
ALTER TABLE libraries ADD COLUMN name text; UPDATE libraries SET name = denominazioni.ufficiale;
ALTER TABLE libraries RENAME "anno-censimento" TO "year";
ALTER TABLE libraries RENAME "data-aggiornamento" TO "last_update";
ALTER TABLE libraries RENAME "codici-identificativi" TO identifiers;
ALTER TABLE libraries RENAME "denominazioni" TO "names";
ALTER TABLE libraries RENAME "indirizzo" TO "location";
ALTER TABLE libraries RENAME "contatti" TO "contacts";
ALTER TABLE libraries RENAME "stato-registrazione" TO "status";
ALTER TABLE libraries RENAME "tipologia-amministrativa" TO "administrative-type";
ALTER TABLE libraries RENAME "tipologia-funzionale" TO "functional-type";
ALTER TABLE libraries RENAME "ente" TO "institution";
ALTER TABLE libraries ADD COLUMN "assets" STRUCT(asset VARCHAR, "type" VARCHAR, items VARCHAR)[];
ALTER TABLE libraries ADD COLUMN "fonds" STRUCT("name" VARCHAR, dewey STRUCT("label" VARCHAR, code VARCHAR), description VARCHAR)[];
CREATE INDEX id_idx ON libraries (id);`

func main() {
	db, err := sql.Open("duckdb", "")
	if err != nil {
		slog.Error("open db", "msg", err)
		os.Exit(1)
	}
	defer db.Close()

	fmt.Println("- load libraries")
	_, err = db.Exec(migrations)
	if err != nil {
		slog.Error("apply migrations", "msg", err)
		os.Exit(1)
	}

	fmt.Println("- load special fonds")
	fonds, err := ParseFonds("sources/fondi-speciali.xml")
	if err != nil {
		slog.Error("parse special fonds", "msg", err)
	}
	bar1 := pb.StartNew(len(fonds))
	for isil, fond := range fonds {
		_, err = db.Exec(`UPDATE libraries SET fonds = ? WHERE id = ?`, fond, isil)
		if err != nil {
			slog.Error("update fonds", "msg", err)
		}
		bar1.Increment()
	}
	bar1.Finish()

	fmt.Println("- load assets")
	assets, err := ParseAssets("sources/patrimonio.xml")
	if err != nil {
		slog.Error("parse assets", "msg", err)
	}
	bar2 := pb.StartNew(len(assets))
	for isil, asset := range assets {
		_, err = db.Exec(`UPDATE libraries SET assets = ? WHERE id = ?`, asset, isil)
		if err != nil {
			slog.Error("update assets", "msg", err)
		}
		bar2.Increment()
	}
	bar2.Finish()

	fmt.Println("- export to parquet")
	_, err = db.Exec("COPY libraries TO 'abi.parquet' (FORMAT PARQUET);")
	if err != nil {
		slog.Error("parquet export", "msg", err)
		os.Exit(1)
	}

}
