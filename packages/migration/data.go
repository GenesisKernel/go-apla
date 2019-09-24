// Copyright (C) 2017, 2018, 2019 EGAAS S.A.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or (at
// your option) any later version.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301, USA.

package migration

//go:generate go run ./gen/contracts.go

var (
	migrationInitial = `
	{{headseq "migration_history"}}
		t.Column("id", "int", {"default_raw": "nextval('migration_history_id_seq')"})
		t.Column("version", "string", {"default": "", "size":255})
		t.Column("date_applied", "int", {})
	{{footer "seq" "primary"}}
`
	migrationInitialTables = `
	{{head "transactions_status"}}
		t.Column("hash", "bytea", {"default": ""})
		t.Column("time", "int", {"default": "0"})
		t.Column("type", "int", {"default": "0"})
		t.Column("ecosystem", "int", {"default": "1"})
		t.Column("wallet_id", "bigint", {"default": "0"})
		t.Column("block_id", "int", {"default": "0"})
		t.Column("error", "string", {"default": "", "size":255})
	{{footer}}
	sql("ALTER TABLE ONLY transactions_status ADD CONSTRAINT transactions_status_pkey PRIMARY KEY (hash);")
	`

	migrationInitialSchema = `
		
		DROP TABLE IF EXISTS "confirmations"; CREATE TABLE "confirmations" (
		"block_id" bigint  NOT NULL DEFAULT '0',
		"good" int  NOT NULL DEFAULT '0',
		"bad" int  NOT NULL DEFAULT '0',
		"time" int  NOT NULL DEFAULT '0'
		);
		ALTER TABLE ONLY "confirmations" ADD CONSTRAINT confirmations_pkey PRIMARY KEY (block_id);
		
		DROP TABLE IF EXISTS "block_chain"; CREATE TABLE "block_chain" (
		"id" int NOT NULL DEFAULT '0',
		"hash" bytea  NOT NULL DEFAULT '',
		"rollbacks_hash" bytea NOT NULL DEFAULT '',
		"data" bytea NOT NULL DEFAULT '',
		"ecosystem_id" int  NOT NULL DEFAULT '0',
		"key_id" bigint  NOT NULL DEFAULT '0',
		"node_position" bigint  NOT NULL DEFAULT '0',
		"time" int NOT NULL DEFAULT '0',
		"tx" int NOT NULL DEFAULT '0'
		);
		ALTER TABLE ONLY "block_chain" ADD CONSTRAINT block_chain_pkey PRIMARY KEY (id);
		
		DROP TABLE IF EXISTS "log_transactions"; CREATE TABLE "log_transactions" (
		"hash" bytea  NOT NULL DEFAULT '',
		"block" int NOT NULL DEFAULT '0'
		);
		ALTER TABLE ONLY "log_transactions" ADD CONSTRAINT log_transactions_pkey PRIMARY KEY (hash);

		DROP TABLE IF EXISTS "queue_tx"; CREATE TABLE "queue_tx" (
		"hash" bytea  NOT NULL DEFAULT '',
		"data" bytea NOT NULL DEFAULT '',
		"from_gate" int NOT NULL DEFAULT '0'
		);
		ALTER TABLE ONLY "queue_tx" ADD CONSTRAINT queue_tx_pkey PRIMARY KEY (hash);
		
		CREATE TABLE "system_contracts" (
		"id" bigint NOT NULL  DEFAULT '0',
		"value" text  NOT NULL DEFAULT '',
		"wallet_id" bigint NOT NULL DEFAULT '0',
		"token_id" bigint NOT NULL DEFAULT '0',
		"active" character(1) NOT NULL DEFAULT '0',
		"conditions" text  NOT NULL DEFAULT ''
		);
		ALTER TABLE ONLY "system_contracts" ADD CONSTRAINT system_contracts_pkey PRIMARY KEY (id);
		
		
		CREATE TABLE "system_tables" (
		"name" varchar(100)  NOT NULL DEFAULT '',
		"permissions" jsonb,
		"columns" jsonb,
		"conditions" text  NOT NULL DEFAULT ''
		);
		ALTER TABLE ONLY "system_tables" ADD CONSTRAINT system_tables_pkey PRIMARY KEY (name);
		
		DROP TABLE IF EXISTS "info_block"; CREATE TABLE "info_block" (
		"hash" bytea  NOT NULL DEFAULT '',
		"rollbacks_hash" bytea NOT NULL DEFAULT '',
		"block_id" int NOT NULL DEFAULT '0',
		"node_position" int  NOT NULL DEFAULT '0',
		"ecosystem_id" bigint NOT NULL DEFAULT '0',
		"key_id" bigint NOT NULL DEFAULT '0',
		"time" int  NOT NULL DEFAULT '0',
		"current_version" varchar(50) NOT NULL DEFAULT '0.0.1',
		"sent" smallint NOT NULL DEFAULT '0'
		);

		DROP TABLE IF EXISTS "queue_blocks"; CREATE TABLE "queue_blocks" (
		"hash" bytea  NOT NULL DEFAULT '',
		"full_node_id" bigint NOT NULL DEFAULT '0',
		"block_id" int NOT NULL DEFAULT '0'
		);
		ALTER TABLE ONLY "queue_blocks" ADD CONSTRAINT queue_blocks_pkey PRIMARY KEY (hash);
		
		DROP TABLE IF EXISTS "transactions"; CREATE TABLE "transactions" (
		"hash" bytea  NOT NULL DEFAULT '',
		"data" bytea NOT NULL DEFAULT '',
		"used" smallint NOT NULL DEFAULT '0',
		"high_rate" smallint NOT NULL DEFAULT '0',
		"type" smallint NOT NULL DEFAULT '0',
		"key_id" bigint NOT NULL DEFAULT '0',
		"counter" smallint NOT NULL DEFAULT '0',
		"sent" smallint NOT NULL DEFAULT '0',
		"attempt" smallint NOT NULL DEFAULT '0',
		"verified" smallint NOT NULL DEFAULT '1'
		);
		ALTER TABLE ONLY "transactions" ADD CONSTRAINT transactions_pkey PRIMARY KEY (hash);
		
		DROP SEQUENCE IF EXISTS rollback_tx_id_seq CASCADE;
		CREATE SEQUENCE rollback_tx_id_seq START WITH 1;
		DROP TABLE IF EXISTS "rollback_tx"; CREATE TABLE "rollback_tx" (
		"id" bigint NOT NULL  default nextval('rollback_tx_id_seq'),
		"block_id" bigint NOT NULL DEFAULT '0',
		"tx_hash" bytea  NOT NULL DEFAULT '',
		"table_name" varchar(255) NOT NULL DEFAULT '',
		"table_id" varchar(255) NOT NULL DEFAULT '',
		"data" TEXT NOT NULL DEFAULT ''
		);
		ALTER SEQUENCE rollback_tx_id_seq owned by rollback_tx.id;
		ALTER TABLE ONLY "rollback_tx" ADD CONSTRAINT rollback_tx_pkey PRIMARY KEY (id);
		CREATE INDEX "rollback_tx_table" ON "rollback_tx" (table_name, table_id);


		DROP TABLE IF EXISTS "install"; CREATE TABLE "install" (
		"progress" varchar(10) NOT NULL DEFAULT ''
		);
		
		
		DROP TYPE IF EXISTS "my_node_keys_enum_status" CASCADE;
		CREATE TYPE "my_node_keys_enum_status" AS ENUM ('my_pending','approved');
		DROP SEQUENCE IF EXISTS my_node_keys_id_seq CASCADE;
		CREATE SEQUENCE my_node_keys_id_seq START WITH 1;
		DROP TABLE IF EXISTS "my_node_keys"; CREATE TABLE "my_node_keys" (
		"id" int NOT NULL  default nextval('my_node_keys_id_seq'),
		"add_time" int NOT NULL DEFAULT '0',
		"public_key" bytea  NOT NULL DEFAULT '',
		"private_key" varchar(3096) NOT NULL DEFAULT '',
		"status" my_node_keys_enum_status  NOT NULL DEFAULT 'my_pending',
		"my_time" int NOT NULL DEFAULT '0',
		"time" bigint NOT NULL DEFAULT '0',
		"block_id" int NOT NULL DEFAULT '0'
		);
		ALTER SEQUENCE my_node_keys_id_seq owned by my_node_keys.id;
		ALTER TABLE ONLY "my_node_keys" ADD CONSTRAINT my_node_keys_pkey PRIMARY KEY (id);
		
		DROP TABLE IF EXISTS "stop_daemons"; CREATE TABLE "stop_daemons" (
		"stop_time" int NOT NULL DEFAULT '0'
		);
		
		CREATE OR REPLACE FUNCTION next_id(table_name TEXT, OUT result INT) AS
		$$
		BEGIN
			EXECUTE FORMAT('SELECT COUNT(*) + 1 FROM "%s"', table_name)
			INTO result;
			RETURN;
		END
		$$
		LANGUAGE plpgsql;`
)
