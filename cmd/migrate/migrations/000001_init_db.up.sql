CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TABLE IF NOT EXISTS listado_definitivo
(
  id bigserial primary key,
  nombre text not null,
  rfc text not null unique,
  fecha_publicacion_sat_definitivos_text text,
  fecha_publicacion_dof_definitivos_text text,
  fecha_publicacion_sat_definitivos timestamptz,
  fecha_publicacion_dof_definitivos timestamptz
);

CREATE INDEX trgm_listado_definitivo_idx ON listado_definitivo USING GIST (nombre gist_trgm_ops);