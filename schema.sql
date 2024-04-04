-- logs definition

CREATE TABLE if not exists logs (
	id INTEGER NOT NULL,
	msg TEXT,
	time INTEGER, 
    "level" TEXT,
	CONSTRAINT logs_pk PRIMARY KEY (id)
);


-- log_attrs definition

CREATE TABLE if not exists log_attrs (
	log_id INTEGER,
	"key" TEXT,
	value TEXT,
	CONSTRAINT log_attrs_logs_FK FOREIGN KEY (log_id) REFERENCES logs(id)
);