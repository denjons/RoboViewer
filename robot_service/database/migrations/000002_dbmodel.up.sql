CREATE TABLE Sessions (
  id SERIAL PRIMARY KEY,
  uuid VARCHAR(36) UNIQUE NOT NULL,
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  robotId BIGINT REFERENCES Robots(id),
  floorId BIGINT REFERENCES Floors(id)
);

CREATE TABLE Points (
  sessionId BIGINT NOT NULL,
  sequence BIGINT NOT NULL,
  x INTEGER NOT NULL,
  y INTEGER NOT NULL,
  PRIMARY KEY(sessionId, sequence)
);

CREATE UNIQUE INDEX sessions_id_index ON Sessions (uuid);