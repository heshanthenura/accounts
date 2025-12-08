CREATE TABLE IF NOT EXISTS Roles (
    name VARCHAR(20) PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS UserRoles (
    userId UUID NOT NULL,
    roleName VARCHAR(20) NOT NULL,
    PRIMARY KEY (userId, roleName),
    FOREIGN KEY (userId) REFERENCES Users(id),
    FOREIGN KEY (roleName) REFERENCES Roles(name)
);

INSERT INTO Roles (name) VALUES ('admin');