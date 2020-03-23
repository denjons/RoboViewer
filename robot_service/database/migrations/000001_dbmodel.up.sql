CREATE TABLE Floors (
  id VARCHAR(40) PRIMARY KEY,
  name VARCHAR(32) NOT NULL,
  width INTEGER NOT NULL,
  grid INTEGER[] NOT NULL
);

CREATE TABLE Robots (
  id VARCHAR(40) PRIMARY KEY,
  name VARCHAR(32) NOT NULL,
  width INTEGER NOT NULL,
  grid INTEGER[] NOT NULL,
  floorId VARCHAR(40) REFERENCES Floors(id)
);

CREATE UNIQUE INDEX floor_id_index ON Floors (id);

CREATE UNIQUE INDEX robot_id_index ON Robots (id);