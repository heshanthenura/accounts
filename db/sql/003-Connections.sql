CREATE TABLE IF NOT EXISTS Connections (
    provider VARCHAR(20) PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS UserConnections (
    userId UUID NOT NULL,
    provider VARCHAR(20) NOT NULL,
    providerUserId VARCHAR(50) NOT NULL,
    providerAccountEmail VARCHAR(255) NOT NULL,
    linkedAt TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (userId, provider),
    FOREIGN KEY (userId) REFERENCES Users(id),
    FOREIGN KEY (provider) REFERENCES Connections(provider)
);
